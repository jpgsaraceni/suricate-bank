package entities

import (
	"time"

	"github.com/google/uuid"
)

type ( // TODO consider moving these types out of entities
	Name   string
	Cpf    int
	Secret string
	Money  int
)

type Account struct {
	Id        uuid.UUID
	Name      Name
	Cpf       Cpf
	Secret    Secret
	Balance   Money
	CreatedAt time.Time
}

func (a Account) NewAccount(name Name, cpf Cpf, secret Secret) Account {
	return Account{
		Id:        uuid.New(),
		Name:      name,
		Cpf:       cpf,    // TODO validate this
		Secret:    secret, // TODO hash this
		Balance:   0,
		CreatedAt: time.Now(),
	}
}
