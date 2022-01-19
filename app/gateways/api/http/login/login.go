package loginroute

import (
	"encoding/json"
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
)

func (h handler) Login(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	var loginRequest schemas.LoginRequest

	if err := json.NewDecoder(r.Body).Decode(&loginRequest); err != nil {
		response.BadRequest(responses.ErrInvalidRequestPayload).SendJSON()

		return
	}

	// TODO: check cpf, call authenticate service and write response
}
