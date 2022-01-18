package accountsroute

import (
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
)

func (h handler) Fetch(w http.ResponseWriter, r *http.Request) {
	var response responses.Response
	response.Writer = w

	accountList, err := h.usecase.Fetch(r.Context())

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	accountListResponse := schemas.AccountsToResponse(accountList)

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Payload = accountListResponse

	response.Ok().SendJSON()
}
