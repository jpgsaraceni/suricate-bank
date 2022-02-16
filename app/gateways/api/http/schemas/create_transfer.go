package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type CreateTransferRequest struct {
	AccountDestinationID string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

type CreateTransferResponse struct {
	ID                   string    `json:"transfer_id"`
	AccountOriginID      string    `json:"account_origin_id"`
	AccountDestinationID string    `json:"account_destination_id"`
	Amount               string    `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

func CreatedTransferToResponse(createdTransfer transfer.Transfer) CreateTransferResponse {
	return CreateTransferResponse{
		ID:                   createdTransfer.ID.String(),
		AccountOriginID:      createdTransfer.AccountOriginID.String(),
		AccountDestinationID: createdTransfer.AccountDestinationID.String(),
		Amount:               createdTransfer.Amount.BRL(),
		CreatedAt:            createdTransfer.CreatedAt,
	}
}

func (r CreateTransferRequest) Validate(response responses.Response, originID account.ID) (transfer.Transfer, responses.Response) {
	if r.AccountDestinationID == "" || r.Amount == 0 {
		return transfer.Transfer{}, response.BadRequest(responses.ErrMissingFieldsTransferPayload)
	}

	if r.Amount < 0 {
		return transfer.Transfer{}, response.BadRequest(responses.ErrInvalidAmount)
	}

	amount, err := money.NewMoney(r.Amount)
	if err != nil {
		return transfer.Transfer{}, response.InternalServerError(err)
	}

	destinationID, err := account.ParseAccountID(r.AccountDestinationID)
	if err != nil {
		return transfer.Transfer{}, response.BadRequest(responses.ErrInvalidDestinationID)
	}

	if destinationID == originID {
		return transfer.Transfer{}, response.BadRequest(responses.ErrSameAccounts)
	}

	transferInstance, err := transfer.NewTransfer(amount, originID, destinationID)
	if err != nil {
		return transfer.Transfer{}, response.InternalServerError(err)
	}

	return transferInstance, response
}
