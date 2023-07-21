package kit

import (
	"context"
	"encoding/json"
	"go.uber.org/zap"
	"net/http"
)

type ErrorHandler struct {
	Logger *zap.Logger
}

func (th *ErrorHandler) Handle(_ context.Context, err error) {
	if cErr, ok := err.(CodedError); ok {
		th.Logger.Warn(
			"http request: handled error ",
			zap.String("message", cErr.Error()),
			zap.String("code", cErr.Code()),
			zap.Int("status", getStatusCode(err)),
		)

		return
	}

	th.Logger.Error("http request: unexpected error", zap.Error(err))
}

func (th *ErrorHandler) ErrorEncoder(_ context.Context, err error, w http.ResponseWriter) {
	var statusCode = 500
	body := map[string]interface{}{
		"code":    "internal_server_error",
		"message": "internal server error",
	}

	if cErr, ok := err.(CodedError); ok {
		statusCode = getStatusCode(err)
		body["message"] = cErr.Error()
		body["code"] = cErr.Code()
	}
	if bErr, ok := err.(*BadRequestError); ok && len(bErr.Params) > 0 {
		body["params"] = bErr.Params
	}

	w.Header().Set("content-type", "application/json")

	b, err := json.Marshal(body)
	if err != nil {
		w.WriteHeader(500)
		_, _ = w.Write([]byte(`{"code":"internal_server_error","message":"internal server error"}`))
		return
	}

	w.WriteHeader(statusCode)
	_, _ = w.Write(b)
}

func getStatusCode(err error) int {
	if v, ok := err.(HttpError); ok {
		return v.HttpCode()
	}

	return 500
}
