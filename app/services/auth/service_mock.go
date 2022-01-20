package auth

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

type MockService struct {
	OnAuthenticate func(ctx context.Context, cpf cpf.Cpf, secret string) (string, error)
}

var _ Service = (*MockService)(nil)

func (mock MockService) Authenticate(ctx context.Context, cpf cpf.Cpf, secret string) (string, error) {
	return mock.OnAuthenticate(ctx, cpf, secret)
}
