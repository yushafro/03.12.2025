package main

import (
	"context"
	"errors"
	"net/http"
	"sync"

	"github.com/yushafro/03.12.2025/internal/config"
	"github.com/yushafro/03.12.2025/internal/status"
	"github.com/yushafro/03.12.2025/pkg/logger"
	"go.uber.org/zap"
)

type server struct {
	cfg     *config.Config
	wg      *sync.WaitGroup
	shSrvCh <-chan struct{}
}

func initServer(ctx context.Context, s *server) {
	log := logger.FromContext(ctx)

	repo := status.NewFileRepository(s.cfg.FileRepositoryPath)
	service := status.NewService(repo, s.cfg.ClientTO)
	server := status.NewServer(service, s.cfg)
	server.RegisterHandlers()

	go func() {
		<-s.shSrvCh
		log.Info(ctx, "server stopping")

		err := server.Stop(ctx)
		if err != nil {
			log.Fatal(ctx, "Error stopping server", zap.Error(err))
		}

		log.Info(ctx, "server stopped")
		s.wg.Done()
	}()

	log.Info(ctx, "server started")
	err := server.Start()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(ctx, "Error starting server", zap.Error(err))

		s.wg.Done()
	}
}
