package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"

type Repository interface {
	GetKeyValue(key string) (responses.Response, error)
	SetKeyValue(key string, res responses.Response) error
}
