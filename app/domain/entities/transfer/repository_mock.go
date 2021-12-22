package transfer

type MockRepository struct {
	OnCreate func(transfer *Transfer) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) Create(transfer *Transfer) error {
	return mock.OnCreate(transfer)
}
