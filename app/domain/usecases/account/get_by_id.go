package accountuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc Usecase) GetById(ctx context.Context, id account.AccountId) (account.Account, error) {
	account, err := uc.Repository.GetById(ctx, id)

	if err != nil {

		return account, fmt.Errorf("%w: %s", ErrGetAccount, err.Error())
	}

	return account, nil
}
