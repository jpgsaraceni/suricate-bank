package accountuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc Usecase) Fetch() ([]account.Account, error) {
	accountList, err := uc.Repository.Fetch()

	if err != nil {

		return []account.Account{}, fmt.Errorf("%w: %s", ErrFetchAccounts, err.Error())
	}

	if len(accountList) == 0 {

		return accountList, ErrNoAccountsToFetch
	}

	return accountList, nil
}
