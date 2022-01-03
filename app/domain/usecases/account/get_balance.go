package accountuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc Usecase) GetBalance(ctx context.Context, id account.AccountId) (int, error) {
	balance, err := uc.Repository.GetBalance(ctx, id)

	if err != nil {

		return 0, fmt.Errorf("%w: %s", ErrGetBalance, err.Error())
	}

	return balance, nil
}
