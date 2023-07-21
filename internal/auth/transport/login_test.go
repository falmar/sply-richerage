package transport

import (
	"context"
	"github.com/falmar/richerage-api/internal/auth/endpoint"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestLogin_RequestDecoder(t *testing.T) {
	testBody := `{"username": "test", "password": "12345"}`
	r, err := http.NewRequest("POST", "/login", strings.NewReader(testBody))
	if err != nil {
		t.Fatal(err)
	}

	out, err := LoginRequestDecoder(context.Background(), r)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	req, ok := out.(*endpoint.LoginRequest)
	if !ok || req == nil {
		t.Errorf("expected request to be of type LoginRequest, got %T", out)
		return
	}

	if req.Username != "test" {
		t.Errorf("expected username to be test, got %s", req.Username)
	}
	if req.Password != "12345" {
		t.Errorf("expected password to be 12345, got %s", req.Password)
	}
}

func TestLogin_RequestDecoder_Empty(t *testing.T) {
	testBody := ``
	r, err := http.NewRequest("POST", "/login", strings.NewReader(testBody))
	if err != nil {
		t.Fatal(err)
	}

	out, err := LoginRequestDecoder(context.Background(), r)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	req, ok := out.(*endpoint.LoginRequest)
	if !ok || req == nil {
		t.Errorf("expected request to be of type LoginRequest, got %T", out)
		return
	}

	if req.Username != "" {
		t.Errorf("expected username to be empty, got %s", req.Username)
	}
	if req.Password != "" {
		t.Errorf("expected password to be empty, got %s", req.Password)
	}
}

func TestLogin_ResponseEncoder(t *testing.T) {
	w := httptest.NewRecorder()

	resp := &endpoint.LoginResponse{
		Token: "sometoken",
	}

	err := LoginResponseEncoder(context.Background(), w, resp)
	if err != nil {
		t.Error("expected error to be nil, got", err)
	}

	got := w.Body.String()
	expect := `{"token":"sometoken"}` + "\n"

	if got != expect {
		t.Errorf("got %v, want %v", got, expect)
	}
}
