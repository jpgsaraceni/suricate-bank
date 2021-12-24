package accountspg

import "errors"

var (
	errQuery        = errors.New("failed to run query")
	errScanningRows = errors.New("failed to scan rows returned from query")
)
