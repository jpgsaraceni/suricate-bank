package accountspg

import "errors"

var (
	ErrQuery        = errors.New("failed to run query")
	ErrScanningRows = errors.New("failed to scan rows returned from query")
)
