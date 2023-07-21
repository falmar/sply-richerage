//go:build test

package tickers

import (
	"context"
	"github.com/falmar/richerage-api/internal/storage"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"testing"
	"time"
)

func TestTickers_History_Sort(t *testing.T) {
	var h []types.TickerHistory
	dates := []time.Time{
		time.Now().Add(time.Hour * 24),
		time.Now().Add(time.Hour * 96),
		time.Now().Add(time.Hour * 72),
		time.Now().Add(time.Hour * 48),
	}
	expected := []time.Time{
		dates[1],
		dates[2],
		dates[3],
		dates[0],
	}
	for _, d := range dates {
		h = append(h, types.TickerHistory{
			Price: 0,
			Date:  d,
		})
	}

	sortHistory(h)

	for i, v := range h {
		if v.Date.IsZero() {
			t.Errorf("expected not zero date")
		}

		if v.Date != expected[i] {
			t.Errorf("expected %v, got %v", expected[i], v.Date)
		}
	}
}

func TestTickers_History(t *testing.T) {
	ctx := context.Background()
	st := storage.NewMock()

	date, err := time.Parse("2006-01-02", "2023-07-21")

	st.(*storage.MockStorage).GetHistoryFunc = func(ctx context.Context, symbol string, before time.Time) ([]types.TickerHistory, error) {
		if symbol != "AAPL" {
			t.Errorf("expected AAPL, got %s", symbol)
		}

		return []types.TickerHistory{
			{
				Price: 100,
				Date:  date,
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

	out, err := svc.GetTickerHistory(ctx, &GetTickerHistoryInput{
		Symbol: "AAPL",
	})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if out == nil {
		t.Errorf("expected tickers to be set, got nil")
		return
	}

	if out.History == nil || len(out.History) != 1 {
		t.Errorf("expected 1 ticker, got %d", len(out.History))
		return
	}

	if out.History[0].Date != date {
		t.Errorf("expected ticker date to be %v, got %v", date, out.History[0].Date)
	} else if out.History[0].Price != 100 {
		t.Errorf("expected ticker price to be 100, got %f", out.History[0].Price)
	}
}

func TestTickers_History_Error(t *testing.T) {
	// return error raised by storage

	ctx := context.Background()

	stErr := storage.ErrMockUncalledFor

	st := storage.NewMock()
	st.(*storage.MockStorage).GetHistoryFunc = func(ctx context.Context, symbol string, before time.Time) ([]types.TickerHistory, error) {
		return nil, stErr
	}

	svc, err := New(&Config{
		Storage: st,
	})
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	_, err = svc.GetTickerHistory(ctx, &GetTickerHistoryInput{
		Symbol: "test",
	})
	if err != stErr {
		t.Errorf("expected error to be %T, got %T", stErr, err)
	}
}
