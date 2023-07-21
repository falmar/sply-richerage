//go:build test

package tickers

import (
	"context"
	"errors"
)

var _ Service = (*MockService)(nil)
var ErrMockUncalledFor = errors.New("uncalled for")

func NewMockService() Service {
	return &MockService{
		GetTickersFunc: func(ctx context.Context, in *GetTickersInput) (*GetTickersOutput, error) {
			return nil, ErrMockUncalledFor
		},
		GetTickerHistoryFunc: func(ctx context.Context, in *GetTickerHistoryInput) (*GetTickerHistoryOutput, error) {
			return nil, ErrMockUncalledFor
		},
	}
}

type MockService struct {
	GetTickersFunc       func(ctx context.Context, in *GetTickersInput) (*GetTickersOutput, error)
	GetTickerHistoryFunc func(ctx context.Context, in *GetTickerHistoryInput) (*GetTickerHistoryOutput, error)
}

func (m *MockService) GetTickers(ctx context.Context, in *GetTickersInput) (*GetTickersOutput, error) {
	return m.GetTickersFunc(ctx, in)
}

func (m *MockService) GetTickerHistory(ctx context.Context, in *GetTickerHistoryInput) (*GetTickerHistoryOutput, error) {
	return m.GetTickerHistoryFunc(ctx, in)
}
