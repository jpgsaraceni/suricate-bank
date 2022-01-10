package transferspg

import "errors"

var (
	ErrQuery        = errors.New("failed to run query")
	ErrScanningRows = errors.New("failed to scan rows returned from query")
	ErrEmptyFetch   = errors.New("fetch returned empty")
	ErrParse        = errors.New("failed to parse query return")
)
