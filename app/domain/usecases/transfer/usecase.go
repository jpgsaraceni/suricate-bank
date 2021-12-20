package transferuc

import (
	"errors"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

type Usecase struct {
	Repository transfer.Repository
}

var (
	ErrCreateTransfer           = errors.New("failed transfer")
	ErrCreateTransferRepository = errors.New("repository error")
)
