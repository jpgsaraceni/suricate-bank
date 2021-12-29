package transferuc

import "errors"

var (
	ErrSameAccounts      = errors.New("origin and destination must be different accounts")
	ErrDebitOrigin       = errors.New("failed to debit origin account")
	ErrCreditDestination = errors.New("failed to credit detination account")
	ErrCreateTransfer    = errors.New("failed to save transfer")
)
