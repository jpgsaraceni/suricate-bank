package accountspg

import "errors"

var (
	ErrQuery        = errors.New("failed to run query")
	errScanningRows = errors.New("failed to scan rows returned from query")
	errParse        = errors.New("failed to parse query return")
)
