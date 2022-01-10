package transfer

import "context"

type Repository interface {
	Create(ctx context.Context, transfer *Transfer) error
	Fetch(ctx context.Context) ([]Transfer, error)
}
