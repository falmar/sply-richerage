package http

import (
	"context"
	"encoding/json"
	"github.com/falmar/richerage-api/internal/bootstrap"
	"github.com/falmar/richerage-api/internal/pkg/zaplogger"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"github.com/spf13/viper"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttp_Tickers_Method(t *testing.T) {
	methods := []string{
		"POST",
		"PUT",
		"DELETE",
		"PATCH",
		"HEAD",
		"OPTIONS",
		"TRACE",
		"CONNECT",
		"WAT",
	}

	ctx := context.Background()

	// bootstrap config
	v := viper.New()
	v.Set("port", "8080")
	logger := zaplogger.New(true)

	config, err := bootstrap.New(ctx, v, logger)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
	}

	handler, err := Handler(ctx, config)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	req := &http.Request{
		URL: &url.URL{
			Scheme: "http",
			Host:   server.Listener.Addr().String(),
			Path:   "/tickers",
		},
		Header: http.Header{},
		Body:   nil,
	}

	for _, method := range methods {
		req.Method = method

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("unexpected error to be nil, got: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusMethodNotAllowed {
			t.Errorf("unexpected status code to be %d, got: %d", http.StatusMethodNotAllowed, resp.StatusCode)
		}
	}
}

func TestHttp_Tickers_NoAuth(t *testing.T) {
	tokens := []string{
		"",
		"W2uBwQJSdSbsme9gIJ9lEeXl++PjDrOiEt90qSc=.dGVzMnQ=.1708107850",
		"W2uBwQJSdSbsm.e9gIJ9lEeXl+1708107850",
	}

	ctx := context.Background()

	// bootstrap config
	v := viper.New()
	v.Set("port", "8080")
	logger := zaplogger.New(true)

	config, err := bootstrap.New(ctx, v, logger)
	if err != nil {
		t.Fatalf("unexpected error to be nil, got: %v", err)
	}

	handler, err := Handler(ctx, config)
	if err != nil {
		t.Fatalf("unexpected error to be nil, got: %v", err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   server.Listener.Addr().String(),
			Path:   "/tickers",
		},
		Header: http.Header{},
		Body:   nil,
	}

	for _, token := range tokens {
		req.SetBasicAuth(token, "")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("unexpected error to be nil, got: %v", err)
		}
		resp.Body.Close()

		if resp.StatusCode != http.StatusUnauthorized {
			t.Fatalf("expected StatusUnauthorized, got %v", resp.Status)
		}
	}
}

func TestHttp_Tickers_Auth(t *testing.T) {
	// different tokens means different output
	tokens := []string{
		"6YR6GMnnrpzr/V5vw3/j+Z/n78sNNWOoAXcgsIpEur8=.dGVzdA==.1708107377",
		"W2uBwQJSdSb2nuJmme9gIJ9lEeXl++PjDrOiEt90qSc=.dGVzMnQ=.1708107850",
		"XoseHrqVheDWBNiZkpGt4BP8RkkieV6TnL7nNzNheUw=.MXRlczJ0.1708107864",
		"7XEYg1NflpD0fGxugiQcUzL4aQnwWxagyPidaUWNCpc=.MXRlZnMydA==.1708107874",
	}

	ctx := context.Background()

	// bootstrap config
	v := viper.New()
	v.Set("port", "8080")
	logger := zaplogger.New(true)

	config, err := bootstrap.New(ctx, v, logger)
	if err != nil {
		t.Fatalf("unexpected error to be nil, got: %v", err)
	}

	handler, err := Handler(ctx, config)
	if err != nil {
		t.Fatalf("unexpected error to be nil, got: %v", err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   server.Listener.Addr().String(),
			Path:   "/tickers",
		},
		Header: http.Header{},
		Body:   nil,
	}

	var prevTickers []types.Ticker

	for i, token := range tokens {
		req.SetBasicAuth(token, "")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("unexpected error to be nil, got: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			t.Fatalf("expected StatusOK, got %v", resp.Status)
		}

		var tickers []types.Ticker
		err = json.NewDecoder(resp.Body).Decode(&tickers)
		if err != nil {
			resp.Body.Close()
			t.Fatal(err)
		}

		resp.Body.Close()

		if len(tickers) < 1 {
			t.Fatalf("expected tickers to be not empty")
		} else if len(tickers) > 10 {
			t.Fatalf("expected tickers to be less than 10, got %v", len(tickers))
		}

		// Check the response to see if it's not empty and the price is above zero
		for _, ticker := range tickers {
			if ticker.Symbol == "" {
				t.Errorf("expected symbol to be not empty, got empty string")
			}
			if ticker.Price <= 0 {
				t.Errorf("expected price to be above 0, got %v", ticker.Price)
			}
		}

		if i == 0 {
			continue
		}

		if len(prevTickers) != len(tickers) {
			prevTickers = tickers
			continue
		}

		for z, ticker := range tickers {
			if ticker.Symbol == prevTickers[z].Symbol {
				t.Errorf("expected symbol to be different, got %s", ticker.Symbol)
			}
			if ticker.Price == prevTickers[z].Price {
				t.Errorf("expected price to be different, got %v", ticker.Price)
			}
		}

		prevTickers = tickers
	}
}

func TestHttp_Tickers_Constant(t *testing.T) {
	ctx := context.Background()

	// bootstrap config
	v := viper.New()
	v.Set("port", "8080")
	logger := zaplogger.New(true)

	config, err := bootstrap.New(ctx, v, logger)
	if err != nil {
		t.Fatalf("unexpected error to be nil, got: %v", err)
	}

	handler, err := Handler(ctx, config)
	if err != nil {
		t.Fatalf("unexpected error to be nil, got: %v", err)
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	req := &http.Request{
		Method: "GET",
		URL: &url.URL{
			Scheme: "http",
			Host:   server.Listener.Addr().String(),
			Path:   "/tickers",
		},
		Header: http.Header{},
		Body:   nil,
	}

	req.SetBasicAuth("7XEYg1NflpD0fGxugiQcUzL4aQnwWxagyPidaUWNCpc=.MXRlZnMydA==.1708107874", "")

	var prevTickers []types.Ticker

	// given token should be same output
	for i := 0; i < 10; i++ {
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("unexpected error to be nil, got: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			t.Fatalf("expected StatusOK, got %v", resp.Status)
		}

		var tickers []types.Ticker
		err = json.NewDecoder(resp.Body).Decode(&tickers)
		if err != nil {
			resp.Body.Close()
			t.Fatal(err)
		}

		resp.Body.Close()

		if i == 0 {
			prevTickers = tickers
		}

		// check if the response is the same as the previous one
		for z, ticker := range tickers {
			if ticker.Symbol != prevTickers[z].Symbol {
				t.Errorf("expected symbol to be %s, got %s", prevTickers[z].Symbol, ticker.Symbol)
			}
			if ticker.Price != prevTickers[z].Price {
				t.Errorf("expected price to be %v, got %v", prevTickers[z].Price, ticker.Price)
			}
		}
	}
}
