package accountuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func (uc Usecase) GetBalance(id account.AccountId) (int, error) {
	balance, err := uc.Repository.GetBalance(id)

	if err != nil {

		return 0, fmt.Errorf("repository failed to get account balance: %w", err)
	}

	return balance, nil
}
