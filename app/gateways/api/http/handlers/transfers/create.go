package transfersroute

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/middlewares"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	var createRequest schemas.CreateTransferRequest

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		response.BadRequest(responses.ErrInvalidCreateTransferPayload).SendJSON()

		return
	}

	if createRequest.AccountDestinationId == "" || createRequest.Amount == 0 {
		response.BadRequest(responses.ErrMissingFieldsTransferPayload).SendJSON()

		return
	}

	if createRequest.Amount < 0 {
		response.BadRequest(responses.ErrInvalidAmount).SendJSON()
	}

	originId, ok := middlewares.OriginIdFromContext(r.Context())

	if !ok {
		response.InternalServerError(errors.New("failed to parse origin id token")).SendJSON()

		return
	}

	amount, err := money.NewMoney(createRequest.Amount)

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	destinationId, err := account.ParseAccountId(createRequest.AccountDestinationId)

	if err != nil {
		response.BadRequest(responses.ErrInvalidDestinationId).SendJSON()

		return
	}

	if destinationId == originId {
		response.BadRequest(responses.ErrSameAccounts).SendJSON()

		return
	}

	createdTransfer, err := h.usecase.Create(r.Context(), amount, originId, destinationId)

	if err != nil {
		if errors.Is(err, account.ErrInsufficientFunds) {
			response.UnprocessableEntity(responses.ErrInsuficientFunds).SendJSON()

			return
		}

		if errors.Is(err, account.ErrIdNotFound) {
			response.NotFound(responses.ErrAccountNotFound).SendJSON()

			return
		}

		response.InternalServerError(err).SendJSON()

		return
	}

	response.Created(schemas.CreatedTransferToResponse(createdTransfer)).SendJSON()
}
