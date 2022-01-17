package accountsroute

import (
	"errors"
	"net/http"
	"strings"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func (h handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var response responses.Response
	response.Writer = w

	accountId, err := getAccountIdFromPath(r)

	if err != nil {
		response.BadRequest(responses.ErrInvalidPathParameter).SendJSON()

		return
	}

	balance, err := h.usecase.GetBalance(r.Context(), accountId)

	if errors.Is(err, accountuc.ErrIdNotFound) {
		response.BadRequest(responses.ErrAccountNotFound).SendJSON()

		return
	}

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Ok(responses.GotBalancePayload(balance)).SendJSON()
}

func getAccountIdFromPath(r *http.Request) (account.AccountId, error) {
	// const pathParamPattern = "/accounts/{account_id}/balance"
	// idParam := getPathParam(r.URL.Path, pathParamPattern)
	idParam := getPathParam(r.URL.Path, 2)
	parsedToUuid, err := uuid.Parse(idParam)

	if err != nil {

		return account.AccountId{}, err
	}

	return account.AccountId(parsedToUuid), nil
}

func getPathParam(url string, position int) string { // TODO: define position using regexp
	parts := strings.Split(url, "/")
	return parts[position]
}
