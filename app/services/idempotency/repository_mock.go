package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"

type MockRepository struct {
	OnGetKeyValue func(key string) (schema.CachedResponse, error)
	OnSetKeyValue func(key string, request schema.CachedResponse) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) GetKeyValue(key string) (schema.CachedResponse, error) {
	return mock.OnGetKeyValue(key)
}

func (mock MockRepository) SetKeyValue(key string, res schema.CachedResponse) error {
	return mock.OnSetKeyValue(key, res)
}
