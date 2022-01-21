package transfersroute

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
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

	if createRequest.Amount <= 0 {
		response.BadRequest(responses.ErrInvalidAmount).SendJSON()
	}

	authHeader := r.Header.Get("Authorization")
	requestToken := strings.ReplaceAll(authHeader, "Bearer ", "")

	if requestToken == "" {
		response.Unauthorized(responses.ErrMissingAuthorizationHeader).SendJSON()

		return
	}

	originId, err := token.Verify(requestToken)

	if err != nil {
		response.Forbidden(responses.ErrInvalidToken).SendJSON()

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

	createdTransfer, err := h.usecase.Create(r.Context(), amount, originId, destinationId)

	if errors.Is(err, transferuc.ErrSameAccounts) {
		response.BadRequest(responses.ErrSameAccounts).SendJSON()

		return
	}

	if errors.Is(err, money.ErrInsuficientFunds) {
		response.UnprocessableEntity(responses.ErrInsuficientFunds).SendJSON()

		return
	}

	// TODO: inexistent origin or destination account

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Created(schemas.CreatedTransferToResponse(createdTransfer)).SendJSON()
}
