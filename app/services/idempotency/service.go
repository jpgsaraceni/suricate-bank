package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis"
)

// service calls Repository to be used in all methods of this package.
type service struct {
	repository redis.Repository
}

type Service interface {
	Idempotency(ctx context.Context, key string) error
}

func NewService(r redis.Repository) Service {
	return &service{
		repository: r,
	}
}
