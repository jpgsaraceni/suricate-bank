package loginroute

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
)

// @Summary Login
// @Tags Login
// @Accept json
// @Produce json
// @Param credentials body schemas.LoginRequest true "Login Credentials"
// @Success 200 {object} schemas.LoginResponse
// @Failure 400 {object} responses.ErrorPayload
// @Failure 401 {object} responses.ErrorPayload
// @Failure 500 {object} responses.ErrorPayload
// @Router /login [post]
func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	var loginRequest schemas.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		response.BadRequest(responses.ErrInvalidLoginPayload).SendJSON()

		return
	}

	account, err := h.service.Authenticate(r.Context(), h.Config, loginRequest.Cpf, loginRequest.Secret)

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
