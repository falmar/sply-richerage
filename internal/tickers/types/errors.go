package types

import "fmt"

type ErrTickerNotFound struct {
	Symbol string
}

func (e *ErrTickerNotFound) HttpCode() int {
	return 404
}

func (e *ErrTickerNotFound) Code() string {
	return "ticker_not_found"
}

func (e *ErrTickerNotFound) Error() string {
	return fmt.Sprintf("ticker %s not found", e.Symbol)
}
