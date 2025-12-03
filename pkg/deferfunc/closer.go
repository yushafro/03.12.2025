package deferfunc

import (
	"context"
	"errors"
	"syscall"

	"github.com/yushafro/03.12.2025/pkg/logger"
	"go.uber.org/zap"
)

func Close(ctx context.Context, c func() error, errMsg string) {
	log := logger.FromContext(ctx)

	err := c()
	if err != nil && !errors.Is(err, syscall.EINVAL) {
		log.Error(ctx, errMsg, zap.Error(err))
	}
}
