package router

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/zap"

	"github.com/laoitdev/random-ml-team/internal/http/handlers"
	"github.com/laoitdev/random-ml-team/internal/http/handlers/health"
)

// New constructs the HTTP router with standard middleware and routes.
func New(logger *zap.Logger, metrics *health.Metrics, teamHandler *handlers.TeamHandler, healthHandler *health.Handler) http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Recoverer)
	r.Use(middleware.Timeout(60 * time.Second))
	r.Use(loggerMiddleware(logger, metrics))

	registerHealthRoutes(r, healthHandler)
	r.Get("/docs/*", httpSwagger.Handler())

	r.Route("/api/v1", func(r chi.Router) {
		teamHandler.Register(r)
		registerHealthRoutes(r, healthHandler)
	})

	return r
}

func loggerMiddleware(logger *zap.Logger, metrics *health.Metrics) func(http.Handler) http.Handler {
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
			if metrics != nil {
				metrics.IncrementRequests()
			}
		})
	}
}

func registerHealthRoutes(r chi.Router, healthHandler *health.Handler) {
	r.Get("/healthz", healthHandler.Healthz)
	r.Route("/health", func(r chi.Router) {
		r.Get("/", healthHandler.Root)
		r.Get("/live", healthHandler.Live)
		r.Get("/ready", healthHandler.Ready)
		r.Get("/status", healthHandler.Status)
	})
}
