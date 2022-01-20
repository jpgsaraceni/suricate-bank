package auth

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func (s service) Authenticate(ctx context.Context, cpf cpf.Cpf, secret string) (string, error) {
	account, err := s.repository.GetByCpf(ctx, cpf)

	if err != nil {

		return "", fmt.Errorf("%w: %s", ErrInexistentCpf, err)
	}

	if !account.Secret.Compare(secret) {

		return "", fmt.Errorf("%w: %s", ErrWrongPassword, err)
	}

	// TODO: generate token using accountId
	return "", nil
}
