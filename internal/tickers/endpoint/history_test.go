//go:build test

package endpoint

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/auth"
	authtypes "github.com/falmar/richerage-api/internal/auth/types"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	"github.com/falmar/richerage-api/internal/tickers"
	richeragetypes "github.com/falmar/richerage-api/internal/tickers/types"
	"testing"
	"time"
)

func getDefaultTickerHistoryRequest() *TickerHistoryRequest {
	return &TickerHistoryRequest{
		Username: "test",
		Symbol:   "AAPL",
		Before:   "2023-07-21T00:00:00Z",
	}
}

func TestEndpointHistory(t *testing.T) {
	ctx := context.Background()

	date, _ := time.Parse(time.RFC3339, "2023-07-21T00:00:00Z")

	req := getDefaultTickerHistoryRequest()
	svc := tickers.NewMockService()
	svc.(*tickers.MockService).GetTickerHistoryFunc = func(ctx context.Context, in *tickers.GetTickerHistoryInput) (*tickers.GetTickerHistoryOutput, error) {
		return &tickers.GetTickerHistoryOutput{
			History: []richeragetypes.TickerHistory{
				{
					Price: 100,
					Date:  date,
				},
			},
		}, nil
	}

	// TickerHistoryEndpoint should call svc.GetTickerHistory and return the response
	resp, err := MakeTickerHistoryEndpoint(svc)(ctx, req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	}

	endRes, ok := resp.(*TickerHistoryResponse)

	if !ok || endRes == nil {
		t.Errorf("expected response to be set, got nil")
		return
	} else if endRes.Tickers == nil {
		t.Errorf("expected response to have tickers set, got nil")
		return
	} else if len(endRes.Tickers) != 1 {
		t.Errorf("expected response to have 1 ticker, got %d", len(endRes.Tickers))
		return
	}

	if endRes.Tickers[0].Price != 100 {
		t.Errorf("expected ticker price to be 100, got %f", endRes.Tickers[0].Price)
	}
	if endRes.Tickers[0].Date != date {
		t.Errorf("expected ticker date to be %s, got %s", date, endRes.Tickers[0].Date)
	}
}

func TestEndpoint_History_Error(t *testing.T) {
	ctx := context.Background()

	svcError := errors.New("svc error")
	symbolError := errors.New("symbol error")
	beforeError := errors.New("before error")

	date, _ := time.Parse(time.RFC3339, "2023-07-21T00:00:00Z")

	req := getDefaultTickerHistoryRequest()
	svc := tickers.NewMockService()
	svc.(*tickers.MockService).GetTickerHistoryFunc = func(ctx context.Context, in *tickers.GetTickerHistoryInput) (*tickers.GetTickerHistoryOutput, error) {
		if in.Symbol != "AAPL" {
			return nil, symbolError
		}
		if in.Before != date {
			return nil, beforeError
		}

		return nil, svcError
	}

	// TickerHistoryEndpoint should return the error raised from the svc.GetTickerHistory
	_, err := MakeTickerHistoryEndpoint(svc)(ctx, req)
	if err != svcError {
		t.Errorf("expected error to be nil, got %T", err)
	}

	req = getDefaultTickerHistoryRequest()
	req.Username = ""
	req.Symbol = ""

	// TickerHistoryEndpoint should return the error raised from verify
	_, err = MakeTickerHistoryEndpoint(svc)(ctx, req)
	if err == nil {
		t.Errorf("expected error to be set, got nil")
	}
}

func TestEndpointHistory_VerifyRequest_Nil(t *testing.T) {
	_, err := verifyTickerHistoryRequest(nil)

	var badRequest *kit.BadRequestError
	if !errors.As(err, &badRequest) {
		t.Errorf("expected error to be of type BadRequestError, got %T", err)
	}
}

func TestEndpointHistory_VerifyRequest_Symbol(t *testing.T) {
	req := getDefaultTickerHistoryRequest()
	req.Symbol = ""

	req, err := verifyTickerHistoryRequest(req)

	var badRequest *kit.BadRequestError
	if !errors.As(err, &badRequest) {
		t.Errorf("expected error to be of type BadRequestError, got %T", err)
	}
	if badRequest.Params == nil {
		t.Errorf("expected bad request parameters to be set, got nil")
	} else if v, ok := badRequest.Params["symbol"]; !ok || v == "" {
		t.Errorf("expected bad request parameter symbol error message, got %s", v)
	}

	req = getDefaultTickerHistoryRequest()
	req.Symbol = "AAPL"

	req, err = verifyTickerHistoryRequest(req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	}
	if req.Symbol != "AAPL" {
		t.Errorf("expected symbol to be unchanged, got %s", req.Symbol)
	}
}

