package auth

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/pkg/hasher"
	"testing"
)

func TestAuth_VerifyToken(t *testing.T) {
	ctx := context.Background()
	mock := hasher.NewMock()

	mock.(*hasher.MockHasher).ValidateTokenFunc = func(ctx context.Context, token []byte) ([]byte, error) {
		if string(token) != "token" {
			t.Errorf("expected token, got %s", token)
		}

		return []byte("test"), nil
	}

	svc, err := New(&Config{
		Hasher: mock,
	})
	if err != nil {
		t.Errorf("unexpected error to be nil, got %v", err)
		return
	}

	out, err := svc.VerifyToken(ctx, &VerifyTokenInput{
		Token: "token",
	})
	if err != nil {
		t.Errorf("unexpected error to be nil, got %v", err)
		return
	} else if out == nil {
		t.Errorf("expected output to be set, got nil")
		return
	}

	if out.Username == "" {
		t.Errorf("expected username to be set, got empty string")
	} else if out.Username != "test" {
		t.Errorf("expected username to be test, got %s", out.Username)
	}
}

func TestAuth_VerifyToken_Error(t *testing.T) {
	ctx := context.Background()
	mock := hasher.NewMock()

	hasherErr := errors.New("hasher error")

	mock.(*hasher.MockHasher).ValidateTokenFunc = func(ctx context.Context, token []byte) ([]byte, error) {
		return nil, hasherErr
	}

	svc, err := New(&Config{
		Hasher: mock,
	})
	if err != nil {
		t.Errorf("unexpected error to be nil, got %v", err)
		return
	}

	// should return error from hasher
	_, err = svc.VerifyToken(ctx, &VerifyTokenInput{
		Token: "token",
	})
	if err != hasherErr {
		t.Errorf("expected error to be %v, got %v", hasherErr, err)
		return
	}
}
