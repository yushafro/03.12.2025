package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
	"github.com/yushafro/03.12.2025/pkg/deferfunc"
	"github.com/yushafro/03.12.2025/pkg/logger"
)

func Logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log := logger.New()
		defer deferfunc.Close(r.Context(), log.Stop, "Error stopping logger")

		id := r.Header.Get("X-Request-ID")
		if id == "" {
			id = uuid.NewString()
		}

		ctx := context.WithValue(r.Context(), logger.RequestID, id)
		ctx = logger.WithLogger(ctx, log)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
