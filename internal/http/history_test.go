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
	"time"
)

func TestHttp_History_Method(t *testing.T) {
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
			Path:   "/tickers/AMZN/history",
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

func TestHttp_History_NoAuth(t *testing.T) {
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
			Path:   "/tickers/AMZN/history",
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

func TestHttp_History_Auth_SymbolDiff(t *testing.T) {
	// same token but different symbols means different output
	matrix := [][]string{
		{"Y04JuAE90RdDLoxdy93+5nuKDnyJr89UdpgV1I3pW+U=.MXRlZnMydA==.1708108610", "AMZN"},
		{"Y04JuAE90RdDLoxdy93+5nuKDnyJr89UdpgV1I3pW+U=.MXRlZnMydA==.1708108610", "AAPL"},
		{"Y04JuAE90RdDLoxdy93+5nuKDnyJr89UdpgV1I3pW+U=.MXRlZnMydA==.1708108610", "MSFT"},
		{"Y04JuAE90RdDLoxdy93+5nuKDnyJr89UdpgV1I3pW+U=.MXRlZnMydA==.1708108610", "GOOG"},
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
		},
		Header: http.Header{},
		Body:   nil,
	}

	var prevHistory []types.TickerHistory

	for i, m := range matrix {
		req.URL.Path = "/tickers/" + m[1] + "/history"
		req.SetBasicAuth(m[0], "")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("unexpected error to be nil, got: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			t.Fatalf("expected StatusOK, got %v", resp.Status)
		}

		var JSONHistory []struct {
			Date  string  `json:"date"`
			Price float64 `json:"price"`
		}
		err = json.NewDecoder(resp.Body).Decode(&JSONHistory)
		if err != nil {
			resp.Body.Close()
			t.Fatal(err)
		}
		resp.Body.Close()

		var history []types.TickerHistory
		for _, ts := range JSONHistory {
			date, err := time.Parse("2006-01-02", ts.Date)
			if err != nil {
				t.Fatal(err)
			}

			history = append(history, types.TickerHistory{
				Date:  date,
				Price: ts.Price,
			})
		}

		if len(history) < 1 {
			t.Fatalf("expected tickers to be not empty")
		} else if len(history) > 100 {
			t.Fatalf("expected tickers to be less than 10, got %v", len(history))
		}

		// Check the response to see if it's not empty and the price is above zero
		for _, ts := range history {
			if ts.Date.IsZero() {
				t.Errorf("expected date to be not zero, got %v", ts.Date)
			}
			if ts.Price <= 0 {
				t.Errorf("expected price to be above 0, got %v", ts.Price)
			}
		}

		if i == 0 || len(prevHistory) != len(history) {
			prevHistory = history
			continue
		}

		for z, ts := range history {
			// date may be the same since its not changing much

			if ts.Price == prevHistory[z].Price {
				t.Errorf("expected price to be different, got %v", ts.Price)
			}
		}

		prevHistory = history
	}
}

func TestHttp_History_Auth_SymbolSame(t *testing.T) {
	// different tokens but same symbols means same output
	matrix := [][]string{
		{"6YR6GMnnrpzr/V5vw3/j+Z/n78sNNWOoAXcgsIpEur8=.dGVzdA==.1708107377", "AMZN"},
		{"W2uBwQJSdSb2nuJmme9gIJ9lEeXl++PjDrOiEt90qSc=.dGVzMnQ=.1708107850", "AMZN"},
		{"XoseHrqVheDWBNiZkpGt4BP8RkkieV6TnL7nNzNheUw=.MXRlczJ0.1708107864", "AMZN"},
		{"7XEYg1NflpD0fGxugiQcUzL4aQnwWxagyPidaUWNCpc=.MXRlZnMydA==.1708107874", "AMZN"},
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
		},
		Header: http.Header{},
		Body:   nil,
	}

	var prevHistory []types.TickerHistory

	for i, m := range matrix {
		req.URL.Path = "/tickers/" + m[1] + "/history"
		req.SetBasicAuth(m[0], "")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("unexpected error to be nil, got: %v", err)
		}

		if resp.StatusCode != http.StatusOK {
			resp.Body.Close()
			t.Fatalf("expected StatusOK, got %v", resp.Status)
		}

		var JSONHistory []struct {
			Date  string  `json:"date"`
			Price float64 `json:"price"`
		}
		err = json.NewDecoder(resp.Body).Decode(&JSONHistory)
		if err != nil {
			resp.Body.Close()
			t.Fatal(err)
		}

		resp.Body.Close()

		var history []types.TickerHistory
		for _, ts := range JSONHistory {
			date, err := time.Parse("2006-01-02", ts.Date)
			if err != nil {
				t.Fatal(err)
			}

			history = append(history, types.TickerHistory{
				Date:  date,
				Price: ts.Price,
			})
		}

		if len(history) < 1 {
			t.Fatalf("expected tickers to be not empty")
		} else if len(history) > 100 {
			t.Fatalf("expected tickers to be less than 10, got %v", len(history))
		}

		// Check the response to see if it's not empty and the price is above zero
		for _, ts := range history {
			if ts.Date.IsZero() {
				t.Errorf("expected date to be not zero, got %v", ts.Date)
			}
			if ts.Price <= 0 {
				t.Errorf("expected price to be above 0, got %v", ts.Price)
			}
		}

		if i == 0 {
			prevHistory = history
			continue
		}

		for z, ts := range history {
			if ts.Date != prevHistory[z].Date {
				t.Errorf("expected date to be %v, got %v", prevHistory[z].Date, ts.Date)
			}
			if ts.Price != prevHistory[z].Price {
				t.Errorf("expected price to be %v, got %v", prevHistory[z].Price, ts.Price)
			}
		}

		history = prevHistory
	}
}

func TestHttp_History_Constant(t *testing.T) {
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
		},
		Header: http.Header{},
		Body:   nil,
	}

	req.URL.Path = "/tickers/AMZN/history"
	req.SetBasicAuth("6YR6GMnnrpzr/V5vw3/j+Z/n78sNNWOoAXcgsIpEur8=.dGVzdA==.1708107377", "")

	var prevHistory []types.TickerHistory

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

		var JSONHistory []struct {
			Date  string  `json:"date"`
			Price float64 `json:"price"`
		}
		err = json.NewDecoder(resp.Body).Decode(&JSONHistory)
		if err != nil {
			resp.Body.Close()
			t.Fatal(err)
		}
		resp.Body.Close()

		var history []types.TickerHistory
		for _, ts := range JSONHistory {
			date, err := time.Parse("2006-01-02", ts.Date)
			if err != nil {
				t.Fatal(err)
			}

			history = append(history, types.TickerHistory{
				Date:  date,
				Price: ts.Price,
			})
		}

		if i == 0 {
			prevHistory = history
		}

		// check if the response is the same as the previous one
		for z, ts := range history {
			if ts.Date != prevHistory[z].Date {
				t.Errorf("expected date to be %v, got %v", prevHistory[z].Date, ts.Date)
			}
			if ts.Price != prevHistory[z].Price {
				t.Errorf("expected price to be %v, got %v", prevHistory[z].Price, ts.Price)
			}
		}
	}
}
