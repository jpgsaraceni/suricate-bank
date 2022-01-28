package transfersroute

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/middlewares"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
)

func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	var createRequest schemas.CreateTransferRequest

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		response.BadRequest(responses.ErrInvalidCreateTransferPayload).SendJSON()

		return
	}

	originId, ok := middlewares.OriginIdFromContext(r.Context())

	if !ok {
		response.InternalServerError(errors.New("failed to parse origin id token")).SendJSON()

		return
	}

	transferInstance, response := createRequest.Validate(response, originId)

	if response.IsComplete() {
		response.SendJSON()

		return
	}

	err := h.usecase.Create(r.Context(), transferInstance)

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

	response.Created(schemas.CreatedTransferToResponse(transferInstance)).SendJSON()
}
