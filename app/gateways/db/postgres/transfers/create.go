package transferspg

import (
	"context"
	"fmt"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

func (r Repository) Create(ctx context.Context, transfer *transfer.Transfer) error {

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
			($1, $2, $3, $4, $5);
	`

	_, err := r.pool.Exec(
		ctx,
		query,
		transfer.Id,
		transfer.AccountOriginId,
		transfer.AccountDestinationId,
		transfer.Amount.Cents(),
		transfer.CreatedAt,
	)

	if err != nil {

		return fmt.Errorf("%w: %s", ErrQuery, err.Error())
	}

	return nil
}
