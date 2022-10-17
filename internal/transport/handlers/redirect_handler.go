package handlers

import (
	"net/http"
	"url-shortener/internal/logger"
	"url-shortener/internal/metrics"
	"url-shortener/internal/metrics/prometheus"
	"url-shortener/internal/service"

	"github.com/valyala/fasthttp"
)

type IRedirectService interface {
	Redirect(shortURL string) string
}

type RedirectHandler struct {
	baseHandler
	redirectService *service.RedirectService
	logger          *logger.Logger
	metricsRecorder *prometheus.MetricsRecorder
}

func NewRedirectHandler(
	redirectService *service.RedirectService,
	logger *logger.Logger,
	metricsRecorder *prometheus.MetricsRecorder) *RedirectHandler {
	return &RedirectHandler{
		redirectService: redirectService,
		logger:          logger,
		metricsRecorder: metricsRecorder,
	}
}

func (h *RedirectHandler) Redirect(ctx *fasthttp.RequestCtx) {
	h.metricsRecorder.RecordRequest(metrics.EventTypeRedirect)

	shortURL := ctx.UserValue("hash").(string)
	url := h.redirectService.Redirect(shortURL)

	ctx.Redirect(url, http.StatusFound)
	h.metricsRecorder.RecordResponse(metrics.StatusOk)
}
