package handlers

import (
	"net/http"

	"github.com/valyala/fasthttp"
)

const jsonContentType = "application/json"

type baseHandler struct {
}

func (h *baseHandler) RespondBadRequest(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusBadRequest)
}

func (h *baseHandler) RespondInternalError(ctx *fasthttp.RequestCtx) {
	ctx.SetStatusCode(http.StatusInternalServerError)
}

func (h *baseHandler) RespondOK(ctx *fasthttp.RequestCtx, responseBody []byte) {
	ctx.SetStatusCode(http.StatusOK)
	ctx.SetContentType(jsonContentType)
	_, _ = ctx.Write(responseBody)

}
