package endpoint

import (
	"context"
	"github.com/falmar/richerage-api/internal/auth"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	kitendpoint "github.com/go-kit/kit/endpoint"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func MakeLoginEndpoint(svc auth.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (response interface{}, err error) {
		req, err := verifyLoginRequest(request)
		if err != nil {
			return nil, err
		}

		out, err := svc.Login(ctx, &auth.LoginInput{
			Username: req.Username,
			Password: req.Password,
		})
		if err != nil {
			return nil, err
		}

		return &LoginResponse{
			Token: out.Token,
		}, nil
	}
}

func verifyLoginRequest(request interface{}) (*LoginRequest, error) {
	req, ok := request.(*LoginRequest)
	if !ok || req == nil {
		return nil, &kit.BadRequestError{
			Message: "invalid request",
		}
	}

	badParams := map[string]string{}

	if req.Username == "" {
		badParams["username"] = "required"
	}
	if req.Password == "" {
		badParams["password"] = "required"
	}

	if len(badParams) > 0 {
		return nil, &kit.BadRequestError{
			Message: "one or more parameters are invalid or missing",
			Params:  badParams,
		}
	}

	return req, nil
}
