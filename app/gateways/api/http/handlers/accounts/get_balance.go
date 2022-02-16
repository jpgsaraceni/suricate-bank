package accountsroute

import (
	"errors"
	"net/http"
	"regexp"
	"strings"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/schemas"
)

func (h handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	response := responses.NewResponse(w)

	accountID, err := getAccountIDFromRequest(r)
	if err != nil {
		response.BadRequest(responses.ErrInvalidPathParameter).SendJSON()

		return
	}

	balance, err := h.usecase.GetBalance(r.Context(), accountID)

	if errors.Is(err, account.ErrIDNotFound) {
		response.NotFound(responses.ErrAccountNotFound).SendJSON()

		return
	}

	if err != nil {
		response.InternalServerError(err).SendJSON()

		return
	}

	response.Ok(schemas.BalanceToResponse(accountID, balance)).SendJSON()
}

func getAccountIDFromRequest(r *http.Request) (account.ID, error) {
	idParam := getPathParam(r.URL.Path)
	parsedToUUID, err := uuid.Parse(idParam)
	if err != nil {
		return account.ID{}, err
	}

	return account.ID(parsedToUUID), nil
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
