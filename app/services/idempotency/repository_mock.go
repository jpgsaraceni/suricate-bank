package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"

type MockRepository struct {
	OnGetKeyValue func(key string) (responses.Response, error)
	OnSetKeyValue func(key string, res responses.Response) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) GetKeyValue(key string) (responses.Response, error) {
	return mock.OnGetKeyValue(key)
}

func (mock MockRepository) SetKeyValue(key string, res responses.Response) error {
	return mock.OnSetKeyValue(key, res)
}
