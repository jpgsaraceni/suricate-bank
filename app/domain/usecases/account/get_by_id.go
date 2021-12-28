package accountuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc Usecase) GetById(id account.AccountId) (account.Account, error) {
	account, err := uc.Repository.GetById(id)

	if err != nil {

		return account, fmt.Errorf("repository failed to get account: %w", err)
	}

	return account, nil
}
