package redis

import (
	"context"
)

func (r Repository) Lock(ctx context.Context, key string) error {
	conn := r.pool.Get()
	defer conn.Close()

	reply, err := conn.Do("SETNX", key, "")

	if reply.(int64) == 0 {

		return errKeyExists
	}

	if err != nil {

		return err
	}

	return nil
}
