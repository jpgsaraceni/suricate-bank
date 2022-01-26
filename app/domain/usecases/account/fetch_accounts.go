package accountuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc usecase) Fetch(ctx context.Context) ([]account.Account, error) {
	accountList, err := uc.repository.Fetch(ctx)

	if err != nil {

		return []account.Account{}, fmt.Errorf("%w: %s", ErrRepository, err.Error())
	}

	return accountList, nil
}
