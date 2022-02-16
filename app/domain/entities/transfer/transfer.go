package transfer

import (
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type (
	ID uuid.UUID
)

type Transfer struct {
	ID                   ID
	AccountOriginID      account.ID
	AccountDestinationID account.ID
	Amount               money.Money
	CreatedAt            time.Time
}

func NewTransfer(amount money.Money, originID, destinationID account.ID) (Transfer, error) {
	if originID == destinationID {
		return Transfer{}, ErrSameAccounts
	}

	if amount.Cents() <= 0 {
		return Transfer{}, ErrAmountNotPositive
	}

	newTransfer := Transfer{
		ID:                   newTransferID(),
		AccountOriginID:      originID,
		AccountDestinationID: destinationID,
		Amount:               amount,
		CreatedAt:            time.Now(),
	}

	return newTransfer, nil
}

func newTransferID() ID {
	return ID(uuid.New())
}

func (id ID) String() string {
	parsedToUUID := uuid.UUID(id)

	return parsedToUUID.String()
}
