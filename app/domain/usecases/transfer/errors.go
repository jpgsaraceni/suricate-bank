package transferuc

import "errors"

var (
	errSameAccounts = errors.New("origin and destination must be different accounts")
	errDebit        = errors.New("failed to debit from origin account")
	errCredit       = errors.New("failed to credit to destination account")
	// errRollback     = errors.New("failed to rollback after transfer error")
	errCreateTransfer           = errors.New("failed transfer")
	errCreateTransferRepository = errors.New("repository error")
	errFetchTransfers           = errors.New("couldn't fetch transfers")
	errNoTransfersToFetch       = errors.New("no transfers to fetch")
)
