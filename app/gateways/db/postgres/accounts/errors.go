package accountspg

import "errors"

var (
	ErrQuery            = errors.New("failed to run query")
	ErrScanningRows     = errors.New("failed to scan rows returned from query")
	ErrParse            = errors.New("failed to parse query return")
	ErrEmptyFetch       = errors.New("fetch returned empty")
	ErrBalanceUnchanged = errors.New("debit failed because amount to debit is greater than balance")
	ErrCpfAlreadyExists = errors.New("cpf already exists in db")
	ErrCpfNotFound      = errors.New("cpf does not exist in db")
)
