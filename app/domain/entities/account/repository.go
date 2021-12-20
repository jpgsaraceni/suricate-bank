package account

type Repository interface {
	Create(account *Account) error
}
