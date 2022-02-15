package redis

import (
	"context"
	"fmt"
	"os"
	"strconv"
)

func (r Repository) Lock(ctx context.Context, key string) error {
	conn := r.pool.Get()
	defer conn.Close()

	ttl, err := strconv.Atoi(os.Getenv("IDEMPOTENCY_TTL"))

	if err != nil {

		return fmt.Errorf("failed to parse env var IDEMPOTENCY_TTL: %s", err)
	}

	reply, err := conn.Do("SET", key, "", "EX", ttl, "NX") // set with ttl and only if not exist

	if reply == nil {

		return errKeyExists
	}

	if err != nil {

		return err
	}

	return nil
}
