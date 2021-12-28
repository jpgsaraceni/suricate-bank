package transferuc

import "errors"

var (
	ErrSameAccounts = errors.New("origin and destination must be different accounts")

	// errRepository is for mocking in tests
	errRepository = errors.New("repository error")
)
