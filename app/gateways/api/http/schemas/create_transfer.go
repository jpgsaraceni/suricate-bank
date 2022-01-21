package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
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
