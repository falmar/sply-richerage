//go:build test

package storage

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"testing"
	"time"
)

func TestStorageSeeder_RandForString(t *testing.T) {
	rnd1 := getRandForString("test")
	rnd2 := getRandForString("test")

	if rnd1.Intn(100) != rnd2.Intn(100) {
		t.Errorf("expected same random number generator, got different")
	}

	rnd1 = getRandForString("test2")

	if rnd1.Intn(100) == rnd2.Intn(100) {
		t.Errorf("expected different random number generator, got same")
	}
}

func TestStorageSeeder_GenValidTickers(t *testing.T) {
	tickers := generatedTickers(
		getRandForString("test"),
	)
	tickers2 := generatedTickers(
		getRandForString("test"),
	)

	validTickers := types.ValidTickers()

	if len(tickers) != len(validTickers) {
		t.Errorf("expected %d tickers, got %d", len(validTickers), len(tickers))
		return
	}

	for _, ticker := range tickers {
		if !isValidSymbol(validTickers, ticker.Symbol) {
			t.Errorf("expected ticker %s to be valid", ticker.Symbol)
		}
	}

	for i := 0; i < len(validTickers); i++ {
		if tickers[i].Symbol != tickers2[i].Symbol {
			t.Errorf("expected ticker %s to be equal to %s", tickers[i].Symbol, tickers2[i].Symbol)
		} else if tickers[i].Price != tickers2[i].Price {
			t.Errorf("expected ticker %f to be equal to %f", tickers[i].Price, tickers2[i].Price)
		}
	}

	tickers2 = generatedTickers(
		getRandForString("test2"),
	)

	for i := 0; i < len(validTickers); i++ {
		if tickers[i].Symbol != tickers2[i].Symbol {
			t.Errorf("expected ticker %s to be equal to %s", tickers[i].Symbol, tickers2[i].Symbol)
		} else if tickers[i].Price == tickers2[i].Price {
			t.Errorf("expected ticker %f to be different to %f", tickers[i].Price, tickers2[i].Price)
		}
	}
}

func TestStorageSeeder_History(t *testing.T) {
	ctx := context.Background()

	// test GetHistory is deterministic on the same seeded_storage
	s1 := NewSeeded()

	h1, err := s1.GetHistory(ctx, "AAPL", time.Time{})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	h2, err := s1.GetHistory(ctx, "AAPL", time.Time{})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if len(h1) != len(h2) {
		t.Errorf("expected history to be equal, got different lengths")
		return
	}

	for i := 0; i < len(h1); i++ {
		if h1[i].Date != h2[i].Date {
			t.Errorf("expected history to be equal, got different dates")
		} else if h1[i].Price != h2[i].Price {
			t.Errorf("expected history to be equal, got different prices")
		}
	}

	// test GetHistory is deterministic on different seeded_storage
	s2 := NewSeeded()

	h3, err := s2.GetHistory(ctx, "AAPL", time.Time{})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if len(h1) != len(h3) {
		t.Errorf("expected history to be equal, got different lengths")
		return
	}

	for i := 0; i < len(h1); i++ {
		if h1[i].Date != h3[i].Date {
			t.Errorf("expected history to be equal, got different dates")
		} else if h1[i].Price != h3[i].Price {
			t.Errorf("expected history to be equal, got different prices")
		}
	}
}

func TestStorageSeeder_History_Invalid(t *testing.T) {
	s1 := NewSeeded()

	_, err := s1.GetHistory(context.Background(), "INVALID", time.Time{})

	var errNotFound *types.ErrTickerNotFound

	if !errors.As(err, &errNotFound) {
		t.Errorf("expected error to be ErrInvalidSymbol, got %T", err)
		return
	} else if errNotFound.Symbol != "INVALID" {
		t.Errorf("expected error to be ErrInvalidSymbol, got %s", errNotFound.Symbol)
	}
}

func TestStorageSeeder_GenTickers(t *testing.T) {
	ctx := context.Background()

	// test GetByUser is deterministic on the same seeded_storage
	s1 := NewSeeded()

	t1, err := s1.GetByUser(ctx, "test")
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	t2, err := s1.GetByUser(ctx, "test")
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if len(t1) != len(t2) {
		t.Errorf("expected tickers to be equal, got different lengths")
		return
	}

	for i := 0; i < len(t1); i++ {
		if t1[i].Symbol != t2[i].Symbol {
			t.Errorf("expected tickers to be equal, got different symbols")
		} else if t1[i].Price != t2[i].Price {
			t.Errorf("expected tickers to be equal, got different prices")
		}
	}

	// test GetByUser is deterministic on different seeded_storage
	s2 := NewSeeded()

	t3, err := s2.GetByUser(ctx, "test")
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if len(t1) != len(t3) {
		t.Errorf("expected tickers to be equal, got different lengths")
		return
	}

	for i := 0; i < len(t1); i++ {
		if t1[i].Symbol != t3[i].Symbol {
			t.Errorf("expected tickers to be equal, got different symbols")
		} else if t1[i].Price != t3[i].Price {
			t.Errorf("expected tickers to be equal, got different prices")
		}
	}

	t4, err := s2.GetByUser(ctx, "test2")
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if len(t1) == len(t4) {
		t.Errorf("expected tickers to be different, got same")
		return
	}
}
