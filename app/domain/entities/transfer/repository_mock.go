package transfer

import "context"

type MockRepository struct {
	OnCreate func(ctx context.Context, transfer *Transfer) error
	OnFetch  func(ctx context.Context) ([]Transfer, error)
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) Create(ctx context.Context, transfer *Transfer) error {
	return mock.OnCreate(ctx, transfer)
}

func (mock MockRepository) Fetch(ctx context.Context) ([]Transfer, error) {
	return mock.OnFetch(ctx)
}
