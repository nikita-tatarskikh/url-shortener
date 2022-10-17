package transport

import (
	"url-shortener/internal/transport/handlers"

	"github.com/fasthttp/router"
	"github.com/valyala/fasthttp"
)

type FastHTTPHandlers struct {
	CreateHandler   *handlers.CreateHandler
	RedirectHandler *handlers.RedirectHandler
}

func NewFastHTTPHandlers(createHandler *handlers.CreateHandler, redirectHandler *handlers.RedirectHandler) *FastHTTPHandlers {
	return &FastHTTPHandlers{
		CreateHandler:   createHandler,
		RedirectHandler: redirectHandler,
	}
}

func NewFastHTTPRouter(h *FastHTTPHandlers) fasthttp.RequestHandler {

	r := router.New()

	r.POST("/create", h.CreateHandler.Create)
	r.GET("/{hash}", h.RedirectHandler.Redirect)

	return r.Handler
}
