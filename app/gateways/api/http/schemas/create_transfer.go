package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type CreateTransferRequest struct {
	AccountDestinationID string `json:"account_destination_id" example:"5738eda2-49f5-4702-83e4-b87b18cf0d31"`
	Amount               int    `json:"amount" example:"100"`
}

type CreateTransferResponse struct {
	ID                   string    `json:"transfer_id" example:"d8e0810f-64aa-4b26-8eab-7dc2ebf0b02b"`
	AccountOriginID      string    `json:"account_origin_id" example:"fbba165f-0382-491d-8a83-b950cb6482c9"`
	AccountDestinationID string    `json:"account_destination_id" example:"5738eda2-49f5-4702-83e4-b87b18cf0d31"`
	Amount               string    `json:"amount" example:"R$1,00"`
	CreatedAt            time.Time `json:"created_at" example:"2022-01-28T19:39:04.585238-03:00"`
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
