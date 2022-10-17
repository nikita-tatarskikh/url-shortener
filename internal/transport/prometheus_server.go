package transport

import (
	"context"
	"net/http"
	"time"
	"url-shortener/internal/logger"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

const DefaultMetricsTimeout = 10 * time.Second

type PrometheusMetricsServerConfig struct {
	Addr string
}

func NewPrometheusMetricsServer(cfg PrometheusMetricsServerConfig, logger *logger.Logger, gatherer prometheus.Gatherer) (server *http.Server, cleanup func()) {
	m := http.NewServeMux()
	m.Handle("/metrics", promhttp.HandlerFor(gatherer, promhttp.HandlerOpts{}))

	server = &http.Server{
		Addr:         cfg.Addr,
		Handler:      m,
		ReadTimeout:  DefaultMetricsTimeout,
		WriteTimeout: DefaultMetricsTimeout,
	}

	cleanup = func() {
		logger.LogInfo("prometheus shut down performed")

		if err := server.Shutdown(context.TODO()); err != nil {
			logger.LogError("error while shutdown", err)
		}
	}

	return server, cleanup
}
