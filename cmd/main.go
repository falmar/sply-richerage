package main

import (
	"context"
	"github.com/falmar/richerage-api/cmd/http"
	"github.com/falmar/richerage-api/internal/bootstrap"
	"github.com/falmar/richerage-api/internal/pkg/zaplogger"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var rootCmd = &cobra.Command{}

func main() {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// set cobra flags
	v := viper.New()
	v.SetDefault("token.expired", true)
	bindFlags(v)

	err := rootCmd.ParseFlags(os.Args[1:])
	if err != nil {
		log.Fatal("parse flags failed:", err)
	}

	// initialize zap logger
	logger := zaplogger.New(v.GetBool("debug"))

	// handle stop signals
	go func() {
		sigChan := make(chan os.Signal)

		signal.Notify(sigChan, syscall.SIGINT)
		signal.Notify(sigChan, syscall.SIGTERM)

		<-sigChan
		cancel()
	}()

	// bootstrap configuration
	cfg, err := bootstrap.New(ctx, v, logger)
	if err != nil {
		logger.Fatal("bootstrap failed", zap.Error(err))
	}

	// add http server
	rootCmd.AddCommand(http.Cmd(ctx, cfg))

	if err := rootCmd.ExecuteContext(ctx); err != nil {
		logger.Error("main: error", zap.Error(err))
		os.Exit(1)
	}
}

func bindFlags(v *viper.Viper) {
	rootCmd.PersistentFlags().BoolP("debug", "d", false, "debug mode")
	v.BindPFlag("debug", rootCmd.PersistentFlags().Lookup("debug"))

	rootCmd.PersistentFlags().StringP("port", "p", "", "http port")
	v.BindPFlag("port", rootCmd.PersistentFlags().Lookup("port"))
}
