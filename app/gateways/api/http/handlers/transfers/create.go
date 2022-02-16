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

// @Summary Create transfer
// @Description Creates a transfer from origin account matching bearer token account id
// @Description to account with request body account ID
// @Tags Transfer
// @Accept json
// @Produce json
// @Param transfer body schemas.CreateTransferRequest true "Transfer"
// @Param Idempotency-Key header string false "Idempotency key"
// @Success 201 {object} schemas.CreateTransferResponse
// @Failure 400 {object} responses.ErrorPayload
// @Failure 401 {object} responses.ErrorPayload
// @Failure 404 {object} responses.ErrorPayload
// @Failure 409 {object} responses.ErrorPayload
// @Failure 422 {object} responses.ErrorPayload
// @Failure 500 {object} responses.ErrorPayload
// @Router /transfers [post]
func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	var createRequest schemas.CreateTransferRequest

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		response.BadRequest(responses.ErrInvalidCreateTransferPayload).SendJSON()

		return
	}

	originID, ok := middlewares.OriginIDFromContext(r.Context())

	if !ok {
		response.InternalServerError(errors.New("failed to parse origin id token")).SendJSON()

		return
	}

	transferInstance, response := createRequest.Validate(response, originID)

	if response.IsComplete() {
		response.SendJSON()

		return
	}

	persistedTransfer, err := h.usecase.Create(r.Context(), transferInstance)
	if err != nil {
		if errors.Is(err, account.ErrInsufficientFunds) {
			response.UnprocessableEntity(responses.ErrInsuficientFunds).SendJSON()

			return
		}

		if errors.Is(err, account.ErrIDNotFound) {
			response.NotFound(responses.ErrAccountNotFound).SendJSON()

			return
		}

		response.InternalServerError(err).SendJSON()

		return
	}

	response.Created(schemas.CreatedTransferToResponse(persistedTransfer)).SendJSON()
}
