package entities

import (
	"errors"
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"github.com/jpgsaraceni/suricate-bank/app/cpf"
)

type (
	AccountId uuid.UUID
	Cpf       string
	Secret    string
	Money     int
)

type Account struct {
	Id        AccountId
	Name      string
	Cpf       string
	Secret    Secret
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

	hashedSecret, err := createHash(secret)

	if err != nil {

		return Account{}, errHash
	}

	return Account{
		Id:        newAccountId(),
		Name:      name,
		Cpf:       cpf.Value(),
		Secret:    hashedSecret,
		Balance:   0,
		CreatedAt: time.Now(),
	}, nil
}

const hashCost = 10

func createHash(secret string) (Secret, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), hashCost)

	if err != nil {

		return "", err
	}

	return Secret(hash), nil
}

func newAccountId() AccountId {
	return AccountId(uuid.New())
}
