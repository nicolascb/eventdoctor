package api

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (a *API) handlerUI(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := a.service.WebDocs(ctx)
	if err != nil {
		a.logger.Error("failed to get topics data", slog.Any("error", err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tplData := map[string]any{
		"Topics": data.Topics,
	}

	a.logger.Info("rendering web template", slog.Int("topics", len(data.Topics)))

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := a.tpl.ExecuteTemplate(c.Writer, "web", tplData); err != nil {
		a.logger.Error("failed to render template", slog.Any("error", err))
		c.String(http.StatusInternalServerError, "error rendering template")
	}
}
