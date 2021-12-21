package accountuc

import "github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"

func (uc Usecase) GetById(id account.AccountId) (account.Account, error) {
	account, err := uc.Repository.GetById(id)

	if err != nil {

		return account, ErrGetBalanceRepository
	}

	return account, nil
}
