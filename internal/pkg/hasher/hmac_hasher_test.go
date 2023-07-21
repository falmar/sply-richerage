package hasher

import (
	"bytes"
	"context"
	"errors"
	"testing"
	"time"
)

func TestHmacHasher(t *testing.T) {
	ctx := context.Background()

	hasher := NewHMAC(&ConfigHMAC{
		Secret:       []byte("secret"),
		TTL:          1 * time.Second,
		CheckExpired: true,
	})

	data := []byte("data")

	token, err := hasher.GenerateToken(ctx, data)
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	if token == nil {
		t.Errorf("expected token to be set, got nil")
	}

	payload, err := hasher.ValidateToken(ctx, token)
	if err != nil {
		t.Errorf("expected error to be nil, got %T %s", err, err.Error())
		return
	}

	if payload == nil {
		t.Errorf("expected payload to be set, got nil")
	} else if !bytes.Equal(payload, data) {
		t.Errorf("expected payload to be %s, got %s", string(data), string(payload))
	}
}

func TestHmacHasher_Invalid(t *testing.T) {
	ctx := context.Background()
	hasher := NewHMAC(&ConfigHMAC{
		Secret:       []byte("secret"),
		TTL:          1 * time.Second,
		CheckExpired: true,
	})

	var data []byte

	_, err := hasher.GenerateToken(ctx, data)

	var errEmptyData *ErrEmptyData

	if !errors.As(err, &errEmptyData) {
		t.Errorf("expected error to be %T, got %T", errEmptyData, err)
	}

	token := []byte("valid.token=")
	_, err = hasher.ValidateToken(ctx, token)

	var errInvalidToken *ErrInvalidToken
	if !errors.As(err, &errInvalidToken) {
		t.Errorf("expected error to be %T, got nil", err)
		return
	}

	errInvalidToken = nil
	token = []byte("ZCBzdHJpbmcgMzIgYnl0ZXM=.YW5vdGhlcg==.1690493268")

	_, err = hasher.ValidateToken(ctx, token)
	if !errors.As(err, &errInvalidToken) {
		t.Errorf("expected error to be %T, got %T", errInvalidToken, err)
		return
	}
}

func TestHmacHasher_Expired(t *testing.T) {
	ctx := context.Background()
	hasher := NewHMAC(&ConfigHMAC{
		Secret:       []byte("secret"),
		TTL:          1 * time.Nanosecond,
		CheckExpired: true,
	})

	token, err := hasher.GenerateToken(ctx, []byte("data"))
	if err != nil {
		t.Errorf("expected error to be nil, got %T", err)
		return
	}

	time.Sleep(1 * time.Second)
	_, err = hasher.ValidateToken(ctx, token)

	var errExpiredToken *ErrExpiredToken
	if !errors.As(err, &errExpiredToken) {
		t.Errorf("expected error to be %T, got %T", errExpiredToken, err)
		return
	}
}
