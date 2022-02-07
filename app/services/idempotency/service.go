package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis"
)

// service calls Repository to be used in all methods of this package.
type service struct {
	repository redis.Repository
}

type Service interface {
	GetKeyValue(ctx context.Context, key string) (responses.Response, error)
	SetKeyValue(key string, res responses.Response) error
}

func NewService(r redis.Repository) Service {
	return &service{
		repository: r,
	}
}
