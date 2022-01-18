package accountsroute

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
)

func (h handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	var response responses.Response
	response.Writer = w

	accountId, err := getAccountIdFromRequest(r)

	if err != nil {
		response.BadRequest(responses.ErrInvalidPathParameter).SendJSON()

		return
	}

	balance, err := h.usecase.GetBalance(r.Context(), accountId)

	if errors.Is(err, accountuc.ErrIdNotFound) {
		response.NotFound(responses.ErrAccountNotFound).SendJSON()

		return
	}

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	balanceJSON, err := schemas.BalanceToResponse(accountId, balance).Marshal()

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Payload = balanceJSON

	response.Ok().SendJSON()
}

func getAccountIdFromRequest(r *http.Request) (account.AccountId, error) {
	idParam := getPathParam(r.URL.Path)
	parsedToUuid, err := uuid.Parse(idParam)

	if err != nil {

		return account.AccountId{}, err
	}

	return account.AccountId(parsedToUuid), nil
}

func getPathParam(url string) string {
	parts := strings.Split(url, "/")
	position := getPathParamPosition()

	return parts[position]
}

const pathPattern = "/accounts/{account_id}/balance"

func getPathParamPosition() int {
	re := regexp.MustCompile(`\{\w*\}`)
	parts := strings.Split(pathPattern, "/")
	for i := range parts {
		if re.MatchString(parts[i]) {

			return i
		}
	}

	return 0
}
