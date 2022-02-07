package idempotency

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func (s service) GetKeyValue(ctx context.Context, key string) (responses.Response, error) {
	// TODO: s.repository.Get(key)

	return responses.Response{}, nil
}
