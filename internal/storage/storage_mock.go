//go:build test

package storage

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"time"
)

var _ Storage = (*MockStorage)(nil)
var ErrMockUncalledFor = errors.New("uncalled for")

func NewMock() Storage {
	return &MockStorage{
		GetByUserFunc: func(ctx context.Context, username string) ([]types.Ticker, error) {
			return nil, ErrMockUncalledFor
		},
		GetHistoryFunc: func(ctx context.Context, symbol string, before time.Time) ([]types.TickerHistory, error) {
			return nil, ErrMockUncalledFor
		},
	}
}

type MockStorage struct {
	GetByUserFunc  func(ctx context.Context, username string) ([]types.Ticker, error)
	GetHistoryFunc func(ctx context.Context, symbol string, before time.Time) ([]types.TickerHistory, error)
}

func (m *MockStorage) GetByUser(ctx context.Context, username string) ([]types.Ticker, error) {
	return m.GetByUserFunc(ctx, username)
}

func (m *MockStorage) GetHistory(ctx context.Context, symbol string, before time.Time) ([]types.TickerHistory, error) {
	return m.GetHistoryFunc(ctx, symbol, before)
}
