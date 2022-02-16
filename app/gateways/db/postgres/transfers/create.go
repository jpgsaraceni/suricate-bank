package transferspg

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func (r Repository) Create(ctx context.Context, transferInstance transfer.Transfer) (transfer.Transfer, error) {
	const query = `
		INSERT INTO
			transfers (
				id,
				account_origin_id,
				account_destination_id,
				amount,
				created_at
			)
		VALUES
			($1, $2, $3, $4, $5)
		RETURNING 
			id,
			account_origin_id,
			account_destination_id,
			amount,
			created_at;
	`

	var transferReturned transfer.Transfer

	err := r.pool.QueryRow(
		ctx,
		query,
		transferInstance.ID,
		transferInstance.AccountOriginID,
		transferInstance.AccountDestinationID,
		transferInstance.Amount.Cents(),
		transferInstance.CreatedAt,
	).Scan(
		&transferReturned.ID,
		&transferReturned.AccountOriginID,
		&transferReturned.AccountDestinationID,
		&transferReturned.Amount,
		&transferReturned.CreatedAt,
	)
	if err != nil {
		return transfer.Transfer{}, fmt.Errorf("%w: %s", ErrCreateTransfer, err.Error())
	}

	return transferReturned, nil
}
