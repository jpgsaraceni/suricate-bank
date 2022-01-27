package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type CreateTransferRequest struct {
	AccountDestinationId string `json:"account_destination_id"`
	Amount               int    `json:"amount"`
}

type CreateTransferResponse struct {
	Id                   string    `json:"transfer_id"`
	AccountOriginId      string    `json:"account_origin_id"`
	AccountDestinationId string    `json:"account_destination_id"`
	Amount               string    `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

func CreatedTransferToResponse(createdTransfer transfer.Transfer) CreateTransferResponse {
	return CreateTransferResponse{
		Id:                   createdTransfer.Id.String(),
		AccountOriginId:      createdTransfer.AccountOriginId.String(),
		AccountDestinationId: createdTransfer.AccountDestinationId.String(),
		Amount:               createdTransfer.Amount.BRL(),
		CreatedAt:            createdTransfer.CreatedAt,
	}
}

func (r CreateTransferRequest) Validate(response *responses.Response, originId account.AccountId) (transfer.Transfer, bool) {
	if r.AccountDestinationId == "" || r.Amount == 0 {

		*response = response.BadRequest(responses.ErrMissingFieldsTransferPayload)
		return transfer.Transfer{}, false
	}

	if r.Amount < 0 {

		*response = response.BadRequest(responses.ErrInvalidAmount)
		return transfer.Transfer{}, false
	}

	amount, err := money.NewMoney(r.Amount)

	if err != nil {

		*response = response.InternalServerError(err)
		return transfer.Transfer{}, false
	}

	destinationId, err := account.ParseAccountId(r.AccountDestinationId)

	if err != nil {

		*response = response.BadRequest(responses.ErrInvalidDestinationId)
		return transfer.Transfer{}, false
	}

	if destinationId == originId {

		*response = response.BadRequest(responses.ErrSameAccounts)
		return transfer.Transfer{}, false
	}

	transferInstance, err := transfer.NewTransfer(amount, originId, destinationId)

	if err != nil {

		*response = response.InternalServerError(err)
		return transfer.Transfer{}, false
	}

	return transferInstance, true
}
