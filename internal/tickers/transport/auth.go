package transport

import (
	"context"
	"net/http"
)

func TokenDecoder(ctx context.Context, r *http.Request) context.Context {
	// let endpoint handle auth checks
	// transport is only responsible for decoding it from request
	u, _, _ := r.BasicAuth()

	return context.WithValue(ctx, "auth_token", u)
}