func TestEndpointHistory_VerifyRequest_Before(t *testing.T) {
	// Test valid format
	req := getDefaultTickerHistoryRequest()
	req.Before = "2023-07-21T00:00:00Z"

	_, err := verifyTickerHistoryRequest(req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	} else if req.Before != "2023-07-21T00:00:00Z" {
		t.Errorf("expected before to be unchanged, got %s", req.Before)
	}

	// Test invalid format
	req = getDefaultTickerHistoryRequest()
	req.Before = "invalid"

	req, err = verifyTickerHistoryRequest(req)

	var badRequest *kit.BadRequestError
	if !errors.As(err, &badRequest) {
		t.Errorf("expected error to be of type BadRequestError, got %T", err)
	}
	if badRequest.Params == nil {
		t.Errorf("expected bad request parameters to be set, got nil")
	} else if v, ok := badRequest.Params["before"]; !ok || v == "" {
		t.Errorf("expected bad request parameter before error message, got %s", v)
	}

	// Test empty value
	req = getDefaultTickerHistoryRequest()
	req.Before = ""

	req, err = verifyTickerHistoryRequest(req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	} else if req.Before != "" {
		t.Errorf("expected before to be unchanged, got %s", req.Before)
	}
}

func TestEndpointHistory_VerifyRequest_Username(t *testing.T) {
	req := getDefaultTickerHistoryRequest()
	req.Username = ""

	req, err := verifyTickerHistoryRequest(req)

	var badRequest *kit.BadRequestError
	if !errors.As(err, &badRequest) {
		t.Errorf("expected error to be of type BadRequestError, got %T", err)
	}
	if badRequest.Params == nil {
		t.Errorf("expected bad request parameters to be set, got nil")
	} else if v, ok := badRequest.Params["username"]; !ok || v == "" {
		t.Errorf("expected bad request parameter username error message, got %s", v)
	}

	req = getDefaultTickerHistoryRequest()
	req.Username = "test"

	req, err = verifyTickerHistoryRequest(req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	}
	if req.Username != "test" {
		t.Errorf("expected username to be unchanged, got %s", req.Username)
	}
}

func TestEndpointHistory_Auth(t *testing.T) {
	ctx := context.Background()

	// use empty request
	req := getDefaultTickerHistoryRequest()
	req.Username = ""

	svc := auth.NewMockService()
	svc.(*auth.MockService).VerifyTokenFunc = func(ctx context.Context, in *auth.VerifyTokenInput) (*auth.VerifyTokenOutput, error) {
		if in.Token != "test" {
			return nil, &authtypes.ErrUnauthorized{}
		}

		return &auth.VerifyTokenOutput{
			Username: "john.doe",
		}, nil
	}

	// AuthEndpoint should call VerifyToken and set the username
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		if req, ok := request.(*TickerHistoryRequest); ok && req != nil {
			if req.Username != "john.doe" {
				t.Errorf("expected username to be john.doe, got %s", req.Username)
			}
		}

		return nil, nil
	}
	ctx = context.WithValue(ctx, "auth_token", "test")

	_, err := MakeTickerHistoryAuthEndpoint(svc, endpoint)(ctx, req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	}
}

func TestEndpointHistory_Auth_Error(t *testing.T) {
	ctx := context.Background()

	req := getDefaultTickerHistoryRequest()
	req.Username = ""

	svc := auth.NewMockService()
	svc.(*auth.MockService).VerifyTokenFunc = func(ctx context.Context, in *auth.VerifyTokenInput) (*auth.VerifyTokenOutput, error) {
		if in.Token != "test" {
			return nil, &authtypes.ErrUnauthorized{}
		}

		return nil, nil
	}
	endpoint := func(ctx context.Context, request interface{}) (interface{}, error) {
		return nil, nil
	}

	// AuthEndpoint should return the error raised from auth service
	_, err := MakeTickerHistoryAuthEndpoint(svc, endpoint)(ctx, req)

	var errUnauthorized *authtypes.ErrUnauthorized
	if !errors.As(err, &errUnauthorized) {
		t.Errorf("expected error to be of type ErrUnauthorized, got %T", err)
	}
}
