package accountsroute

import (
	"net/http"
	"strings"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
)

func (h handler) GetBalance(w http.ResponseWriter, r *http.Request) {
	p := strings.Split(r.URL.Path, "/")
	id, err := uuid.Parse(p[1])

	if err != nil {
		// bad request
		return
	}

	balance, err := h.usecase.GetBalance(r.Context(), account.AccountId(id))

	if err != nil {
		// id doesn't exist or internal server error
		return
	}

	responses.SendJSON(w, responses.GotBalance(balance))
}
