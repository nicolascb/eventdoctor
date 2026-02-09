package api

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"

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
		ReadTimeout:  15 * 1000000000, // 15 seconds in nanoseconds
		WriteTimeout: 15 * 1000000000, // 15 seconds in nanoseconds
		IdleTimeout:  60 * 1000000000, // 60 seconds in nanoseconds
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
	a.mux.HandleFunc("OPTIONS /v1/config", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	a.mux.HandleFunc("OPTIONS /v1/producers", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	a.mux.HandleFunc("OPTIONS /v1/events", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))
	a.mux.HandleFunc("OPTIONS /v1/consumers", a.corsMiddleware(func(w http.ResponseWriter, r *http.Request) {}))

}

func (a *API) writeResponse(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

// stringPtr returns a pointer to the given string
func stringPtr(s string) *string {
	return &s
}

func (a *API) handlerApplyConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := eventdoctor.LoadSpecFromReader(r.Body)
	if err != nil {
		a.logger.Error("failed to load config", slog.Any("error", err))
		a.writeResponse(w, http.StatusBadRequest, ErrorResponse{
			Error:   "invalid configuration",
			Details: err.Error(),
		})
		return
	}

	if err := cfg.Validate(); err != nil {
		a.logger.Error("invalid config", slog.Any("error", err))
		a.writeResponse(w, http.StatusBadRequest, ErrorResponse{
			Error:   "configuration validation failed",
			Details: err.Error(),
		})
		return
	}

	if err := a.service.SaveSpec(r.Context(), cfg); err != nil {
		a.logger.Error("failed to save config", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, ErrorResponse{
			Error:   "failed to persist configuration",
			Details: err.Error(),
		})
		return
	}

	a.writeResponse(w, http.StatusOK, SuccessResponse{
		Message: "config uploaded successfully",
	})
}

func (a *API) handlerListProducers(w http.ResponseWriter, r *http.Request) {
	// Mock data for producers
	mockProducers := []ProducerResponse{
		{
			Topic:       "user.events",
			Name:        "user-service",
			Owner:       true,
			Writes:      true,
			Description: "Produces user-related events",
			Events: []EventResponse{
				{
					Name:        "user.created",
					Version:     stringPtr("1.0.0"),
					Description: "User creation event",
					SchemaURL:   "https://schemas.example.com/user.created.v1.json",
				},
				{
					Name:        "user.updated",
					Version:     stringPtr("1.0.0"),
					Description: "User update event",
					SchemaURL:   "https://schemas.example.com/user.updated.v1.json",
				},
			},
		},
		{
			Topic:       "payment.events",
			Name:        "payment-service",
			Owner:       true,
			Writes:      true,
			Description: "Produces payment-related events",
			Events: []EventResponse{
				{
					Name:        "payment.processed",
					Version:     stringPtr("1.0.0"),
					Description: "Payment processing event",
					SchemaURL:   "https://schemas.example.com/payment.processed.v1.json",
				},
			},
		},
	}

	if err := a.writeResponse(w, http.StatusOK, mockProducers); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerListEvents(w http.ResponseWriter, r *http.Request) {
	// Mock data for events grouped by topic
	mockTopics := []TopicWithEvents{
		{
			TopicName: "user.events",
			Events: []EventResponse{
				{
					Name:        "user.created",
					Version:     stringPtr("1.0.0"),
					Description: "User creation event",
					SchemaURL:   "https://schemas.example.com/user.created.v1.json",
					Headers: []EventHeaderResponse{
						{Name: "x-correlation-id", Description: "Correlation ID for tracing"},
						{Name: "x-source-service", Description: "Service that produced the event"},
					},
				},
				{
					Name:        "user.updated",
					Version:     stringPtr("1.0.0"),
					Description: "User update event",
					SchemaURL:   "https://schemas.example.com/user.updated.v1.json",
					Headers: []EventHeaderResponse{
						{Name: "x-correlation-id", Description: "Correlation ID for tracing"},
					},
				},
			},
		},
		{
			TopicName: "payment.events",
			Events: []EventResponse{
				{
					Name:        "payment.processed",
					Version:     stringPtr("1.0.0"),
					Description: "Payment processing event",
					SchemaURL:   "https://schemas.example.com/payment.processed.v1.json",
				},
			},
		},
		{
			TopicName: "order.events",
			Events: []EventResponse{
				{
					Name:        "order.completed",
					Version:     stringPtr("2.1.0"),
					Description: "Order completion event",
					SchemaURL:   "https://schemas.example.com/order.completed.v2.json",
				},
			},
		},
	}

	if err := a.writeResponse(w, http.StatusOK, mockTopics); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerListConsumers(w http.ResponseWriter, r *http.Request) {
	// Mock data for consumers
	mockConsumers := []ConsumerResponse{
		{
			Group:       "notification-service",
			Description: "Service responsible for sending notifications",
			Topics: []ConsumerTopicResponse{
				{
					Name: "user.events",
					Events: []ConsumerEventResponse{
						{Name: "user.created", Version: stringPtr("1.0.0")},
						{Name: "user.updated", Version: stringPtr("1.0.0")},
					},
				},
			},
		},
		{
			Group:       "analytics-service",
			Description: "Service for analytics and reporting",
			Topics: []ConsumerTopicResponse{
				{
					Name: "payment.events",
					Events: []ConsumerEventResponse{
						{Name: "payment.processed", Version: stringPtr("1.0.0")},
					},
				},
				{
					Name: "user.events",
					Events: []ConsumerEventResponse{
						{Name: "user.created", Version: stringPtr("1.0.0")},
					},
				},
			},
		},
		{
			Group:       "order-fulfillment",
			Description: "Order processing and fulfillment service",
			Topics: []ConsumerTopicResponse{
				{
					Name: "order.events",
					Events: []ConsumerEventResponse{
						{Name: "order.completed", Version: stringPtr("2.1.0")},
					},
				},
			},
		},
	}

	if err := a.writeResponse(w, http.StatusOK, mockConsumers); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}
