package auth

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

// usecase calls Repository to be used in all methods of this package.
type service struct {
	repository account.Repository
}

type Service interface {
	Authenticate(ctx context.Context, cpf cpf.Cpf, secret string) (string, error) // TODO: check cpf here?
}

func NewService(r account.Repository) Service {
	return &service{
		repository: r,
	}
}
