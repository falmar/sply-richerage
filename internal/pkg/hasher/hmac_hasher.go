package hasher

import (
	"bytes"
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"strconv"
	"time"
)

var _ Hasher = (*hmacHasher)(nil)

type ConfigHMAC struct {
	Secret []byte
	TTL    time.Duration

	CheckExpired bool
}

func NewHMAC(cfg *ConfigHMAC) Hasher {
	return &hmacHasher{
		secret:       cfg.Secret,
		ttl:          cfg.TTL,
		checkExpired: cfg.CheckExpired,
	}
}

type hmacHasher struct {
	secret       []byte
	ttl          time.Duration
	checkExpired bool
}

func (m *hmacHasher) GenerateToken(_ context.Context, data []byte) ([]byte, error) {
	if data == nil || len(data) == 0 {
		return nil, &ErrEmptyData{}
	}

	mac := hmac.New(sha256.New, m.secret)

	// make token expire in time.Duration from now
	ts := strconv.FormatInt(
		time.Now().Add(m.ttl).UTC().Unix(),
		10,
	)

	payload := base64.StdEncoding.EncodeToString(data) + "." + ts

	_, err := mac.Write([]byte(payload))
	if err != nil {
		return nil, err
	}

	sig := mac.Sum(nil)

	return []byte(base64.StdEncoding.EncodeToString(sig) + "." + payload), nil
}

func (m *hmacHasher) ValidateToken(_ context.Context, token []byte) ([]byte, error) {
	parts := bytes.Split(token, []byte("."))
	if len(parts) != 3 {
		return nil, &ErrInvalidToken{
			Message: "invalid format",
		}
	}

	mac := hmac.New(sha256.New, m.secret)

	// build payload
	var payload = make([]byte, 0, len(parts[1])+len(parts[2])+1)
	payload = append(payload, parts[1]...)
	payload = append(payload, []byte(".")...)
	payload = append(payload, parts[2]...)

	_, err := mac.Write(payload)
	if err != nil {
		return nil, err
	}

	// decode signature
	uSig, err := base64.StdEncoding.DecodeString(string(parts[0]))
	if err != nil {
		return nil, err
	}

	sig := mac.Sum(nil)

	if !hmac.Equal(sig, uSig) {
		return nil, &ErrInvalidToken{
			Message: "invalid signature",
		}
	}

	// check date
	date, err := strconv.ParseInt(string(parts[2]), 10, 64)

	if err != nil {
		return nil, &ErrInvalidToken{
			Message: "invalid date",
		}
	}

	if m.checkExpired && time.Now().UTC().Unix() > date {
		return nil, &ErrExpiredToken{}
	}

	decoded, err := base64.StdEncoding.DecodeString(string(parts[1]))
	if err != nil {
		return nil, err
	}

	return decoded, nil
}
