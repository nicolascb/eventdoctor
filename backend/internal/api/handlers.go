package api

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/nicolascb/eventdoctor/internal/api/response"
	"github.com/nicolascb/eventdoctor/internal/eventdoctor"
)

func (a *API) handlerApplyConfig(w http.ResponseWriter, r *http.Request) {
	cfg, err := eventdoctor.LoadSpecFromReader(r.Body)
	if err != nil {
		a.logger.Error("failed to load config", slog.Any("error", err))
		response.Write(w, http.StatusBadRequest, response.ErrorResponse{
			Error:   "invalid configuration",
			Details: err.Error(),
		})
		return
	}

	if err := cfg.Validate(); err != nil {
		a.logger.Error("invalid config", slog.Any("error", err))
		response.Write(w, http.StatusBadRequest, response.ErrorResponse{
			Error:   "configuration validation failed",
			Details: err.Error(),
		})
		return
	}

	if err := a.service.SaveSpec(r.Context(), cfg); err != nil {
		a.logger.Error("failed to save config", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to persist configuration",
			Details: err.Error(),
		})
		return
	}

	response.Write(w, http.StatusOK, response.SuccessResponse{
		Message: "config uploaded successfully",
	})
}

func (a *API) handlerListProducers(w http.ResponseWriter, r *http.Request) {
	pagination := getPaginationFromRequest(r)
	search := r.URL.Query().Get("search")
	producers, err := a.service.ListProducers(r.Context(), pagination.Page, pagination.PageSize, search)
	if err != nil {
		a.logger.Error("failed to list producers", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list producers",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, producers); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerGetProducerDetail(w http.ResponseWriter, r *http.Request) {
	serviceIDStr := r.PathValue("service_id")
	topicIDStr := r.PathValue("topic_id")

	serviceID, err := strconv.ParseInt(serviceIDStr, 10, 64)
	if err != nil {
		response.Write(w, http.StatusBadRequest, response.ErrorResponse{
			Error: "invalid service id",
		})
		return
	}

	topicID, err := strconv.ParseInt(topicIDStr, 10, 64)
	if err != nil {
		response.Write(w, http.StatusBadRequest, response.ErrorResponse{
			Error: "invalid topic id",
		})
		return
	}

	detail, err := a.service.GetProducerDetail(r.Context(), serviceID, topicID)
	if err != nil {
		a.logger.Error("failed to get producer detail", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to get producer detail",
			Details: err.Error(),
		})
		return
	}

	if detail == nil {
		response.Write(w, http.StatusNotFound, response.ErrorResponse{
			Error: "producer not found",
		})
		return
	}

	if err := response.Write(w, http.StatusOK, detail); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerListEvents(w http.ResponseWriter, r *http.Request) {
	pagination := getPaginationFromRequest(r)
	search := r.URL.Query().Get("search")
	events, err := a.service.ListEvents(r.Context(), pagination.Page, pagination.PageSize, search)
	if err != nil {
		a.logger.Error("failed to list events", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list events",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, events); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerListConsumers(w http.ResponseWriter, r *http.Request) {
	undocumentedOnly := r.URL.Query().Get("undocumented_only") == "true"

	if undocumentedOnly {
		a.handlerUndocumentedConsumers(w, r)
		return
	}

	pagination := getPaginationFromRequest(r)
	search := r.URL.Query().Get("search")
	consumers, err := a.service.ListConsumers(r.Context(), pagination.Page, pagination.PageSize, search)
	if err != nil {
		a.logger.Error("failed to list consumers", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list consumers",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, consumers); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerOverview(w http.ResponseWriter, r *http.Request) {
	overview, err := a.service.Overview(r.Context())
	if err != nil {
		a.logger.Error("failed to get overview", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to get overview",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, overview); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerGetTopic(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	view, err := a.service.GetTopicView(r.Context(), name)
	if err != nil {
		a.logger.Error("failed to get topic", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to get topic",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, view); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerGetService(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	view, err := a.service.GetServiceView(r.Context(), name)
	if err != nil {
		a.logger.Error("failed to get service", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to get service",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, view); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerListTopics(w http.ResponseWriter, r *http.Request) {
	pagination := getPaginationFromRequest(r)
	topics, err := a.service.ListTopics(r.Context(), pagination.Page, pagination.PageSize)
	if err != nil {
		a.logger.Error("failed to list topics", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list topics",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, topics); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}

func (a *API) handlerUndocumentedConsumers(w http.ResponseWriter, r *http.Request) {
	result, err := a.service.ListUndocumentedConsumers(r.Context())
	if err != nil {
		a.logger.Error("failed to list undocumented consumers", slog.Any("error", err))
		response.Write(w, http.StatusInternalServerError, response.ErrorResponse{
			Error:   "failed to list undocumented consumers",
			Details: err.Error(),
		})
		return
	}

	if err := response.Write(w, http.StatusOK, result); err != nil {
		a.logger.Error("failed to write response", slog.Any("error", err))
	}
}
