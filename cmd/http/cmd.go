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
	"github.com/spf13/cobra"
	"go.uber.org/zap"
	"net/http"
)

func Cmd(_ context.Context, config *bootstrap.Config) *cobra.Command {
	return &cobra.Command{
		Use:   "http",
		Short: "Start HTTP server",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			router := chi.NewRouter()

			loggerHandler := &kit.LoggerHandler{
				Logger: config.Logger,
			}
			transportHandler := &kit.ErrorHandler{
				Logger: config.Logger,
			}

			loginEndpoint := authendpoints.MakeLoginEndpoint(config.AuthService)
			router.Method("POST", "/login", kithttp.NewServer(
				loginEndpoint,
				authtransport.LoginRequestDecoder,
				authtransport.LoginResponseEncoder,
				kithttp.ServerErrorEncoder(transportHandler.ErrorEncoder),
				kithttp.ServerErrorHandler(transportHandler),
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
				kithttp.ServerErrorEncoder(transportHandler.ErrorEncoder),
				kithttp.ServerErrorHandler(transportHandler),
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
				kithttp.ServerErrorEncoder(transportHandler.ErrorEncoder),
				kithttp.ServerErrorHandler(transportHandler),
				kithttp.ServerBefore(loggerHandler.Before),
				kithttp.ServerAfter(loggerHandler.After),
			))

			port := config.Viper.GetString("port")
			if port == "" {
				port = "8080"
			}

			server := &http.Server{
				Addr: ":" + port,
			}
			server.Handler = router

			go func() {
				<-ctx.Done()
				config.Logger.Info("http: shutdown signal received")
				_ = server.Shutdown(ctx)
			}()

			config.Logger.Info("http: starting server", zap.String("port", port))

			err := server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				return err
			}

			return nil
		},
	}
}
