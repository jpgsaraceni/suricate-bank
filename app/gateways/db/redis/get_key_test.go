package redis

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis/redistest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestGetKeyValue(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)

	t.Cleanup(tearDown)

	type testCase struct {
		name    string
		mustSet bool
		key     string
		res     responses.Response
		err     error
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
			name:    "get a response",
			mustSet: true,
			key:     "somekey",
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

			if tt.mustSet {
				testRepo.SetKeyValue(tt.key, tt.res)
			}

			got, err := testRepo.GetKeyValue(tt.key)

			if !errors.Is(err, tt.err) {

				t.Fatalf("\ngot error: \n%s\nexpected error: \n%s\n", err, tt.err)
			}

			if !reflect.DeepEqual(got.Payload, tt.res.Payload) || got.Status != tt.res.Status {

				t.Fatalf("got %v expected %v", got, tt.res)
			}
		})
	}
}
