package api

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/nicolascb/eventdoctor/internal/api/middlewares"
)

type API struct {
	service    Service
	mux        *http.ServeMux
	httpServer *http.Server
	logger     *slog.Logger
	port       string
}

func NewAPI(port string, svc Service, logger *slog.Logger) *API {
	api := &API{service: svc, logger: logger, port: port}
	api.routes()
	return api
}

func (a *API) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down API server")
	if err := a.httpServer.Shutdown(ctx); err != nil {
		return fmt.Errorf("error closing server: %w", err)
	}

	a.logger.Info("API server shut down gracefully")
	return nil
}

func (a *API) Run() error {
	a.logger.Info("starting API server", slog.String("port", a.port))
	a.httpServer = &http.Server{
		Addr:         a.port,
		Handler:      middlewares.Logging(a.logger, a.corsMiddleware(a.mux)),
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return a.httpServer.ListenAndServe()
}
