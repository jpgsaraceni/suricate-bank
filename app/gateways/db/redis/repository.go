package redis

import "github.com/gomodule/redigo/redis"

type Repository struct {
	pool *redis.Pool
}

func NewRepository(pool *redis.Pool) *Repository {
	return &Repository{pool}
}
