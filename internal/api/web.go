package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func (a *API) handlerUI(c *gin.Context) {
	ctx := c.Request.Context()

	data, err := a.service.WebDocs(ctx)
	if err != nil {
		a.logger.Error("failed to get topics data", zap.Error(err))
		c.String(http.StatusInternalServerError, err.Error())
		return
	}

	tplData := map[string]any{
		"Topics": data.Topics,
	}

	a.logger.Info("rendering web template", zap.Int("topics", len(data.Topics)))

	c.Header("Content-Type", "text/html; charset=utf-8")
	if err := a.tpl.ExecuteTemplate(c.Writer, "web", tplData); err != nil {
		a.logger.Error("failed to render template", zap.Error(err))
		c.String(http.StatusInternalServerError, "error rendering template")
	}
}
