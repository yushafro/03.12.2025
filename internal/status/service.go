package status

import (
	"net/http"
	"time"
)

type service struct {
	repo   *fileRepository
	client *http.Client
}

func NewService(repo *fileRepository, clientTO time.Duration) *service {
	client := &http.Client{
		Timeout: clientTO,
	}

	return &service{
		repo:   repo,
		client: client,
	}
}
