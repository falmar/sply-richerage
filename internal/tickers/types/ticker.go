package types

import "time"

type Ticker struct {
	Symbol string  `json:"symbol"`
	Price  float64 `json:"price"`
}

type TickerHistory struct {
	Date  time.Time `json:"date"`
	Price float64   `json:"price"`
}

func ValidTickers() []string {
	return []string{
		"AAPL",
		"MSFT",
		"GOOG",
		"AMZN",
		"FB",
		"TSLA",
		"NVDA",
		"JPM",
		"BABA",
		"JNJ",
		"WMT",
		"PG",
		"PYPL",
		"DIS",
		"ADBE",
		"PFE",
		"V",
		"MA",
		"CRM",
		"NFLX",
	}
}
