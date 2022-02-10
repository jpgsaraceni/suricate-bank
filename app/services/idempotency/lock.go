package idempotency

import (
	"context"
	"fmt"
)

// Lock creates the key in storage so if another request arrives before
// the server completes the operation, it will know that the request is already
// being processed.
func (s service) Lock(ctx context.Context, key string) error {

	err := s.repository.Lock(ctx, key)

	if err != nil {

		return fmt.Errorf("%w:%s", ErrRepository, err)
	}

	return nil
}
