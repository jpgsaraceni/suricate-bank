package account

import (
	"errors"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/hash"
)

type (
	AccountId uuid.UUID
	Money     int
)

type Account struct {
	Id        AccountId
	Name      string
	Cpf       string
	Secret    string
	Balance   Money
	CreatedAt time.Time
}

var (
	errCpf    = errors.New("invalid cpf")
	errHash   = errors.New("hash failed")
	errName   = errors.New("empty name")
	errSecret = errors.New("empty name")
)

func NewAccount(name string, cpfInput string, secret string) (Account, error) {
	if len(name) == 0 {

		return Account{}, errName
	}

	if len(secret) == 0 {

		return Account{}, errSecret
	}

	cpf, err := cpf.NewCpf(cpfInput)

	if err != nil {

		return Account{}, errCpf
	}

	hashedSecret, err := hash.NewHash(secret)

	if err != nil {

		return Account{}, errHash
	}

	return Account{
		Id:        newAccountId(),
		Name:      name,
		Cpf:       cpf.Value(),
		Secret:    hashedSecret.Value(),
		Balance:   0,
		CreatedAt: time.Now(),
	}, nil
}

func newAccountId() AccountId {
	return AccountId(uuid.New())
}
