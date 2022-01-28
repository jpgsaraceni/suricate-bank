package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

type FetchTransfersResponse struct {
	Transfers []FetchedTransfer `json:"transfers"`
}

type FetchedTransfer struct {
	Id                   string    `json:"transfer_id"`
	AccountOriginId      string    `json:"account_origin_id"`
	AccountDestinationId string    `json:"account_destination_id"`
	Amount               string    `json:"amount"`
	CreatedAt            time.Time `json:"created_at"`
}

func TransfersToResponse(transferList []transfer.Transfer) FetchTransfersResponse {
	transferResponse := make([]FetchedTransfer, 0, len(transferList))
	for _, transfer := range transferList {
		transferResponse = append(transferResponse, FetchedTransfer{
			Id:                   transfer.Id.String(),
			AccountOriginId:      transfer.AccountOriginId.String(),
			AccountDestinationId: transfer.AccountDestinationId.String(),
			Amount:               transfer.Amount.BRL(),
			CreatedAt:            transfer.CreatedAt,
		})
	}
	return FetchTransfersResponse{Transfers: transferResponse}
}
