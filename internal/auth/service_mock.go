//go:build test

package auth

import (
	"context"
	"errors"
)

var _ Service = (*MockService)(nil)
var ErrMockUncalledFor = errors.New("uncalled for")

func NewMockService() Service {
	return &MockService{
		LoginFunc: func(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
			return nil, ErrMockUncalledFor
		},
		VerifyTokenFunc: func(ctx context.Context, in *VerifyTokenInput) (*VerifyTokenOutput, error) {
			return nil, ErrMockUncalledFor
		},
	}
}

type MockService struct {
	LoginFunc       func(ctx context.Context, in *LoginInput) (*LoginOutput, error)
	VerifyTokenFunc func(ctx context.Context, in *VerifyTokenInput) (*VerifyTokenOutput, error)
}

func (m MockService) Login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	return m.LoginFunc(ctx, in)
}

func (m MockService) VerifyToken(ctx context.Context, in *VerifyTokenInput) (*VerifyTokenOutput, error) {
	return m.VerifyTokenFunc(ctx, in)
}
