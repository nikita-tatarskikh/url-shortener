package transport

import (
	"time"
	"url-shortener/internal/logger"

	"github.com/valyala/fasthttp"
)

const (
	DefaultReadBufferSize = 65536
	DefaultTimeout        = 10 * time.Second
)

func NewFastHTTPServer(handler fasthttp.RequestHandler, logger *logger.Logger) (server *fasthttp.Server, cleanup func()) {
	server = &fasthttp.Server{
		Handler:        handler,
		ReadBufferSize: DefaultReadBufferSize,
		ReadTimeout:    DefaultTimeout,
		WriteTimeout:   DefaultTimeout,
	}
	cleanup = func() {
		logger.LogInfo("shuts down gracefully")

		if err := server.Shutdown(); err != nil {
			logger.LogError("error while shutdown:", err)
		}
	}

	return
}
