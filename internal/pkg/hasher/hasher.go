package hasher

import "context"

type Hasher interface {
	GenerateToken(ctx context.Context, data []byte) ([]byte, error)

	ValidateToken(ctx context.Context, token []byte) ([]byte, error)
}
