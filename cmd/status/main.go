package main

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/yushafro/03.12.2025/pkg/deferfunc"
	"github.com/yushafro/03.12.2025/pkg/logger"
)

const (
	servers = 1
	shTO    = 30
)

func main() {
	log := logger.New()
	ctx, stop := context.WithTimeout(logger.WithLogger(context.Background(), log), shTO)
	defer deferfunc.Close(ctx, log.Stop, "Error stopping logger")
	defer stop()

	cfg := initConfig(ctx)

	wg := &sync.WaitGroup{}
	shSrvCh := make(chan struct{})

	wg.Add(servers)
	go initServer(ctx, &server{
		cfg:     cfg,
		wg:      wg,
		shSrvCh: shSrvCh,
	})

	shCh := make(chan os.Signal, 1)
	signal.Notify(shCh, os.Interrupt, syscall.SIGINT)
	<-shCh
	log.Info(ctx, "shutdown signal received")

	close(shSrvCh)
	wg.Wait()
	log.Info(ctx, "all servers stopped")
}
