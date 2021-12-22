package account

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type (
	AccountId uuid.UUID
)

type Account struct {
	Id        AccountId
	Name      string
	Cpf       cpf.Cpf
	Secret    hash.Secret
	Balance   money.Money
	CreatedAt time.Time
}

var (
	errInvalidCpf  = errors.New("invalid cpf")
	errNewHash     = errors.New("hash failed")
	errEmptyName   = errors.New("empty name")
	errEmptySecret = errors.New("empty secret")
	errCredit      = errors.New("failed to credit account")
	errDebit       = errors.New("failed to debit account")
)

func NewAccount(name string, cpfInput string, secret string) (Account, error) {
	if len(name) == 0 {

		return Account{}, errEmptyName
	}

	if len(secret) == 0 {

		return Account{}, errEmptySecret
	}

	cpf, err := cpf.NewCpf(cpfInput)

	if err != nil {

		return Account{}, errInvalidCpf
	}

	hashedSecret, err := hash.NewHash(secret)

	if err != nil {

		return Account{}, errNewHash
	}

	newMoney, _ := money.NewMoney(0)

	return Account{
		Id:        newAccountId(),
		Name:      name,
		Cpf:       cpf,
		Secret:    hashedSecret,
		Balance:   newMoney,
		CreatedAt: time.Now(),
	}, nil
}

func newAccountId() AccountId {

	return AccountId(uuid.New())
}
