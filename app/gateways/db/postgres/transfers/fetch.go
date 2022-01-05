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

		return nil, fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	defer rows.Close()
	var transferList []transfer.Transfer

	for rows.Next() {
		var transferReturned queryReturn
		err := rows.Scan(
			&transferReturned.id,
			&transferReturned.originId,
			&transferReturned.destinationId,
			&transferReturned.amount,
			&transferReturned.createdAt,
		)

		if err != nil {

			return nil, fmt.Errorf("%w: %s", ErrScanningRows, err.Error())
		}

		parsedTransfer, err := transferReturned.parse()

		if err != nil {

			return nil, fmt.Errorf("%w: %s", ErrParse, err.Error())
		}

		transferList = append(transferList, parsedTransfer)
	}

	return transferList, nil
}
