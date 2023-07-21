package endpoint

import (
	"context"
	"github.com/falmar/richerage-api/internal/auth"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	"github.com/falmar/richerage-api/internal/tickers"
	"github.com/falmar/richerage-api/internal/tickers/types"
	kitendpoint "github.com/go-kit/kit/endpoint"
)

type TickersRequest struct {
	Username string
}

type TickersResponse struct {
	Tickers []types.Ticker
}

func MakeTickersEndpoint(svc tickers.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := verifyTickersRequest(request)
		if err != nil {
			return nil, err
		}

		out, err := svc.GetTickers(ctx, &tickers.GetTickersInput{
			Username: req.Username,
		})
		if err != nil {
			return nil, err
		}

		return &TickersResponse{
			Tickers: out.Tickers,
		}, nil
	}
}

func MakeTickersAuthEndpoint(svc auth.Service, e kitendpoint.Endpoint) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		token, _ := ctx.Value("auth_token").(string)

		// assume that the username is valid and exists in the database
		out, err := svc.VerifyToken(ctx, &auth.VerifyTokenInput{
			Token: token,
		})
		if err != nil {
			return nil, err
		}

		if req, ok := request.(*TickersRequest); ok && req != nil {
			req.Username = out.Username
		}

		return e(ctx, request)
	}
}

func verifyTickersRequest(request interface{}) (*TickersRequest, error) {
	req, ok := request.(*TickersRequest)
	if !ok || req == nil {
		return nil, &kit.BadRequestError{
			Message: "invalid request",
		}
	}

	badParams := map[string]string{}

	if req.Username == "" {
		badParams["username"] = "required"
	}

	if len(badParams) > 0 {
		return nil, &kit.BadRequestError{
			Message: "one or more parameters are invalid or missing",
			Params:  badParams,
		}
	}

	return req, nil
}
