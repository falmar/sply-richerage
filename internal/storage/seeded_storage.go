package storage

import (
	"context"
	"crypto/sha1"
	"encoding/binary"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"math/rand"
	"time"
)

var _ Storage = (*seededStorage)(nil)

func NewSeeded() Storage {
	return &seededStorage{}
}

type seededStorage struct{}

func (s *seededStorage) GetByUser(_ context.Context, username string) ([]types.Ticker, error) {
	genTickers := generatedTickers(getRandForString("tickers"))
	rnd := getRandForString(username)
	maxTickers := rnd.Intn(len(genTickers))

	if maxTickers == 0 {
		// at least one ticker
		maxTickers = 1
	} else if maxTickers > 10 {
		// at most 10 tickers
		maxTickers = 10
	}

	// allocate
	tickers := make([]types.Ticker, 0, maxTickers)

	for i := 0; i < maxTickers; i++ {
		index := rnd.Intn(len(genTickers))
		var repeat bool

		for _, v := range tickers {
			if v.Symbol == genTickers[index].Symbol {
				repeat = true
				break
			}
		}

		if repeat {
			continue
		}

		tickers = append(tickers, genTickers[index])
	}

	return tickers, nil
}

func (s *seededStorage) GetHistory(_ context.Context, symbol string, before time.Time) ([]types.TickerHistory, error) {
	if !isValidSymbol(types.ValidTickers(), symbol) {
		return nil, &types.ErrTickerNotFound{
			Symbol: symbol,
		}
	}

	genTickers := generatedTickers(getRandForString("tickers"))

	// obtain a deterministic random number for price and date given the symbol
	rnd := getRandForString(symbol)

	// start from today at 00:00:00
	date := time.Now().Truncate(time.Hour * 24)

	// up to 100 records
	records := rnd.Intn(100)
	if records == 0 {
		// at least one record
		records = 1
	}

	// allocate records of history
	history := make([]types.TickerHistory, 0, records)

	// add the current price for today
	for _, v := range genTickers {
		if v.Symbol == symbol {
			history = append(history, types.TickerHistory{
				Date:  date,
				Price: v.Price,
			})
			break
		}
	}

	for i := 0; i < records; i++ {
		base := rnd.Intn(1000)
		decimals := rnd.Intn(100)

		// from today, down to 90 days ago

		var historyDate time.Time

	dateLoop:
		for {
			historyDate = date.Add(-time.Hour * 24 * time.Duration(rnd.Intn(90)+1))

			// check if the date is already in the history
			for _, v := range history {
				if v.Date == historyDate {
					continue dateLoop
				}
			}

			break dateLoop
		}

		history = append(history, types.TickerHistory{
			Price: float64(base) + (float64(decimals) / 100),
			Date:  historyDate,
		})
	}

	if !before.IsZero() {
		// remove all dates after before
		for i, v := range history {
			if v.Date.After(before) {
				history = history[:i]
				break
			}
		}
	}

	return history, nil
}

func generatedTickers(rnd *rand.Rand) []types.Ticker {
	validTickers := types.ValidTickers()

	// allocate
	tickers := make([]types.Ticker, 0, len(validTickers))

	for i := 0; i < len(validTickers); i++ {
		base := rnd.Intn(1000)
		decimals := rnd.Intn(100)

		tickers = append(tickers, types.Ticker{
			Symbol: validTickers[i],
			Price:  float64(base) + (float64(decimals) / 100),
		})
	}

	return tickers
}

func getRandForString(symbol string) *rand.Rand {
	hasher := sha1.New()
	hasher.Write([]byte(symbol))
	hash := hasher.Sum(nil)
	seed := binary.BigEndian.Uint64(hash)

	return rand.New(rand.NewSource(int64(seed)))
}
