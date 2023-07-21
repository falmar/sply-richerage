package tickers

import (
	"context"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"time"
)

type GetTickerHistoryInput struct {
	Symbol string
	Before time.Time
}

type GetTickerHistoryOutput struct {
	History []types.TickerHistory
}

func (s *service) GetTickerHistory(ctx context.Context, in *GetTickerHistoryInput) (*GetTickerHistoryOutput, error) {
	history, err := s.storage.GetHistory(ctx, in.Symbol, in.Before)
	if err != nil {
		return nil, err
	}

	sortHistory(history)

	return &GetTickerHistoryOutput{
		History: history,
	}, nil
}

func sortHistory(history []types.TickerHistory) {
	for i := 0; i < len(history); i++ {
		for j := 0; j < len(history)-i-1; j++ {
			if history[j].Date.Before(history[j+1].Date) {
				history[j], history[j+1] = history[j+1], history[j]
			}
		}
	}
}
