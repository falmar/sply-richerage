package transport

import (
	"context"
	"encoding/json"
	"github.com/falmar/richerage-api/internal/tickers/endpoint"
	"github.com/go-chi/chi/v5"
	"net/http"
)

func TickerHistoryRequestDecoder(_ context.Context, r *http.Request) (interface{}, error) {
	req := &endpoint.TickerHistoryRequest{
		Symbol: chi.URLParam(r, "symbol"),
	}

	return req, nil
}

func TickerHistoryResponseEncoder(_ context.Context, w http.ResponseWriter, response interface{}) error {
	res := response.(*endpoint.TickerHistoryResponse)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	// format date at transport output
	var tickers []interface{}

	for _, ticker := range res.Tickers {
		tickers = append(tickers, map[string]interface{}{
			"price": ticker.Price,
			"date":  ticker.Date.Format("2006-01-02"),
		})
	}

	return json.NewEncoder(w).Encode(tickers)
}
