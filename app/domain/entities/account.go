package entities

import (
	"time"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type ( // TODO consider moving these types out of entities
	AccountId uuid.UUID
	Cpf       string // TODO get from cpf package
	Secret    string
	Money     int
)

type Account struct {
	Id        AccountId
	Name      string
	Cpf       Cpf
	Secret    Secret
	Balance   Money
	CreatedAt time.Time
}

func NewAccount(name string, cpf Cpf, secret string) (Account, error) {
	hashedSecret, err := createHash(secret)

	if err != nil {

		return Account{}, err
	}

	return Account{
		Id:        AccountId(uuid.New()),
		Name:      name,
		Cpf:       cpf, // TODO validate this
		Secret:    hashedSecret,
		Balance:   0,
		CreatedAt: time.Now(),
	}, nil
}

func createHash(secret string) (Secret, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(secret), 10)

	if err != nil {

		return "", err
	}

	return Secret(hash), nil
}
