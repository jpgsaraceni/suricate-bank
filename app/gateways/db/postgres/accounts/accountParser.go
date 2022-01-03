package accountspg

import (
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

type queryReturn struct {
	id        uuid.UUID
	name      string
	cpf       string
	secret    string
	balance   int
	createdAt time.Time
}

func (q *queryReturn) parse() (account.Account, error) {
	var parsedAccount account.Account

	parsedAccount.Id = account.AccountId(q.id)
	parsedAccount.Name = q.name
	parsedAccount.Cpf, _ = cpf.NewCpf(q.cpf)
	parsedAccount.Balance, _ = money.NewMoney(q.balance)
	parsedAccount.Secret = hash.Parse(q.secret)
	parsedAccount.CreatedAt = q.createdAt

	return parsedAccount, nil
}
