package transfer

import "errors"

var (
	ErrSameAccounts = errors.New("origin and destination must be different accounts")
	ErrAmountZero   = errors.New("cannot transfer zero money")
)
