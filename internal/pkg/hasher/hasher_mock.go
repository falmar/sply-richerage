//go:build test

package hasher

import (
	"context"
	"errors"
)

var _ Hasher = (*MockHasher)(nil)
var ErrMockUncalledFor = errors.New("uncalled for")

func NewMock() Hasher {
	return &MockHasher{
		GenerateTokenFunc: func(ctx context.Context, data []byte) ([]byte, error) {
			return nil, ErrMockUncalledFor
		},
		ValidateTokenFunc: func(ctx context.Context, token []byte) ([]byte, error) {
			return nil, ErrMockUncalledFor
		},
	}
}

type MockHasher struct {
	GenerateTokenFunc func(ctx context.Context, data []byte) ([]byte, error)
	ValidateTokenFunc func(ctx context.Context, token []byte) ([]byte, error)
}

func (m *MockHasher) GenerateToken(ctx context.Context, data []byte) ([]byte, error) {
	return m.GenerateTokenFunc(ctx, data)
}

func (m *MockHasher) ValidateToken(ctx context.Context, token []byte) ([]byte, error) {
	return m.ValidateTokenFunc(ctx, token)
}
