package auth

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/config"
)

type MockService struct {
	OnAuthenticate func(ctx context.Context, cfg config.Config, cpfInput, secret string) (string, error)
}

var _ Service = (*MockService)(nil)

func (mock MockService) Authenticate(ctx context.Context, cfg config.Config, cpfInput, secret string) (string, error) {
	return mock.OnAuthenticate(ctx, cfg, cpfInput, secret)
}
