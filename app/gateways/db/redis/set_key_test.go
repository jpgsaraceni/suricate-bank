package redis

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis/redistest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestSetKeyValue(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)

	t.Cleanup(tearDown)

	type testCase struct {
		name string
		key  string
		res  responses.Response
		err  error
	}

	testAccount := account.Account{
		Id:        account.AccountId(uuid.New()),
		Name:      "nice name",
		Cpf:       cpf.Random(),
		Balance:   money.Money{},
		CreatedAt: time.Now(),
	}

	testCases := []testCase{
		{
			name: "set a response",
			key:  "somekey",
			res: responses.Response{
				Status: 200,
				Payload: map[string]interface{}{
					"account_id": testAccount.Id.String(),
					"name":       testAccount.Name,
					"cpf":        testAccount.Cpf.Masked(),
					"balance":    testAccount.Balance.BRL(),
					"created_at": testAccount.CreatedAt.Format(time.RFC3339Nano),
				},
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := testRepo.SetKeyValue(tt.key, tt.res)

			if !errors.Is(err, tt.err) {

				t.Fatalf("\ngot error: \n%s\nexpected error: \n%s\n", err, tt.err)
			}
		})
	}
}
