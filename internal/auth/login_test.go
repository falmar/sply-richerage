package auth

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/pkg/hasher"
	"testing"
)

func TestAuth_Login(t *testing.T) {
	ctx := context.Background()
	mock := hasher.NewMock()

	mock.(*hasher.MockHasher).GenerateTokenFunc = func(ctx context.Context, data []byte) ([]byte, error) {
		if string(data) != "test" {
			t.Errorf("expected test, got %s", data)
		}

		return []byte("token"), nil
	}

	svc, err := New(&Config{
		Hasher: mock,
	})
	if err != nil {
		t.Errorf("unexpected error to be nil, got %v", err)
		return
	}

	out, err := svc.Login(ctx, &LoginInput{
		Username: "test",
	})
	if err != nil {
		t.Errorf("unexpected error to be nil, got %v", err)
		return
	} else if out == nil {
		t.Errorf("expected output to be set, got nil")
		return
	}

	if out.Token == "" {
		t.Errorf("expected token to be set, got empty string")
	} else if out.Token != "token" {
		t.Errorf("expected token to be token, got %s", out.Token)
	}
}

func TestAuth_Login_Error(t *testing.T) {
	ctx := context.Background()
	mock := hasher.NewMock()

	hasherErr := errors.New("hasher error")

	mock.(*hasher.MockHasher).GenerateTokenFunc = func(ctx context.Context, data []byte) ([]byte, error) {
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
	_, err = svc.Login(ctx, &LoginInput{
		Username: "test",
	})
	if err != hasherErr {
		t.Errorf("expected error to be %v, got %v", hasherErr, err)
		return
	}
}
