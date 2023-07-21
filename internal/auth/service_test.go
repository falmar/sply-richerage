package auth

import (
	"github.com/falmar/richerage-api/internal/pkg/hasher"
	"testing"
)

func TestAuth_New(t *testing.T) {
	_, err := New(nil)
	if err != ErrInvalidConfig {
		t.Errorf("expected %v, got %v", ErrInvalidConfig, err)
	}

	svc, err := New(&Config{})
	if err != ErrInvalidConfig {
		t.Errorf("expected %v, got %v", ErrInvalidConfig, err)
	}

	svc, err = New(&Config{
		Hasher: hasher.NewMock(),
	})
	if svc == nil {
		t.Error("service is nil")
	}
}
