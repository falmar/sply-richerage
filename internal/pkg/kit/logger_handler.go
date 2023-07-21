package kit

import (
	"context"
	"go.uber.org/zap"
	"net/http"
)

type LoggerHandler struct {
	Logger *zap.Logger
}

func (h *LoggerHandler) Before(ctx context.Context, req *http.Request) context.Context {
	reqID := req.Header.Get("x-request-id")
	logger := h.Logger.
		With(zap.String("request_id", reqID)).
		With(zap.String("path", req.URL.Path)).
		With(zap.String("method", req.Method))

	ctx = context.WithValue(ctx, "request_id", reqID)

	logger.Info("http: request received")

	return ctx
}
func (h *LoggerHandler) After(ctx context.Context, _ http.ResponseWriter) context.Context {
	h.Logger.Info("http: request processed")

	return ctx
}
