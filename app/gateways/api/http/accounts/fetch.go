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

	defer func(r *responses.Response) {
		responses.SendJSON(w, *r)
	}(&response)

	accountList, err := h.usecase.Fetch(r.Context())

	if errors.Is(err, accountuc.ErrNoAccountsToFetch) { // TODO: reavaliar se deve ser um erro mesmo
		response = responses.Ok(responses.NoAccounts)

		return
	}

	if err != nil {
		response = responses.InternalServerError(err)

		return
	}

	j, err := json.Marshal(accountList)

	if err != nil {
		response = responses.InternalServerError(err)

		return
	}

	response = responses.Ok(responses.FetchedAccountsPayload(j))
}
