package api

import (
	"context"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/nicolascb/eventdoctor/internal/config"
	"go.uber.org/zap"
)

type API struct {
	store  Store
	router *gin.Engine
	logger *zap.Logger
	port   string
}

func NewAPI(port string, store Store, logger *zap.Logger) *API {
	api := &API{store: store, logger: logger, port: port}
	api.routes()
	return api
}

func (a *API) Run() error {
	a.logger.Info("Starting API server", zap.String("port", a.port))
	return a.router.Run(a.port)
}

func (a *API) routes() {
	a.router = gin.Default()
	a.router.POST("/config", a.handleUploadConfig)
}

func (a *API) handleUploadConfig(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "campo 'file' obrigatório"})
		return
	}

	f, err := fileHeader.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao abrir arquivo"})
		return
	}
	defer f.Close()

	cfg, err := config.LoadConfigFromReader(f)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "falha ao parsear configuração YAML"})
		return
	}

	a.logger.Info("Config parsed", zap.Any("config", cfg))

	if err := a.saveSpec(c.Request.Context(), cfg); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "falha ao persistir configuração"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "config uploaded successfully",
		"size":    c.Request.ContentLength,
	})
}

func (a *API) saveSpec(ctx context.Context, cfg config.Config) error {
	// Implementar lógica para salvar a spec via store
	return nil
}
