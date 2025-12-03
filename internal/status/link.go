package status

import (
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"

	"github.com/yushafro/03.12.2025/pkg/deferfunc"
	"github.com/yushafro/03.12.2025/pkg/logger"
	"go.uber.org/zap"
)

type (
	requestRegisterLinks struct {
		Links []string `json:"links"`
	}
	responseRegisterLinks struct {
		Links    Links `json:"links"`
		LinksNum int   `json:"links_num"`
	}
)

func (s *server) RegisterLinks(w http.ResponseWriter, r *http.Request) {
	select {
	case <-r.Context().Done():
		http.Error(w, "Request canceled", http.StatusInternalServerError)
	default:
		log := logger.FromContext(r.Context())

		if r.ContentLength == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			log.Error(r.Context(), "Request body is empty")

			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(
				r.Context(),
				"Error reading request body",
				zap.Error(err),
			)

			return
		}
		defer deferfunc.Close(r.Context(), r.Body.Close, "Error closing request body")

		var links requestRegisterLinks
		err = json.Unmarshal(body, &links)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			log.Error(
				r.Context(),
				"Error unmarshalling request body",
				zap.Error(err),
				zap.String("body", string(body)),
			)

			return
		}

		if len(links.Links) == 0 {
			http.Error(w, "links is empty", http.StatusBadRequest)
			log.Error(r.Context(), "links is empty")

			return
		}

		resp, err := s.service.RegisterLinks(r.Context(), links)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(r.Context(), "Error registering links", zap.Error(err))

			return
		}

		body, err = json.Marshal(resp)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)

			return
		}

		_, err = w.Write(body)
		if err != nil {
			http.Error(w, "Error writing response body", http.StatusInternalServerError)
		}
	}
}

func (s *service) RegisterLinks(
	ctx context.Context,
	links requestRegisterLinks,
) (responseRegisterLinks, error) {
	select {
	case <-ctx.Done():
		return responseRegisterLinks{}, ctx.Err()
	default:
		log := logger.FromContext(ctx)
		result := responseRegisterLinks{
			Links: make(Links),
		}

		for _, link := range links.Links {
			u, err := url.Parse(link)
			if err != nil {
				return responseRegisterLinks{}, err
			}

			if u.Scheme == "" {
				u.Scheme = "https"
				link = u.String()
			}

			resp, err := s.client.Get(u.String())
			if err != nil {
				if errors.Is(err, context.DeadlineExceeded) {
					result.Links[link] = unavailable
					log.Info(
						ctx,
						"Link is not available (Deadline exceeded)",
						zap.String("link", link),
					)

					continue
				}

				return responseRegisterLinks{}, err
			}
			defer deferfunc.Close(ctx, resp.Body.Close, "Error closing response body")

			if resp.StatusCode == http.StatusOK {
				result.Links[link] = available

				log.Info(ctx, "Link is available", zap.String("link", link))

				continue
			}

			result.Links[link] = unavailable
			log.Info(ctx, "Link is not available", zap.String("link", link))
		}

		result, err := s.repo.RegisterLinks(ctx, result)
		if err != nil {
			return result, err
		}

		return result, nil
	}
}

func (fr *fileRepository) RegisterLinks(
	ctx context.Context,
	links responseRegisterLinks,
) (responseRegisterLinks, error) {
	fr.mu.Lock()
	data, err := fr.read(ctx)
	fr.mu.Unlock()

	if err != nil {
		return links, err
	}

	links.LinksNum = len(data) + 1

	fr.mu.Lock()
	data[links.LinksNum] = links.Links
	err = fr.write(ctx, data)
	fr.mu.Unlock()

	if err != nil {
		return links, err
	}

	return links, nil
}
