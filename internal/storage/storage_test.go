//go:build test

package storage

import (
	"github.com/falmar/richerage-api/internal/tickers/types"
	"testing"
)

func TestStorage_ValidSymbol(t *testing.T) {
	v := isValidSymbol(types.ValidTickers(), "AAPL")
	if !v {
		t.Errorf("expected true, got %v", v)
	}

	v = isValidSymbol(types.ValidTickers(), "INVALID")
	if v {
		t.Errorf("expected false, got %v", v)
	}

	v = isValidSymbol([]string{}, "AAPL")
	if v {
		t.Errorf("expected false, got %v", v)
	}
}
