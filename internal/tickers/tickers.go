package tickers

import (
	"context"
	"github.com/falmar/richerage-api/internal/tickers/types"
)

type GetTickersInput struct {
	Username string
}

type GetTickersOutput struct {
	Tickers []types.Ticker
}

func (s *service) GetTickers(ctx context.Context, in *GetTickersInput) (*GetTickersOutput, error) {
	tickers, err := s.storage.GetByUser(ctx, in.Username)
	if err != nil {
		return nil, err
	}

	return &GetTickersOutput{
		Tickers: tickers,
	}, nil
}
