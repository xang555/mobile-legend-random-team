package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/laoitdev/random-ml-team/internal/http/handlers"
)

// New constructs the HTTP router with standard middleware and routes.
func New(logger *zap.Logger, teamHandler *handlers.TeamHandler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(loggerMiddleware(logger))

	r.Get("/healthz", healthz)
	r.Get("/docs/*", httpSwagger.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		teamHandler.Register(r)
	})

	return r
}

func loggerMiddleware(logger *zap.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			next.ServeHTTP(w, r)
			duration := time.Since(start)
			logger.Info("request completed",
				zap.String("method", r.Method),
				zap.String("path", r.URL.Path),
				zap.Duration("duration", duration),
			)
		})
	}
}

// healthz godoc
// @Summary Health check
// @Tags health
// @Produce plain
// @Success 200 {string} string "ok"
// @Router /healthz [get]
func healthz(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
