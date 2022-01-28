package accountuc

import (
	"context"
	"errors"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc usecase) Create(ctx context.Context, accountInstance account.Account) (account.Account, error) {
	persistedAccount, err := uc.repository.Create(ctx, accountInstance)

	if err != nil {
		if errors.Is(err, account.ErrDuplicateCpf) {

			return account.Account{}, err
		}

		return account.Account{}, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return persistedAccount, nil
}
