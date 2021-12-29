package transferuc

import "errors"

var (
	ErrSameAccounts = errors.New("origin and destination must be different accounts")
)
