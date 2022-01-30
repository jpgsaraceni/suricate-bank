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

func NewAccount(name string, cpfInput string, secret string) (Account, error) {
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

	newMoney, _ := money.NewMoney(1000)

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

func (id AccountId) String() string {
	parsedToUUID := uuid.UUID(id)
	return parsedToUUID.String()
}

func ParseAccountId(id string) (AccountId, error) {
	accountId, err := uuid.Parse(id)

	if err != nil {
		return AccountId{}, err
	}

	return AccountId(accountId), nil
}
