package bootstrap

import (
	"context"
	"encoding/base64"
	"github.com/falmar/richerage-api/internal/auth"
	"github.com/falmar/richerage-api/internal/pkg/hasher"
	"github.com/falmar/richerage-api/internal/storage"
	"github.com/falmar/richerage-api/internal/tickers"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"strings"
	"time"
)

type Config struct {
	Logger *zap.Logger
	Viper  *viper.Viper

	AuthService      auth.Service
	RicherageService tickers.Service
}

func New(_ context.Context, v *viper.Viper, logger *zap.Logger) (*Config, error) {
	var err error = nil
	cfg := &Config{}

	cfg.Logger = logger
	cfg.Viper = v

	cfg.Viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))
	cfg.Viper.AutomaticEnv()

	// DISCLAMER: I know that this is not a good practice; however, I'm doing this for the sake of simplicity
	// could regenerate with crypt/rand at server restart but would require re-logins when testing the api
	// in order to save your time, I'll just use a static key instead of using config/env files through viper
	sk, _ := base64.StdEncoding.DecodeString("ZCBzZWNyZXQga2V5IDMyIGJ5dGVz")

	cfg.AuthService, err = auth.New(&auth.Config{
		Hasher: hasher.NewHMAC(&hasher.ConfigHMAC{
			Secret: sk,

			// 1 week
			TTL:          time.Hour * 24 * 7 * 30,
			CheckExpired: v.GetBool("token.expired"),
		}),
	})

	// bootstrap dependencies
	cfg.RicherageService, err = tickers.New(&tickers.Config{
		Storage: storage.NewSeeded(),
	})
	if err != nil {
		return nil, err
	}

	return cfg, nil
}
