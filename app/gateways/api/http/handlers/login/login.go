package loginroute

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
)

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	var loginRequest schemas.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		response.BadRequest(responses.ErrInvalidLoginPayload).SendJSON()

		return
	}

	account, err := h.service.Authenticate(r.Context(), loginRequest.Cpf, loginRequest.Secret)

	if errors.Is(err, auth.ErrCredentials) {
		response.Unauthorized(responses.ErrCredentials).SendJSON()

		return
	}

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Ok(schemas.LoginToResponse(account)).SendJSON()
}
