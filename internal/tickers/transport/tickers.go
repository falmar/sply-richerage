package transport

import (
	"context"
	"encoding/json"
	"github.com/falmar/richerage-api/internal/tickers/endpoint"
	"net/http"
)

func TickersRequestDecoder(context.Context, *http.Request) (interface{}, error) {
	return &endpoint.TickersRequest{}, nil
}

func TickersResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(*endpoint.TickersResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	return json.NewEncoder(w).Encode(res.Tickers)
}
