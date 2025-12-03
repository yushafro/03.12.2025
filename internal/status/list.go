package status

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"maps"
	"net/http"

	"github.com/yushafro/03.12.2025/pkg/deferfunc"
	"github.com/yushafro/03.12.2025/pkg/logger"
	"go.uber.org/zap"
)

type requestLinksList struct {
	LinksList []int `json:"links_list"`
}
type responseLinksList struct {
	pdfBytes []byte
}

func (s *server) ListLinks(w http.ResponseWriter, r *http.Request) {
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

		var list requestLinksList
		err = json.Unmarshal(body, &list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(
				r.Context(),
				"Error unmarshalling request body",
				zap.Error(err),
			)

			return
		}

		if len(list.LinksList) == 0 {
			http.Error(w, "Request body is empty", http.StatusBadRequest)
			log.Error(r.Context(), "Request body is empty")

			return
		}

		resp, err := s.service.ListLinks(r.Context(), list)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			log.Error(
				r.Context(),
				"Error listing links",
				zap.Error(err),
			)

			return
		}

		if len(resp.pdfBytes) == 0 {
			http.Error(w, "No links found", http.StatusNotFound)
			log.Info(r.Context(), "No links found")

			return
		}

		w.Write(resp.pdfBytes)
		log.Info(r.Context(), "links listed")
	}
}

func (s *service) ListLinks(ctx context.Context, list requestLinksList) (responseLinksList, error) {
	select {
	case <-ctx.Done():
		return responseLinksList{}, ctx.Err()
	default:
		links, err := s.repo.ListLinks(ctx, list)
		if err != nil {
			return responseLinksList{}, err
		}

		pdfBytes, err := generatePDFReport(ctx, links)
		if err != nil {
			return responseLinksList{}, fmt.Errorf("failed to generate PDF: %w", err)
		}

		return responseLinksList{
			pdfBytes: pdfBytes,
		}, nil
	}
}

func (fr *fileRepository) ListLinks(
	ctx context.Context,
	list requestLinksList,
) (Links, error) {
	select {
	case <-ctx.Done():
		return Links{}, ctx.Err()
	default:
		fr.mu.Lock()
		data, err := fr.read(ctx)
		fr.mu.Unlock()

		if err != nil {
			return Links{}, err
		}

		links := make(Links)
		for _, num := range list.LinksList {
			l := data[num]

			maps.Copy(links, l)
		}

		return links, nil
	}
}
