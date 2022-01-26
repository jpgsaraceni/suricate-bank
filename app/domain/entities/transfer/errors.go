package transfer

import "errors"

var (
	ErrSameAccounts      = errors.New("origin and destination must be different accounts")
	ErrAmountNotPositive = errors.New("transfer amount must be greater than zero")
)
