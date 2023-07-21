package storage

import (
	"context"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"time"
)

// Why is storage not part of the "tickers service"?
// to allow plug and play of different storage implementations
// although, it becomes more boilerplate code, it allows to easily change the storage implementation

type Storage interface {
	GetByUser(ctx context.Context, username string) ([]types.Ticker, error)

	GetHistory(ctx context.Context, symbol string, before time.Time) ([]types.TickerHistory, error)
}

func isValidSymbol(tickers []string, symbol string) bool {
	for _, v := range tickers {
		if v == symbol {
			return true
		}
	}

	return false
}
