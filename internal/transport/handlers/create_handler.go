package handlers

import (
	"url-shortener/internal/logger"
	"url-shortener/internal/metrics"
	"url-shortener/internal/metrics/prometheus"
	"url-shortener/internal/service"

	json "github.com/json-iterator/go"
	"github.com/valyala/fasthttp"
)

type Creator interface {
	Create(req *service.Request) (service.Response, error)
}

type CreateHandler struct {
	baseHandler
	shortURLCreator *service.URLShortener
	logger          *logger.Logger
	metricsRecorder *prometheus.MetricsRecorder
}

func NewCreateHandler(
	shortURLCreator *service.URLShortener,
	logger *logger.Logger,
	metricsRecorder *prometheus.MetricsRecorder) *CreateHandler {
	return &CreateHandler{
		shortURLCreator: shortURLCreator,
		logger:          logger,
		metricsRecorder: metricsRecorder,
	}
}

func (h *CreateHandler) Create(ctx *fasthttp.RequestCtx) {
	var req service.Request
	h.metricsRecorder.RecordRequest(metrics.EventTypeCreate)

	err := json.Unmarshal(ctx.Request.Body(), &req)
	if err != nil {
		h.RespondBadRequest(ctx)
		h.metricsRecorder.RecordResponse(metrics.StatusBadRequest)

		return
	}

	response, err := h.shortURLCreator.Create(&req)
	if err != nil {
		h.RespondInternalError(ctx)
		h.metricsRecorder.RecordResponse(metrics.StatusInternalError)

		return
	}

	responseBody, _ := json.Marshal(response)

	h.RespondOK(ctx, responseBody)
	h.metricsRecorder.RecordResponse(metrics.StatusOk)
}
