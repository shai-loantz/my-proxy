package proxy

import (
	"context"
	"io"
	"net/http"
)

type proxyRequest struct {
	method         string
	URLPath        string
	body           []byte
	header         http.Header
	responseWriter http.ResponseWriter
	ctx            context.Context
	cancelReqCtx   context.CancelFunc
}

func NewProxyRequest(responseWriter http.ResponseWriter, req *http.Request, reqCtx context.Context, cancelReqCtx context.CancelFunc) proxyRequest {
	defer req.Body.Close()

	body, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(responseWriter, "Error reading request body", http.StatusInternalServerError)
		return proxyRequest{}
	}
	return proxyRequest{
		method:         req.Method,
		URLPath:        req.RequestURI,
		body:           body,
		header:         req.Header,
		responseWriter: responseWriter,
		ctx:            reqCtx,
		cancelReqCtx:   cancelReqCtx,
	}
}
