package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

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

func (a *API) handleDocs(c *gin.Context) {
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
