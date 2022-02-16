package auth

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
)

func (s service) Authenticate(ctx context.Context, cpfInput, secret string) (string, error) {
	validatedCpf, err := cpf.NewCpf(cpfInput)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrCredentials, err)
	}

	account, err := s.repository.GetByCpf(ctx, validatedCpf)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrCredentials, err)
	}

	if !account.Secret.Compare(secret) {
		return "", fmt.Errorf("%w: %s", ErrCredentials, err)
	}

	jwt, err := token.Sign(account.ID)
	if err != nil {
		return "", fmt.Errorf("%w: %s", ErrSignToken, err)
	}

	return jwt.Value(), nil
}
