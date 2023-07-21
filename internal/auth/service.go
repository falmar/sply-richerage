package auth

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/pkg/hasher"
)

var ErrInvalidConfig = errors.New("invalid auth service config")

var _ Service = (*service)(nil)

type Service interface {
	Login(ctx context.Context, in *LoginInput) (*LoginOutput, error)
	VerifyToken(ctx context.Context, in *VerifyTokenInput) (*VerifyTokenOutput, error)
}

type Config struct {
	Hasher hasher.Hasher
}

func New(cfg *Config) (Service, error) {
	if cfg == nil {
		return nil, ErrInvalidConfig
	}

	return &service{
		hasher: cfg.Hasher,
	}, nil
}

type service struct {
	hasher hasher.Hasher
}
