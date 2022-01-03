package accountuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc Usecase) Fetch(ctx context.Context) ([]account.Account, error) {
	accountList, err := uc.Repository.Fetch(ctx)

	if err != nil {

		return []account.Account{}, fmt.Errorf("%w: %s", ErrFetchAccounts, err.Error())
	}

	if len(accountList) == 0 {

		return accountList, ErrNoAccountsToFetch
	}

	return accountList, nil
}
