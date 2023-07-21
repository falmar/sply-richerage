package transport

import (
	"context"
	"encoding/json"
	"github.com/falmar/richerage-api/internal/tickers/endpoint"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTickers_RequestDecoder(t *testing.T) {
	r, err := http.NewRequest("GET", "/tickers", nil)
	if err != nil {
		t.Fatal(err)
	}

	out, err := TickersRequestDecoder(context.Background(), r)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	_, ok := out.(*endpoint.TickersRequest)
	if !ok {
		t.Errorf("expected request to be of type TickersRequest, got %T", out)
	}
}

func TestTickers_ResponseEncoder(t *testing.T) {
	w := httptest.NewRecorder()

	tickers := []types.Ticker{
		{Symbol: "BTC", Price: 45000.0},
		{Symbol: "ETH", Price: 3000.0},
	}
	resp := &endpoint.TickersResponse{
		Tickers: tickers,
	}

	err := TickersResponseEncoder(context.Background(), w, resp)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	var got []types.Ticker
	err = json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Error("expected error to be nil, got", err)
		return
	}

	if len(got) != len(tickers) {
		t.Errorf("got length %d, want length %d", len(got), len(tickers))
		return
	}

	for i, ticker := range tickers {
		if got[i].Symbol != ticker.Symbol || got[i].Price != ticker.Price {
			t.Errorf("got ticker %+v, want ticker %+v", got[i], ticker)
		}
	}

	tickers = []types.Ticker{}
	resp = &endpoint.TickersResponse{
		Tickers: tickers,
	}

	w = httptest.NewRecorder()
	err = TickersResponseEncoder(context.Background(), w, resp)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	got = []types.Ticker{}
	err = json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Error("expected error to be nil, got", err)
		return
	}

	if len(got) != len(tickers) {
		t.Errorf("got length %d, want length %d", len(got), len(tickers))
		return
	}
}

func TestTickers_ResponseEncoder_Empty(t *testing.T) {
	w := httptest.NewRecorder()

	tickers := []types.Ticker{}
	resp := &endpoint.TickersResponse{
		Tickers: tickers,
	}

	err := TickersResponseEncoder(context.Background(), w, resp)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	var got []types.Ticker
	err = json.NewDecoder(w.Body).Decode(&got)
	if err != nil {
		t.Error("expected error to be nil, got", err)
		return
	}

	if len(got) != len(tickers) {
		t.Errorf("got length %d, want length %d", len(got), len(tickers))
		return
	}
}
