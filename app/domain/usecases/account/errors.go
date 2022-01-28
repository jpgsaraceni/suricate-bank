package accountuc

import (
	"errors"
)

var (
	ErrNameLength  = errors.New("invalid name length")
	ErrShortSecret = errors.New("invalid password length")
	ErrAmount      = errors.New("invalid amount")

	ErrRepository = errors.New("accounts repository error")
)
