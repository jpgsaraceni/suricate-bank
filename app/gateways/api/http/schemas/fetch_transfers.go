package schemas

import (
	"time"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
)

type FetchTransfersResponse struct {
	Transfers []FetchedTransfer `json:"transfers"`
}

type FetchedTransfer struct {
	ID                   string    `json:"transfer_id" example:"d8e0810f-64aa-4b26-8eab-7dc2ebf0b02b"`
	AccountOriginID      string    `json:"account_origin_id" example:"fbba165f-0382-491d-8a83-b950cb6482c9"`
	AccountDestinationID string    `json:"account_destination_id" example:"5738eda2-49f5-4702-83e4-b87b18cf0d31"`
	Amount               string    `json:"amount" example:"R$1,00"`
	CreatedAt            time.Time `json:"created_at" example:"2022-01-28T19:39:04.585238-03:00"`
}

func TransfersToResponse(transferList []transfer.Transfer) FetchTransfersResponse {
	transferResponse := make([]FetchedTransfer, 0, len(transferList))
	for _, transfer := range transferList {
		transferResponse = append(transferResponse, FetchedTransfer{
			ID:                   transfer.ID.String(),
			AccountOriginID:      transfer.AccountOriginID.String(),
			AccountDestinationID: transfer.AccountDestinationID.String(),
			Amount:               transfer.Amount.BRL(),
			CreatedAt:            transfer.CreatedAt,
		})
	}

	return FetchTransfersResponse{Transfers: transferResponse}
}
