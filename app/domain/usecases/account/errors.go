package accountuc

import "errors"

var (
	errCreditAccountRepository = errors.New("repository error when crediting account")
	errDebitAccountRepository  = errors.New("repository error when debiting account")
	errCreateAccount           = errors.New("failed to create account")
	errNameLength              = errors.New("invalid name length")
	errShortSecret             = errors.New("invalid password length")
	errCreateAccountRepository = errors.New("repository error")
	errGetBalanceRepository    = errors.New("failed to get balance")
	errFetchAccounts           = errors.New("couldn't fetch accounts")
	errNoAccountsToFetch       = errors.New("no accounts to fetch")
	errAccountNotFound         = errors.New("account not found")
	errNotPositive             = errors.New("amount must be greater than zero")
	errInsuficientFunds        = errors.New("insuficient funds in balance")
)
