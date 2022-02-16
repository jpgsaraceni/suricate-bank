package auth

import (
	"context"
)

type MockService struct {
	OnAuthenticate func(ctx context.Context, cpfInput, secret string) (string, error)
}

var _ Service = (*MockService)(nil)

func (mock MockService) Authenticate(ctx context.Context, cpfInput, secret string) (string, error) {
	return mock.OnAuthenticate(ctx, cpfInput, secret)
}
