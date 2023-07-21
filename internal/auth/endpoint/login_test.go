//go:build test

package endpoint

import (
	"context"
	"errors"
	"github.com/falmar/richerage-api/internal/auth"
	"github.com/falmar/richerage-api/internal/auth/types"
	"github.com/falmar/richerage-api/internal/pkg/kit"
	"testing"
)

func getDefaultLoginRequest() *LoginRequest {
	return &LoginRequest{
		Username: "test",
		Password: "secret",
	}
}

func TestEndpointLogin(t *testing.T) {
	ctx := context.Background()

	svc := auth.NewMockService()
	svc.(*auth.MockService).LoginFunc = func(ctx context.Context, in *auth.LoginInput) (*auth.LoginOutput, error) {
		return &auth.LoginOutput{
			Token: "test",
		}, nil
	}

	req := getDefaultLoginRequest()

	// LoginEndpoint should return the token returned by the service
	resp, err := MakeLoginEndpoint(svc)(ctx, req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	}

	if r, ok := resp.(*LoginResponse); ok && r != nil {
		if r.Token != "test" {
			t.Errorf("expected token to be test, got %s", r.Token)
		}
	} else {
		t.Errorf("expected response to be of type LoginResponse, got %T", resp)
	}

	req = getDefaultLoginRequest()
	req.Username = ""

	_, err = MakeLoginEndpoint(svc)(ctx, req)
	var eBadRequest *kit.BadRequestError

	if !errors.As(err, &eBadRequest) {
		t.Errorf("expected error to be of type BadRequestError, got nil")
		return
	}
}

func TestEndpointLogin_Error(t *testing.T) {
	ctx := context.Background()
	svc := auth.NewMockService()
	svc.(*auth.MockService).LoginFunc = func(ctx context.Context, in *auth.LoginInput) (*auth.LoginOutput, error) {
		return nil, &types.ErrUnauthorized{}
	}

	req := getDefaultLoginRequest()

	// LoginEndpoint should return the error returned by the service
	_, err := MakeLoginEndpoint(svc)(ctx, req)

	var errUnauthorized *types.ErrUnauthorized
	if !errors.As(err, &errUnauthorized) {
		t.Errorf("expected error to be of type ErrUnauthorized, got %T", err)
	}
}

func TestEndpointLogin_VerifyRequest(t *testing.T) {
	var eBadRequest *kit.BadRequestError

	_, err := verifyLoginRequest(nil)

	if !errors.As(err, &eBadRequest) {
		t.Errorf("expected error to be of type BadRequestError, got nil")
	}

	req := getDefaultLoginRequest()

	req, err = verifyLoginRequest(req)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
	} else if req.Username != "test" {
		t.Errorf("expected username to be test, got %s", req.Username)
	} else if req.Password != "secret" {
		t.Errorf("expected password to be secret, got %s", req.Username)
	}

	req = getDefaultLoginRequest()
	req.Username = ""

	_, err = verifyLoginRequest(req)
	eBadRequest = nil
	if !errors.As(err, &eBadRequest) {
		t.Errorf("expected error to be of type BadRequestError, got nil")
		return
	}

	if len(eBadRequest.Params) != 1 {
		t.Errorf("expected error params to have length 1, got %d", len(eBadRequest.Params))
	}
	if v, ok := eBadRequest.Params["username"]; !ok || v == "" {
		t.Errorf("expected error param username to be required, got %s", eBadRequest.Params["username"])
	}

	req = getDefaultLoginRequest()
	req.Password = ""

	_, err = verifyLoginRequest(req)
	eBadRequest = nil
	if !errors.As(err, &eBadRequest) {
		t.Errorf("expected error to be of type BadRequestError, got nil")
		return
	}

	if len(eBadRequest.Params) != 1 {
		t.Errorf("expected error params to have length 1, got %d", len(eBadRequest.Params))
	}
	if v, ok := eBadRequest.Params["password"]; !ok || v == "" {
		t.Errorf("expected error param password to be required, got %s", eBadRequest.Params["password"])
	}
}
