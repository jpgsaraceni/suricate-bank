package accountsroute

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
)

// @Summary Create a new account
// @Description Creates an account with BRL$10.00 initial balance
// @Tags Account
// @Accept json
// @Produce json
// @Param account body schemas.CreateAccountRequest true "Account"
// @Param Idempotency-Key header string false "Idempotency key"
// @Success 201 {object} schemas.CreateAccountResponse
// @Failure 400 {object} responses.ErrorPayload
// @Failure 404 {object} responses.ErrorPayload
// @Failure 409 {object} responses.ErrorPayload
// @Failure 500 {object} responses.ErrorPayload
// @Router /accounts [post]
func (h handler) Create(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	var createRequest schemas.CreateAccountRequest

	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		response.BadRequest(responses.ErrInvalidCreateAccountPayload).SendJSON()

		return
	}

	accountInstance, response := createRequest.Validate(response)

	if response.IsComplete() {
		response.SendJSON()

		return
	}

	persistedAccount, err := h.usecase.Create(r.Context(), accountInstance)

	if errors.Is(err, account.ErrInvalidCpf) {
		response.BadRequest(responses.ErrInvalidCpf).SendJSON()

		return
	}

	if errors.Is(err, account.ErrDuplicateCpf) {
		response.BadRequest(responses.ErrCpfAlreadyExists).SendJSON()

		return
	}

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Created(schemas.CreatedAccountToResponse(persistedAccount)).SendJSON()
}
