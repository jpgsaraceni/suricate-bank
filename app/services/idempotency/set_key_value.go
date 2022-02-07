package idempotency

import "github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"

func (s service) SetKeyValue(key string, res responses.Response) error {
	return s.repository.SetKeyValue(key, res)
}
