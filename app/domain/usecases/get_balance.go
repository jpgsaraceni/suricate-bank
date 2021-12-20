package usecase

import "github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"

func (uc Usecase) GetBalance(id account.AccountId) (int, error) {

	return uc.Repository.GetBalance(id)
}
