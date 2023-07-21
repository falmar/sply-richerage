//go:build test

package endpoint

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/auth"
	types2 "github.com/falmar/richerage-api/internal/auth/types"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	"github.com/falmar/richerage-api/internal/tickers"
	"github.com/falmar/richerage-api/internal/tickers/types"
	"testing"
)

func getDefaultTickersRequest() *TickersRequest {
	return &TickersRequest{
		Username: "test",
	}
}

func TestEndpointTickers(t *testing.T) {
	ctx := context.Background()

	req := getDefaultTickersRequest()
	svc := tickers.NewMockService()
	svc.(*tickers.MockService).GetTickersFunc = func(ctx context.Context, in *tickers.GetTickersInput) (*tickers.GetTickersOutput, error) {
		return &tickers.GetTickersOutput{
			Tickers: []types.Ticker{
				{
					Symbol: "BTC",
					Price:  100,
				},
			},
		}, nil
	}

	// TickersEndpoint should call GetAll and return the response
	resp, err := MakeTickersEndpoint(svc)(ctx, req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	}
	if resp == nil {
		t.Errorf("expected response to be set, got nil")
		return
	} else if resp.(*TickersResponse).Tickers == nil {
		t.Errorf("expected response to have tickers set, got nil")
		return
	} else if len(resp.(*TickersResponse).Tickers) != 1 {
		t.Errorf("expected response to have 1 ticker, got %d", len(resp.(*TickersResponse).Tickers))
		return
	}

	if resp.(*TickersResponse).Tickers[0].Symbol != "BTC" {
		t.Errorf("expected ticker symbol to be BTC, got %s", resp.(*TickersResponse).Tickers[0].Symbol)
	}
	if resp.(*TickersResponse).Tickers[0].Price != 100 {
		t.Errorf("expected ticker price to be 100, got %f", resp.(*TickersResponse).Tickers[0].Price)
	}
}

func TestEndpointTickers_Error(t *testing.T) {
	ctx := context.Background()

	svcErr := errors.New("test error")

	req := getDefaultTickersRequest()
	svc := tickers.NewMockService()
	svc.(*tickers.MockService).GetTickersFunc = func(ctx context.Context, in *tickers.GetTickersInput) (*tickers.GetTickersOutput, error) {
		return nil, svcErr
	}

	// TickersEndpoint should return the error raised from the service
	_, err := MakeTickersEndpoint(svc)(ctx, req)
	if err != svcErr {
		t.Errorf("expected error got %T", err)
	}

	req = getDefaultTickersRequest()
	req.Username = ""

	// TickersEndpoint should return an error if the request is invalid
	_, err = MakeTickersEndpoint(svc)(ctx, req)
	if err == nil {
		t.Errorf("expected error to be set, got nil")
	}
}

func TestEndpointTickers_VerifyRequest(t *testing.T) {
	_, err := verifyTickersRequest(nil)
	if err == nil {
		t.Errorf("expected error to be set, got nil")
	}

	req := getDefaultTickersRequest()

	req, err = verifyTickersRequest(req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	} else if req.Username != "test" {
		t.Errorf("expected username to be test, got %s", req.Username)
	}

	req = getDefaultTickersRequest()
	req.Username = ""

	// TickersEndpoint should return an error if the request is invalid
	_, err = verifyTickersRequest(req)
	if err == nil {
		t.Errorf("expected error to be set, got nil")
		return
	}

	var errBadRequest *kit.BadRequestError
	if !errors.As(err, &errBadRequest) {
		t.Errorf("expected error to be of type BadRequestError, got %T", err)
	}
}

func TestEndpointTickers_Auth(t *testing.T) {
	ctx := context.Background()

	req := getDefaultTickersRequest()
	svc := auth.NewMockService()
	svc.(*auth.MockService).VerifyTokenFunc = func(ctx context.Context, in *auth.VerifyTokenInput) (*auth.VerifyTokenOutput, error) {
		if in.Token != "test" {
			return nil, &types2.ErrUnauthorized{}
		}

		return &auth.VerifyTokenOutput{
			Username: "john.doe",
		}, nil
	}

	// AuthEndpoint should call VerifyToken and set the username
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		if req, ok := request.(*TickersRequest); ok && req != nil {
			if req.Username != "john.doe" {
				t.Errorf("expected username to be john.doe, got %s", req.Username)
			}
		}

		return nil, nil
	}
	ctx = context.WithValue(ctx, "auth_token", "test")

	_, err := MakeTickersAuthEndpoint(svc, endpoint)(ctx, req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	}
}

func TestEndpointTickers_Auth_Error(t *testing.T) {
	ctx := context.Background()

	req := getDefaultTickersRequest()
	svc := auth.NewMockService()
	svc.(*auth.MockService).VerifyTokenFunc = func(ctx context.Context, in *auth.VerifyTokenInput) (*auth.VerifyTokenOutput, error) {
		if in.Token != "test" {
			return nil, &types2.ErrUnauthorized{}
		}

		return nil, nil
	}
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}

	// AuthEndpoint should return the error raised from auth service
	_, err := MakeTickersAuthEndpoint(svc, endpoint)(ctx, req)

	var errUnauthorized *types2.ErrUnauthorized
	if !errors.As(err, &errUnauthorized) {
		t.Errorf("expected error to be of type ErrUnauthorized, got %T", err)
	}
}
