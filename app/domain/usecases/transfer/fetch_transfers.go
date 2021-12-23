package transferuc

import "github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"

func (uc Usecase) Fetch() ([]transfer.Transfer, error) {
	transferList, err := uc.Repository.Fetch()

	if err != nil {

		return []transfer.Transfer{}, errFetchTransfers
	}

	if len(transferList) == 0 {

		return transferList, errNoTransfersToFetch
	}

	return transferList, nil
}
