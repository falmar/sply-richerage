//go:build test

package tickers

import (
	"context"
	"github.com/falmar/richerage-api/internal/storage"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"testing"
)

func TestTickers_Tickers(t *testing.T) {
	ctx := context.Background()
	st := storage.NewMock()

	st.(*storage.MockStorage).GetByUserFunc = func(ctx context.Context, username string) ([]types.Ticker, error) {
		if username != "test" {
			t.Errorf("expected test, got %s", username)
		}

		return []types.Ticker{
			{
				Price:  100,
				Symbol: "AAPL",
			},
		}, nil
	}

	svc, err := New(&Config{
		Storage: st,
	})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	out, err := svc.GetTickers(ctx, &GetTickersInput{
		Username: "test",
	})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if out == nil {
		t.Errorf("expected tickers to be set, got nil")
		return
	}

	if out.Tickers == nil || len(out.Tickers) != 1 {
		t.Errorf("expected 1 ticker, got %d", len(out.Tickers))
		return
	}

	if out.Tickers[0].Symbol != "AAPL" {
		t.Errorf("expected ticker symbol to be AAPL, got %s", out.Tickers[0].Symbol)
	} else if out.Tickers[0].Price != 100 {
		t.Errorf("expected ticker price to be 100, got %f", out.Tickers[0].Price)
	}
}

func TestTickers_Tickers_Error(t *testing.T) {
	// return error raised by storage

	ctx := context.Background()

	stErr := storage.ErrMockUncalledFor

	st := storage.NewMock()
	st.(*storage.MockStorage).GetByUserFunc = func(ctx context.Context, username string) ([]types.Ticker, error) {
		return nil, stErr
	}

	svc, err := New(&Config{
		Storage: st,
	})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	_, err = svc.GetTickers(ctx, &GetTickersInput{
		Username: "test",
	})
	if err != stErr {
		t.Errorf("expected error to be %T, got %T", stErr, err)
	}
}
