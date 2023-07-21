package http

import (
	"context"
	"github.com/falmar/richerage-api/internal/bootstrap"
	apphttp "github.com/falmar/richerage-api/internal/http"
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

			port := config.Viper.GetString("port")
			if port == "" {
				port = "8080"
			}

			handler, err := apphttp.Handler(ctx, config)
			if err != nil {
				return err
			}

			server := &http.Server{
				Addr: ":" + port,
			}
			server.Handler = handler

			go func() {
				<-ctx.Done()
				config.Logger.Info("http: shutdown signal received")
				_ = server.Shutdown(ctx)
			}()

			config.Logger.Info("http: starting server", zap.String("port", port))

			err = server.ListenAndServe()
			if err != nil && err != http.ErrServerClosed {
				return err
			}

			return nil
		},
	}
}
