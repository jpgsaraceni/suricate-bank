package accountsroute

import (
	"encoding/json"
	"errors"
	"net/http"

	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func (h handler) Fetch(w http.ResponseWriter, r *http.Request) error {
	var response responses.Response
	response.Writer = w

	accountList, err := h.usecase.Fetch(r.Context())

	if errors.Is(err, accountuc.ErrNoAccountsToFetch) { // TODO: reavaliar se deve ser um erro mesmo

		return response.Ok(responses.NoAccounts).SendJSON()
	}

	if err != nil {

		return response.InternalServerError(err).SendJSON()
	}

	j, err := json.Marshal(accountList)

	if err != nil {

		return response.InternalServerError(err).SendJSON()
	}

	return response.Ok(responses.FetchedAccountsPayload(j)).SendJSON()
}
