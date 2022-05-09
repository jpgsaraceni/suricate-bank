package auth

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/config"
)

// service calls Repository to be used in all methods of this package.
type service struct {
	repository account.Repository
}

type Service interface {
	Authenticate(ctx context.Context, cfg config.Config, cpfInput, secret string) (string, error)
}

func NewService(r account.Repository) Service {
	return &service{
		repository: r,
	}
}
