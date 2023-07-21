package endpoint

import (
	"context"
	"fmt"
	"github.com/falmar/richerage-api/internal/auth"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	"github.com/falmar/richerage-api/internal/tickers"
	"github.com/falmar/richerage-api/internal/tickers/types"
	kitendpoint "github.com/go-kit/kit/endpoint"
	"time"
)

type TickerHistoryRequest struct {
	Username string

	Symbol string
	Before string
}

type TickerHistoryResponse struct {
	Tickers []types.TickerHistory
}

func MakeTickerHistoryEndpoint(svc tickers.Service) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req, err := verifyTickerHistoryRequest(request)
		if err != nil {
			return nil, err
		}

		var before = time.Time{}
		if req.Before != "" {
			before, _ = time.Parse(time.RFC3339, req.Before)
		}

		out, err := svc.GetTickerHistory(ctx, &tickers.GetTickerHistoryInput{
			Symbol: req.Symbol,
			Before: before,
		})

		if err != nil {
			return nil, err
		}

		return &TickerHistoryResponse{
			Tickers: out.History,
		}, nil
	}
}

func MakeTickerHistoryAuthEndpoint(svc auth.Service, e kitendpoint.Endpoint) kitendpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		token, _ := ctx.Value("auth_token").(string)

		// assume that the username is valid and exists in the database
		out, err := svc.VerifyToken(ctx, &auth.VerifyTokenInput{
			Token: token,
		})
		if err != nil {
			return nil, err
		}

		if req, ok := request.(*TickerHistoryRequest); ok && req != nil {
			req.Username = out.Username
		}

		return e(ctx, request)
	}
}

func verifyTickerHistoryRequest(request interface{}) (*TickerHistoryRequest, error) {
	req, ok := request.(*TickerHistoryRequest)
	if !ok || req == nil {
		return nil, &kit.BadRequestError{
			Message: "invalid request",
		}
	}

	badParams := map[string]string{}

	if req.Username == "" {
		badParams["username"] = "required"
	}
	if req.Symbol == "" {
		badParams["symbol"] = "required"
	}
	if req.Before != "" {
		if _, err := time.Parse(time.RFC3339, req.Before); err != nil {
			badParams["before"] = fmt.Sprintf("invalid format %s: %s", req.Before, err.Error())
		}
	}

	if len(badParams) > 0 {
		return nil, &kit.BadRequestError{
			Message: "one or more parameters are invalid or missing",
			Params:  badParams,
		}
	}

	return req, nil
}
