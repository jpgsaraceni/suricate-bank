package idempotency

import "context"

func (s service) Idempotency(ctx context.Context, key string) error {
	return nil
}
