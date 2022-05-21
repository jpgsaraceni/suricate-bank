package idempotency

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/config"
)

// Lock creates the key in storage so if another request arrives before
// the server completes the operation, it will know that the request is already
// being processed.
func (s service) Lock(ctx context.Context, cfg config.Config, key string) error {
	err := s.repository.Lock(ctx, cfg, key)
	if err != nil {
		return fmt.Errorf("%w:%s", ErrRepository, err)
	}

	return nil
}
