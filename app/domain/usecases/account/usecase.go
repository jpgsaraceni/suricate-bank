package accountuc

import (
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

// Usecase calls Repository to be used in all methods of this package.
type Usecase struct {
	Repository account.Repository
}
