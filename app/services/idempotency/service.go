package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
	"github.com/jpgsaraceni/suricate-bank/config"
)

// service calls Repository to be used in all methods of this package.
type service struct {
	repository Repository
}

type Service interface {
	GetCachedResponse(ctx context.Context, key string) (schema.CachedResponse, error)
	CacheResponse(ctx context.Context, request schema.CachedResponse) error
	Lock(ctx context.Context, cfg config.Config, key string) error
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
