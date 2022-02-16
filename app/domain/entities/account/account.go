package account

import (
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type (
	ID uuid.UUID
)

type Account struct {
	ID        ID
	Name      string
	Cpf       cpf.Cpf
	Secret    hash.Secret
	Balance   money.Money
	CreatedAt time.Time
}

func NewAccount(name, cpfInput, secret string) (Account, error) {
	if len(name) == 0 {
		return Account{}, ErrEmptyName
	}

	if len(secret) == 0 {
		return Account{}, ErrEmptySecret
	}

	cpf, err := cpf.NewCpf(cpfInput)
	if err != nil {
		return Account{}, ErrInvalidCpf
	}

	hashedSecret, err := hash.NewHash(secret)
	if err != nil {
		return Account{}, fmt.Errorf("failed to hash secret: %w", err)
	}

	const testMoneyAmount = 1000
	newMoney, _ := money.NewMoney(testMoneyAmount)

	return Account{
		ID:        newAccountID(),
		Name:      name,
		Cpf:       cpf,
		Secret:    hashedSecret,
		Balance:   newMoney,
		CreatedAt: time.Now(),
	}, nil
}

func newAccountID() ID {
	return ID(uuid.New())
}

func (id ID) String() string {
	parsedToUUID := uuid.UUID(id)

	return parsedToUUID.String()
}

func ParseAccountID(id string) (ID, error) {
	accountID, err := uuid.Parse(id)
	if err != nil {
		return ID{}, err
	}

	return ID(accountID), nil
}
