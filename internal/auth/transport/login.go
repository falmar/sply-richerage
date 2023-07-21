package transport

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/falmar/richerage-api/internal/auth/endpoint"
	"io"
	"net/http"
)

func LoginRequestDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	req := &endpoint.LoginRequest{}

	// let LoginEndpoint handle the validation of empty body
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil && !errors.Is(err, io.EOF) {
		return nil, err
	}

	return req, nil
}

func LoginResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	resp := response.(*endpoint.LoginResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(resp)
}
