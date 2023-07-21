package tickers

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/storage"
)

var ErrInvalidConfig = errors.New("invalid tickers service config")

var _ Service = (*service)(nil)

type Service interface {
	GetTickers(ctx context.Context, in *GetTickersInput) (*GetTickersOutput, error)
	GetTickerHistory(ctx context.Context, in *GetTickerHistoryInput) (*GetTickerHistoryOutput, error)
}

type Config struct {
	Storage storage.Storage
}

func New(cfg *Config) (Service, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	return &service{
		storage: cfg.Storage,
	}, nil
}

type service struct {
	storage storage.Storage
}
