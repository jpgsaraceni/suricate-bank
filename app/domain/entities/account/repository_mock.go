package account

type MockRepository struct {
	OnCreate     func(account *Account) error
	OnGetBalance func(id AccountId) (int, error)
}

var _ Repository = (*MockRepository)(nil)

func (mock MockRepository) Create(account *Account) error {
	return mock.OnCreate(account)
}

func (mock MockRepository) GetBalance(id AccountId) (int, error) {
	return mock.OnGetBalance(id)
}
