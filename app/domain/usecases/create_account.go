package usecase

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities"
)

func Create(name string, cpf entities.Cpf, secret string) (entities.Account, error) {
	a, err := entities.NewAccount(name, cpf, secret)

	if err != nil {
		return a, err
	}

	return a, nil
}
