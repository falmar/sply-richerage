package http

import (
	"bytes"
	"context"
	"encoding/json"
	"github.com/falmar/richerage-api/internal/bootstrap"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	"github.com/falmar/richerage-api/internal/pkg/zaplogger"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestHttp_Login_Method(t *testing.T) {
	methods := []string{
		"GET",
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
			Path:   "/login",
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: nil,
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

func TestHttp_Login_NoBody(t *testing.T) {
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
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   server.Listener.Addr().String(),
			Path:   "/login",
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: nil,
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Errorf("unexpected status code to be %d, got: %d", http.StatusBadRequest, resp.StatusCode)
	}

	respError := &kit.HttpErrorBody{}
	err = json.NewDecoder(resp.Body).Decode(respError)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
	}

	badRequest := kit.BadRequestError{}
	if respError.Code != badRequest.Code() {
		t.Errorf("unexpected error code to be %s, got: %s", badRequest.Code(), respError.Code)
	}
}

func TestHttp_Login_Token(t *testing.T) {
	ctx := context.Background()

	// bootstrap config
	v := viper.New()
	v.Set("port", "8080")
	logger := zaplogger.New(true)

	config, err := bootstrap.New(ctx, v, logger)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
		return
	}

	handler, err := Handler(ctx, config)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
		return
	}

	server := httptest.NewServer(handler)
	defer server.Close()

	req := &http.Request{
		Method: "POST",
		URL: &url.URL{
			Scheme: "http",
			Host:   server.Listener.Addr().String(),
			Path:   "/login",
		},
		Header: http.Header{
			"Content-Type": []string{"application/json"},
		},
		Body: io.NopCloser(bytes.NewBuffer([]byte(`{"username": "test", "password": "test"}`))),
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
		return
	}

	if resp.StatusCode != http.StatusOK {
		t.Errorf("unexpected status code to be %d, got: %d", http.StatusOK, resp.StatusCode)
		return
	}

	respBody := map[string]interface{}{}

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		t.Errorf("unexpected error to be nil, got: %v", err)
		return
	}

	if v, ok := respBody["token"]; !ok {
		t.Errorf("unexpected token to be present")
		return
	} else if v == "" {
		t.Errorf("unexpected token to be not empty")
		return
	}
}

// not testing credential validation here because its not part of the code auth provided
