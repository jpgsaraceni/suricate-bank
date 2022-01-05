package transferspg

import (
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type queryReturn struct {
	id            uuid.UUID
	originId      uuid.UUID
	destinationId uuid.UUID
	amount        int
	createdAt     time.Time
}

// parse parses values returned from queries into types expected by transfer entity
func (q *queryReturn) parse() (transfer.Transfer, error) {
	var parsedTransfer transfer.Transfer

	parsedTransfer.Id = transfer.TransferId(q.id)
	parsedTransfer.AccountOriginId = account.AccountId(q.originId)
	parsedTransfer.AccountDestinationId = account.AccountId(q.destinationId)
	parsedTransfer.Amount, _ = money.NewMoney(q.amount)
	parsedTransfer.CreatedAt = q.createdAt

	return parsedTransfer, nil
}
