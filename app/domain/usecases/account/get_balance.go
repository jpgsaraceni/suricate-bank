package accountuc

import (
	"context"
	"errors"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc usecase) GetBalance(ctx context.Context, id account.AccountId) (int, error) {
	balance, err := uc.repository.GetBalance(ctx, id)

	if errors.Is(err, account.ErrIdNotFound) {

		return 0, err
	}

	if err != nil {

		return 0, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return balance, nil
}
