package transferspg

import "github.com/jackc/pgx/v4/pgxpool"

type Repository struct {
	pool *pgxpool.Pool
}

func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool}
}
