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

var errCpf = errors.New("invalid cpf")
var errHash = errors.New("hash failed")

func NewAccount(name string, cpfInput string, secret string) (Account, error) {
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
