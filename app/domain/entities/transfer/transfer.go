package transfer

import (
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type (
	TransferId uuid.UUID
)

type Transfer struct {
	Id                   TransferId
	AccountOriginId      account.AccountId
	AccountDestinationId account.AccountId
	Amount               money.Money
	CreatedAt            time.Time
}

func NewTransfer(amount money.Money, originId, destinationId account.AccountId) (Transfer, error) {
	newTransfer := Transfer{
		Id:                   newTransferId(),
		AccountOriginId:      originId,
		AccountDestinationId: destinationId,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}

	return newTransfer, nil
}

func newTransferId() TransferId {

	return TransferId(uuid.New())
}
