package accountuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc usecase) GetByID(ctx context.Context, id account.ID) (account.Account, error) {
	account, err := uc.repository.GetByID(ctx, id)
	if err != nil {
		return account, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return account, nil
}
