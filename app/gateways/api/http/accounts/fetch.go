package accountsroute

import (
	"encoding/json"
	"errors"
	"net/http"

	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func (h handler) Fetch(w http.ResponseWriter, r *http.Request) {
	var response responses.Response
	response.Writer = w

	accountList, err := h.usecase.Fetch(r.Context())

	if errors.Is(err, accountuc.ErrNoAccountsToFetch) { // TODO: reavaliar se deve ser um erro mesmo
		response.Ok(responses.NoAccounts).SendJSON()

		return
	}

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	j, err := json.Marshal(accountList)

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}
	response.Ok(responses.FetchedAccountsPayload(j)).SendJSON()
}
