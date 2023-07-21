package auth

import (
	"context"
)

type LoginInput struct {
	Username string
	Password string
}

type LoginOutput struct {
	Token string
}

func (s *service) Login(ctx context.Context, in *LoginInput) (*LoginOutput, error) {
	// assumes user is found and password is correct
	// generate token
	token, err := s.hasher.GenerateToken(ctx, []byte(in.Username))
	if err != nil {
		return nil, err
	}

	return &LoginOutput{
		Token: string(token),
	}, nil
}
