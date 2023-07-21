//go:build test

package tickers

import "testing"

func TestTicker_New(t *testing.T) {
	_, err := New(nil)
	if err != ErrInvalidConfig {
		t.Errorf("expected %v, got %v", ErrInvalidConfig, err)
	}

	svc, err := New(&Config{})
	if svc == nil {
		t.Error("service is nil")
	}
}
