package account

type Repository interface {
	Create(account *Account) error
	GetBalance(id AccountId) (int, error)
	Fetch() ([]Account, error)
	GetById(id AccountId) (Account, error)
}
