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
	p := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(p[1])

	var response responses.Response

	defer func(r *responses.Response) {
		responses.SendJSON(w, *r)
	}(&response)

	if err != nil {
		response = responses.BadRequest(responses.ErrInvalidPathParameter)

		return
	}

	balance, err := h.usecase.GetBalance(r.Context(), account.AccountId(id))

	if errors.Is(err, accountuc.ErrIdNotFound) {
		response = responses.BadRequest(responses.ErrAccountNotFound)

		return
	}

	if err != nil {
		response = responses.InternalServerError(err)

		return
	}

	response = responses.GotBalance(balance)
}
