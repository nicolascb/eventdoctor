package api

import (
	"context"
	"html/template"
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
	tpl    *template.Template
}

func NewAPI(port string, store Store, logger *zap.Logger) *API {
	api := &API{store: store, logger: logger, port: port}
	api.parseTemplates()
	api.routes()
	return api
}

func (a *API) parseTemplates() {
	// templateDocs definido em templates.go (go:embed)
	tpl := template.Must(template.New("docs").Parse(templateDocs))
	a.tpl = tpl
}

func (a *API) Run() error {
	a.logger.Info("Starting API server", zap.String("port", a.port))
	return a.router.Run(a.port)
}

func (a *API) routes() {
	a.router = gin.Default()
	a.router.POST("/config", a.handleUploadConfig)
	a.router.GET("/docs", a.handleDocs)
	a.router.GET("/", a.handleDocs)
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

func (a *API) handleDocs(c *gin.Context) {
	if a.store == nil {
		c.String(http.StatusInternalServerError, "store não configurado")
		return
	}

	type ProducerView struct {
		Name       string
		Repository *string
	}
	type ConsumerView struct {
		Name       string
		Repository *string
	}
	type EventView struct {
		Type       string
		Version    *string
		SchemaURL  string
		Deprecated bool
		Producers  []ProducerView
		Consumers  []ConsumerView
	}
	type TopicView struct {
		Name   string
		Owner  string
		Events []EventView
	}

	topics, err := a.store.ListTopics()
	if err != nil {
		c.String(http.StatusInternalServerError, "erro listando tópicos")
		return
	}

	var data []TopicView
	for _, t := range topics {
		evs, err := a.store.GetEventsByTopic(t.ID)
		if err != nil {
			c.String(http.StatusInternalServerError, "erro listando eventos")
			return
		}
		tv := TopicView{Name: t.Name, Owner: t.Owner}
		for _, ev := range evs {
			prods, _ := a.store.GetProducersByEvent(ev.ID)
			cons, _ := a.store.GetConsumersByEvent(ev.ID)

			evv := EventView{
				Type:       ev.EventType,
				Version:    ev.SchemaVersion,
				SchemaURL:  ev.SchemaURL,
				Deprecated: ev.Deprecated,
			}
			for _, p := range prods {
				evv.Producers = append(evv.Producers, ProducerView{Name: p.Producer, Repository: p.Repository})
			}
			for _, cn := range cons {
				evv.Consumers = append(evv.Consumers, ConsumerView{Name: cn.Consumer, Repository: cn.Repository})
			}
			tv.Events = append(tv.Events, evv)
		}
		data = append(data, tv)
	}

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := a.tpl.ExecuteTemplate(c.Writer, "docs", map[string]any{
		"Topics": data,
	}); err != nil {
		c.String(http.StatusInternalServerError, "erro renderizando template")
	}
}

func (a *API) saveSpec(ctx context.Context, cfg config.Config) error {
	// Implementar lógica para salvar a spec via store
	return nil
}
