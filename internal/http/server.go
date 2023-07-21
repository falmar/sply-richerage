package http

import (
	"context"
	authendpoints "github.com/falmar/richerage-api/internal/auth/endpoint"
	authtransport "github.com/falmar/richerage-api/internal/auth/transport"
	"github.com/falmar/richerage-api/internal/bootstrap"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	tickersendpoints "github.com/falmar/richerage-api/internal/tickers/endpoint"
	tickerstransport "github.com/falmar/richerage-api/internal/tickers/transport"
	"github.com/go-chi/chi/v5"
	kithttp "github.com/go-kit/kit/transport/http"
	"net/http"
)

func Handler(_ context.Context, config *bootstrap.Config) (http.Handler, error) {
	router := chi.NewRouter()

	loggerHandler := &kit.LoggerHandler{
		Logger: config.Logger,
	}
	errorHandler := &kit.ErrorHandler{
		Logger: config.Logger,
	}

	loginEndpoint := authendpoints.MakeLoginEndpoint(config.AuthService)
	router.Method("POST", "/login", kithttp.NewServer(
		loginEndpoint,
		authtransport.LoginRequestDecoder,
		authtransport.LoginResponseEncoder,
		kithttp.ServerErrorEncoder(errorHandler.ErrorEncoder),
		kithttp.ServerErrorHandler(errorHandler),
		kithttp.ServerBefore(loggerHandler.Before),
		kithttp.ServerAfter(loggerHandler.After),
	))

	tickerEndpoint := tickersendpoints.MakeTickersEndpoint(config.RicherageService)
	tickerEndpoint = tickersendpoints.MakeTickersAuthEndpoint(config.AuthService, tickerEndpoint)
	router.Method("GET", "/tickers", kithttp.NewServer(
		tickerEndpoint,
		tickerstransport.TickersRequestDecoder,
		tickerstransport.TickersResponseEncoder,
		kithttp.ServerBefore(tickerstransport.TokenDecoder),
		kithttp.ServerErrorEncoder(errorHandler.ErrorEncoder),
		kithttp.ServerErrorHandler(errorHandler),
		kithttp.ServerBefore(loggerHandler.Before),
		kithttp.ServerAfter(loggerHandler.After),
	))

	historyEndpoint := tickersendpoints.MakeTickerHistoryEndpoint(config.RicherageService)
	historyEndpoint = tickersendpoints.MakeTickerHistoryAuthEndpoint(config.AuthService, historyEndpoint)
	router.Method("GET", "/tickers/{symbol}/history", kithttp.NewServer(
		historyEndpoint,
		tickerstransport.TickerHistoryRequestDecoder,
		tickerstransport.TickerHistoryResponseEncoder,
		kithttp.ServerBefore(tickerstransport.TokenDecoder),
		kithttp.ServerErrorEncoder(errorHandler.ErrorEncoder),
		kithttp.ServerErrorHandler(errorHandler),
		kithttp.ServerBefore(loggerHandler.Before),
		kithttp.ServerAfter(loggerHandler.After),
	))

	return router, nil
}
