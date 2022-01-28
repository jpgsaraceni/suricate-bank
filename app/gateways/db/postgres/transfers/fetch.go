package transferspg

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func (r Repository) Fetch(ctx context.Context) ([]transfer.Transfer, error) {
	const query = `
		SELECT 
			id,
			account_origin_id,
			account_destination_id,
			amount,
			created_at
		FROM transfers;
	`

	rows, err := r.pool.Query(ctx, query)

	if err != nil {

		return nil, fmt.Errorf("%w: %s", ErrFetch, err.Error())
	}

	defer rows.Close()
	var transferList []transfer.Transfer

	for rows.Next() {
		var transferReturned transfer.Transfer
		err := rows.Scan(
			&transferReturned.Id,
			&transferReturned.AccountOriginId,
			&transferReturned.AccountDestinationId,
			&transferReturned.Amount,
			&transferReturned.CreatedAt,
		)

		if err != nil {

			return nil, fmt.Errorf("%w: %s", ErrScanningRows, err.Error())
		}

		transferList = append(transferList, transferReturned)
	}

	return transferList, nil
}
