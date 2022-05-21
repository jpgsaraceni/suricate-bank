package redis

import (
	"context"

	"github.com/jpgsaraceni/suricate-bank/config"
)

func (r Repository) Lock(_ context.Context, cfg config.Config, key string) error {
	conn := r.pool.Get()
	defer conn.Close()
	ttl := cfg.IdempotencyKeyTTL

	reply, err := conn.Do("SET", key, "", "EX", ttl, "NX") // set with ttl and only if not exist

	if reply == nil {
		return errKeyExists
	}

	if err != nil {
		return err
	}

	return nil
}
