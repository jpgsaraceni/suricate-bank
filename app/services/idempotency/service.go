package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

// service calls Repository to be used in all methods of this package.
type service struct {
	repository Repository
}

type Service interface {
	GetKeyValue(ctx context.Context, key string) (schema.CachedResponse, error)
	SetKeyValue(key string, request schema.CachedResponse) error
}

func NewService(r Repository) Service {
	return &service{
		repository: r,
	}
}
