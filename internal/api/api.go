package api

import (
	"context"
	"encoding/json"
	"fmt"
	"html/template"
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
	tpl        *template.Template
}

func NewAPI(port string, svc Service, logger *slog.Logger) *API {
	api := &API{service: svc, logger: logger, port: port}
	api.parseTemplates()
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

func (a *API) parseTemplates() {
	funcMap := template.FuncMap{
		"subtract": func(a, b int) int {
			return a - b
		},
	}
	tpl := template.Must(template.New("web").Funcs(funcMap).Parse(templateWeb))
	a.tpl = tpl
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

func (a *API) routes() {
	a.mux = http.NewServeMux()

	a.mux.HandleFunc("GET /", a.handlerUI)
	a.mux.HandleFunc("GET /ui/", a.handlerUI)
	a.mux.HandleFunc("POST /api/v1/config", a.handlerApplyConfig)
}

func (a *API) writeResponse(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func (a *API) handlerUI(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	data, err := a.service.WebDocs(ctx)
	if err != nil {
		a.logger.Error("failed to load topics data", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to load topics data",
		})
		return
	}

	tplData := map[string]any{
		"Topics": data.Topics,
	}

	a.logger.Info("rendering web template", slog.Int("topics", len(data.Topics)))

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := a.tpl.ExecuteTemplate(w, "web", tplData); err != nil {
		a.logger.Error("failed to render template", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, map[string]string{
			"error": "failed to render template",
		})
		return
	}
}

func (a *API) handlerApplyConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := eventdoctor.LoadSpecFromReader(r.Body)
	if err != nil {
		a.logger.Error("failed to load config", slog.Any("error", err))
		a.writeResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "invalid configuration",
			"details": err.Error(),
		})
		return
	}

	if err := cfg.Validate(); err != nil {
		a.logger.Error("invalid config", slog.Any("error", err))
		a.writeResponse(w, http.StatusBadRequest, map[string]string{
			"error":   "configuration validation failed",
			"details": err.Error(),
		})
		return
	}

	if err := a.service.SaveSpec(r.Context(), cfg); err != nil {
		a.logger.Error("failed to save config", slog.Any("error", err))
		a.writeResponse(w, http.StatusInternalServerError, map[string]string{
			"error":   "failed to persist configuration",
			"details": err.Error(),
		})
		return
	}

	a.writeResponse(w, http.StatusOK, map[string]string{
		"message": "config uploaded successfully",
	})
}
