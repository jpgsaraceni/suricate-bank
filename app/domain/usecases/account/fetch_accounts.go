package accountuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc Usecase) Fetch() ([]account.Account, error) {
	accountList, err := uc.Repository.Fetch()

	if err != nil {

		return []account.Account{}, fmt.Errorf("failed to fetch accounts: %w", err)
	}

	if len(accountList) == 0 {

		return accountList, ErrNoAccountsToFetch
	}

	return accountList, nil
}
