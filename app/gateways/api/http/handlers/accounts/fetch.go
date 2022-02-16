package accountsroute

import (
	"net/http"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
)

func (h handler) Fetch(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	accountList, err := h.usecase.Fetch(r.Context())
	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Ok(schemas.AccountsToResponse(accountList)).SendJSON()
}
