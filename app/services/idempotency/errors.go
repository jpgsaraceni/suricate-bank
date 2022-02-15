package idempotency

import "errors"

var (
	ErrRepository     = errors.New("idempotency repository error")
	ErrResponseExists = errors.New("response already cached")
)
