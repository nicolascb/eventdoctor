package api

import (
	"context"
	"fmt"
	"html/template"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
	"go.uber.org/zap"
)

type API struct {
	service Service
	router  *gin.Engine
	http    *http.Server
	logger  *zap.Logger
	port    string
	tpl     *template.Template
}

func NewAPI(port string, svc Service, logger *zap.Logger) *API {
	api := &API{service: svc, logger: logger, port: port}
	api.parseTemplates()
	api.routes()
	return api
}

func (a *API) Shutdown(ctx context.Context) error {
	a.logger.Info("shutting down API server")
	if err := a.http.Shutdown(ctx); err != nil {
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
	a.logger.Info("starting API server", zap.String("port", a.port))
	a.http = &http.Server{
		Addr:    a.port,
		Handler: a.router,
	}

	return a.http.ListenAndServe()
}

func (a *API) routes() {
	a.router = gin.Default()

	ui := a.router.Group("/ui")
	ui.GET("/", a.handlerUI)

	api := a.router.Group("/api/v1")
	api.POST("/config", a.handlerApplyConfig)
}

func (a *API) handlerApplyConfig(c *gin.Context) {
	body := c.Request.Body

	cfg, err := eventdoctor.LoadSpecFromReader(body)
	if err != nil {
		a.logger.Error("failed to load config", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid configuration", "details": err.Error()})
		return
	}

	if err := cfg.Validate(); err != nil {
		a.logger.Error("invalid config", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "configuration validation failed", "details": err.Error()})
		return
	}

	if err := a.service.SaveSpec(c.Request.Context(), cfg); err != nil {
		a.logger.Error("failed to save config", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to persist configuration", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "config uploaded successfully",
	})
}
