package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"time"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
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
		Handler:      a.mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	return a.httpServer.ListenAndServe()
}

// corsMiddleware adds CORS headers to handle cross-origin requests
func (a *API) corsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Set CORS headers
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
		w.Header().Set("Access-Control-Max-Age", "3600")

		// Handle preflight requests
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}

		next(w, r)
	}
}

func (a *API) routes() {
	a.mux = http.NewServeMux()
	a.mux.HandleFunc("POST /v1/config", a.corsMiddleware(a.handlerApplyConfig))
	a.mux.HandleFunc("GET /v1/producers", a.corsMiddleware(a.handlerListProducers))
	a.mux.HandleFunc("GET /v1/events", a.corsMiddleware(a.handlerListEvents))
	a.mux.HandleFunc("GET /v1/consumers", a.corsMiddleware(a.handlerListConsumers))
	a.mux.HandleFunc("GET /v1/overview", a.corsMiddleware(a.handlerOverview))
	a.mux.HandleFunc("OPTIONS /v1/config", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	a.mux.HandleFunc("OPTIONS /v1/producers", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	a.mux.HandleFunc("OPTIONS /v1/events", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	a.mux.HandleFunc("OPTIONS /v1/consumers", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	a.mux.HandleFunc("OPTIONS /v1/overview", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))

}

func (a *API) writeResponse(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (a *API) handlerApplyConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := eventdoctor.LoadSpecFromReader(r.Body)
	if err != nil {
		a.logger.Error("failed to load config", slog.Any("error", err))
		a.writeResponse(w, http.StatusBadRequest, response.ErrorResponse{
			Error:   "invalid configuration",
			Details: err.Error(),
		})
		return
	}

	if err := cfg.Validate(); err != nil {
		a.logger.Error("invalid config", slog.Any("error", err))
		a.writeResponse(w, http.StatusBadRequest, response.ErrorResponse{
			Error:   "configuration validation failed",
			Details: err.Error(),
		})
		return
	}

	if err := a.service.SaveSpec(r.Context(), cfg); err != nil {
		a.logger.Error("failed to save config", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to persist configuration",
			Details: err.Error(),
		})
		return
	}

	a.writeResponse(w, http.StatusOK, response.SuccessResponse{
		Message: "config uploaded successfully",
	})
}

func (a *API) handlerListProducers(w http.ResponseWriter, r *http.Request) {
	producers, err := a.service.ListProducers(r.Context())
	if err != nil {
		a.logger.Error("failed to list producers", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list producers",
			Details: err.Error(),
		})
		return
	}

	if err := a.writeResponse(w, http.StatusOK, producers); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerListEvents(w http.ResponseWriter, r *http.Request) {
	events, err := a.service.ListEvents(r.Context())
	if err != nil {
		a.logger.Error("failed to list events", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list events",
			Details: err.Error(),
		})
		return
	}

	if err := a.writeResponse(w, http.StatusOK, events); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerListConsumers(w http.ResponseWriter, r *http.Request) {
	consumers, err := a.service.ListConsumers(r.Context())
	if err != nil {
		a.logger.Error("failed to list consumers", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list consumers",
			Details: err.Error(),
		})
		return
	}

	if err := a.writeResponse(w, http.StatusOK, consumers); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := a.service.Overview(r.Context())
	if err != nil {
		a.logger.Error("failed to get overview", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to get overview",
			Details: err.Error(),
		})
		return
	}

	if err := a.writeResponse(w, http.StatusOK, overview); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}
