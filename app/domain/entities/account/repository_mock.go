package account

type MockRepository struct {
	OnCreate func(account *Account) error
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) Create(account *Account) error {
	return mock.OnCreate(account)
}
