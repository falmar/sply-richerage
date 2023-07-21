package auth

import (
	"context"
)

type VerifyTokenInput struct {
	Token string
}

type VerifyTokenOutput struct {
	Username string
}

func (s *service) VerifyToken(ctx context.Context, in *VerifyTokenInput) (*VerifyTokenOutput, error) {
	c, err := s.hasher.ValidateToken(ctx, []byte(in.Token))
	if err != nil {
		return nil, err
	}

	return &VerifyTokenOutput{
		Username: string(c),
	}, nil
}
