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

func (h handler) GetBalance(w http.ResponseWriter, r *http.Request) error {
	p := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(p[2])

	var response responses.Response
	response.Writer = w

	if err != nil {

		return response.BadRequest(responses.ErrInvalidPathParameter).SendJSON()
	}

	balance, err := h.usecase.GetBalance(r.Context(), account.AccountId(id))

	if errors.Is(err, accountuc.ErrIdNotFound) {

		return response.BadRequest(responses.ErrAccountNotFound).SendJSON()
	}

	if err != nil {

		return response.InternalServerError(err).SendJSON()
	}

	return response.Ok(responses.GotBalancePayload(balance)).SendJSON()
}
