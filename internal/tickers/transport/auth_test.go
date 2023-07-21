package transport

import (
	"context"
	"net/http"
	"testing"
)

func TestTokenDecoder(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("expected error to be nil, got", err)
		return
	}

	const username = "testuser"
	req.SetBasicAuth(username, "")

	ctx := TokenDecoder(context.Background(), req)
	token, ok := ctx.Value("auth_token").(string)

	if !ok {
		t.Error("expected token to be a string")
		return
	}

	if token != username {
		t.Errorf("expected token to be %q, got %q", username, token)
	}
}

func TestTokenDecoder_Empty(t *testing.T) {
	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Error("expected error to be nil, got", err)
		return
	}

	ctx := TokenDecoder(context.Background(), req)
	token, ok := ctx.Value("auth_token").(string)

	if !ok {
		t.Error("expected token to be a string")
		return
	}

	if token != "" {
		t.Errorf("expected token to be empty, got %q", token)
	}
}
