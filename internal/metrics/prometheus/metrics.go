package prometheus

import (
	"url-shortener/internal/metrics"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/collectors"
)

const (
	MetricResponse = "response_total"
	MetricRequest  = "request_total"
)

type MetricsRecorder struct {
	Registry *prometheus.Registry
	request  *prometheus.CounterVec
	response *prometheus.CounterVec
}

type MetricsConfig struct {
	Namespace string
	Subsystem string
}

const (
	LabelRequestType = "request_type"
)

func NewMetricsRecorder(cfg MetricsConfig) *MetricsRecorder {
	mtx := MetricsRecorder{}
	mtx.Registry = prometheus.NewRegistry()

	labelRequestType := []string{LabelRequestType}
	labelResponseType := []string{LabelRequestType}

	mtx.request = newCounter(
		cfg, MetricRequest, "The url-shortener cumulative request total counter.", labelRequestType)

	mtx.response = newCounter(
		cfg, MetricResponse, "The url-shortener cumulative response total counter.", labelResponseType)

	mtx.Registry.MustRegister(
		collectors.NewGoCollector(),
		collectors.NewProcessCollector(collectors.ProcessCollectorOpts{}),
		mtx.response,
		mtx.request,
	)

	return &mtx
}

func newSimpleCounter(cfg MetricsConfig, name string, help string) prometheus.Counter {
	opts := prometheus.CounterOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      name,
		Help:      help,
	}

	return prometheus.NewCounter(opts)
}

func newCounter(cfg MetricsConfig, name string, help string, labels []string) *prometheus.CounterVec {
	opts := prometheus.CounterOpts{
		Namespace: cfg.Namespace,
		Subsystem: cfg.Subsystem,
		Name:      name,
		Help:      help,
	}

	return prometheus.NewCounterVec(opts, labels)
}

func (m *MetricsRecorder) RecordRequest(reqType metrics.EventType) {
	m.request.WithLabelValues(string(reqType)).Inc()
}

func (m *MetricsRecorder) RecordResponse(resType metrics.ResponseType) {
	m.response.WithLabelValues(string(resType)).Inc()
}
