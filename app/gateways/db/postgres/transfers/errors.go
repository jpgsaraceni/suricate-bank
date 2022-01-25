package transferspg

import "errors"

var (
	ErrCreateTransfer = errors.New("failed to create transfer in db")
	ErrFetch          = errors.New("failed to fetch transfers from db")

	ErrScanningRows = errors.New("failed to scan rows returned from query")
)
