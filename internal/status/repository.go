package status

import (
	"context"
	"encoding/json"
	"io"
	"os"
	"sync"

	"github.com/yushafro/03.12.2025/pkg/deferfunc"
	"github.com/yushafro/03.12.2025/pkg/logger"
	"go.uber.org/zap"
)

type repository interface {
	RegisterLinks(ctx context.Context, links responseRegisterLinks) (responseRegisterLinks, error)
	ListLinks(ctx context.Context, list requestLinksList) (responseLinksList, error)
}

type (
	fileRepositoryData map[int]Links
	fileRepository     struct {
		repository

		path string
		mu   *sync.Mutex
	}
)

func NewFileRepository(path string) *fileRepository {
	mu := &sync.Mutex{}

	return &fileRepository{
		path: path,
		mu:   mu,
	}
}

func (fr *fileRepository) read(ctx context.Context) (fileRepositoryData, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		log := logger.FromContext(ctx)

		file, err := os.Open(fr.path)
		if err != nil {
			if os.IsNotExist(err) {
				log.Info(ctx, "file does not exist")

				return make(fileRepositoryData), nil
			}

			log.Error(ctx, "Error opening file", zap.Error(err))

			return nil, err
		}
		defer deferfunc.Close(ctx, file.Close, "Error closing file")

		bytes, err := io.ReadAll(file)
		if err != nil {
			log.Error(ctx, "Error reading file", zap.Error(err))

			return nil, err
		}

		data := make(fileRepositoryData)
		err = json.Unmarshal(bytes, &data)
		if err != nil {
			log.Error(ctx, "Error unmarshalling file", zap.Error(err))

			return nil, err
		}

		log.Info(ctx, "file read")

		return data, nil
	}
}

func (fr *fileRepository) write(ctx context.Context, data fileRepositoryData) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		log := logger.FromContext(ctx)

		file, err := os.Create(fr.path)
		if err != nil {
			log.Error(ctx, "Error creating file", zap.Error(err))

			return err
		}
		defer deferfunc.Close(ctx, file.Close, "Error closing file")

		bytes, err := json.Marshal(data)
		if err != nil {
			log.Error(ctx, "Error marshalling file", zap.Error(err))

			return err
		}

		_, err = file.Write(bytes)
		if err != nil {
			log.Error(ctx, "Error writing file", zap.Error(err))

			return err
		}

		log.Info(ctx, "file written")

		return nil
	}
}
