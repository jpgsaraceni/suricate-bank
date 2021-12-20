package usecase

import "github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"

func (uc Usecase) GetBalance(id account.AccountId) (int, error) {
	balance, err := uc.Repository.GetBalance(id)

	if err != nil {

		return 0, ErrGetBalanceRepository
	}

	return balance, nil
}