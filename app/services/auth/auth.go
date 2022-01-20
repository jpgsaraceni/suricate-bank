package auth

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
)

func (s service) Authenticate(ctx context.Context, cpf cpf.Cpf, secret string) (string, error) {
	account, err := s.repository.GetByCpf(ctx, cpf)

	if err != nil {

		return "", fmt.Errorf("%w: %s", ErrInexistentCpf, err)
	}

	if !account.Secret.Compare(secret) {

		return "", fmt.Errorf("%w: %s", ErrWrongPassword, err)
	}

	jwt, err := token.Sign(account.Id)

	if err != nil {

		return "", fmt.Errorf("%w: %s", ErrSignToken, err)
	}

	return jwt.Value(), nil
}
