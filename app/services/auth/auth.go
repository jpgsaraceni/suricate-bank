package auth

import (
	"context"
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func (s service) Authenticate(ctx context.Context, cpf cpf.Cpf, secret string) (string, error) {
	account, err := s.repository.GetByCpf(ctx, cpf)

	if err != nil {

		return "", errors.New("inexistent cpf") // TODO: create error variable
	}

	if !account.Secret.Compare(secret) {

		return "", errors.New("wrong password") // TODO: create error variable
	}

	// generate token using accountId
	return "", nil
}
