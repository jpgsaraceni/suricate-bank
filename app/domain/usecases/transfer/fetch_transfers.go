package transferuc

import (
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func (uc Usecase) Fetch() ([]transfer.Transfer, error) {
	transferList, err := uc.Repository.Fetch()

	if err != nil {

		return []transfer.Transfer{}, fmt.Errorf("%w: %s", ErrFetchTransfers, err.Error())
	}

	if len(transferList) == 0 {

		return transferList, ErrNoTransfersToFetch
	}

	return transferList, nil
}
