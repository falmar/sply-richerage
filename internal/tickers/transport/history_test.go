package transport

import (
	"context"
	"encoding/json"
	"github.com/falmar/richerage-api/internal/tickers/endpoint"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"github.com/go-chi/chi/v5"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestTickerHistory_RequestDecoder(t *testing.T) {
	r, _ := http.NewRequest("GET", "/ticker/BTC/history", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("symbol", "BTC")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	out, err := TickerHistoryRequestDecoder(context.Background(), r)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	req, ok := out.(*endpoint.TickerHistoryRequest)
	if !ok || req == nil {
		t.Errorf("expected request to be of type TickerHistoryRequest, got %T", out)
		return
	}

	if req.Symbol != "BTC" {
		t.Errorf("expected symbol to be BTC, got %s", req.Symbol)
	}
}

func TestTickerHistory_RequestDecoder_Empty(t *testing.T) {
	r, _ := http.NewRequest("GET", "/ticker//history", nil)

	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("symbol", "")
	r = r.WithContext(context.WithValue(r.Context(), chi.RouteCtxKey, rctx))

	out, err := TickerHistoryRequestDecoder(context.Background(), r)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	req, ok := out.(*endpoint.TickerHistoryRequest)
	if !ok || req == nil {
		t.Errorf("expected request to be of type TickerHistoryRequest, got %T", out)
		return
	}

	if req.Symbol != "" {
		t.Errorf("expected symbol to be empty, got %s", req.Symbol)
	}
}
func TestTickerHistory_ResponseEncoder(t *testing.T) {
	w := httptest.NewRecorder()

	tickers := []types.TickerHistory{
		{Price: 45000.0, Date: time.Date(2023, 07, 21, 0, 0, 0, 0, time.UTC)},
		{Price: 46000.0, Date: time.Date(2023, 07, 22, 0, 0, 0, 0, time.UTC)},
	}
	resp := &endpoint.TickerHistoryResponse{
		Tickers: tickers,
	}

	err := TickerHistoryResponseEncoder(context.Background(), w, resp)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	var got []map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	if len(got) != len(tickers) {
		t.Errorf("got length %d, want length %d", len(got), len(tickers))
	}

	// important to check date format
	for i, ticker := range tickers {
		if got[i]["price"].(float64) != ticker.Price || got[i]["date"].(string) != ticker.Date.Format("2006-01-02") {
			t.Errorf("got ticker %+v, want ticker %+v", got[i], ticker)
		}
	}
}

func TestTickerHistory_ResponseEncoder_Empty(t *testing.T) {
	w := httptest.NewRecorder()

	resp := &endpoint.TickerHistoryResponse{
		Tickers: []types.TickerHistory{},
	}

	err := TickerHistoryResponseEncoder(context.Background(), w, resp)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	var got []map[string]interface{}
	err = json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	if len(got) != len(resp.Tickers) {
		t.Errorf("got length %d, want length %d", len(got), len(resp.Tickers))
	}
}
