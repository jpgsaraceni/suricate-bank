package transferuc

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func (uc usecase) Fetch(ctx context.Context) ([]transfer.Transfer, error) {
	transferList, err := uc.Repository.Fetch(ctx)

	if err != nil {

		return []transfer.Transfer{}, fmt.Errorf("%w: %s", ErrFetchTransfers, err.Error())
	}

	if len(transferList) == 0 {

		return transferList, ErrNoTransfersToFetch
	}

	return transferList, nil
}
