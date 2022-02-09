package redis

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis/redistest"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestCacheResponse(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)

	t.Cleanup(tearDown)

	type testCase struct {
		name      string
		runBefore func()
		response  schema.CachedResponse
		err       error
	}

	testAccount := func() account.Account {
		return account.Account{
			Id:        account.AccountId(uuid.New()),
			Name:      "nice name",
			Cpf:       cpf.Random(),
			Balance:   money.Money{},
			CreatedAt: time.Now(),
		}
	}

	testAccounts := []account.Account{
		testAccount(),
		testAccount(),
		testAccount(),
		testAccount(),
	}

	repeatedKey := uuid.NewString()

	createdAccountJson, _ := json.Marshal(testAccounts[0])
	fetchedAccountsJson, _ := json.Marshal(testAccounts)
	createAccountErrorJson, _ := json.Marshal(responses.ErrorPayload{Message: "Super descriptive error message"})

	testCases := []testCase{
		{
			name: "set a created account response",
			response: schema.CachedResponse{
				Key:            uuid.NewString(),
				ResponseStatus: 201,
				ResponseBody:   createdAccountJson,
			},
		},
		{
			name: "set a create account error response",
			response: schema.CachedResponse{
				Key:            uuid.NewString(),
				ResponseStatus: 400,
				ResponseBody:   createAccountErrorJson,
			},
		},
		{
			name: "set a fetched accounts slice response",
			response: schema.CachedResponse{
				Key:            uuid.NewString(),
				ResponseStatus: 200,
				ResponseBody:   fetchedAccountsJson,
			},
		},
		{
			name: "fail to set existent key",
			runBefore: func() {
				testRepo.CacheResponse(context.Background(), schema.CachedResponse{
					Key: repeatedKey,
				})
			},
			response: schema.CachedResponse{
				Key:            repeatedKey,
				ResponseStatus: 201,
			},
			err: errKeyExists,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				tt.runBefore()
			}

			err := testRepo.CacheResponse(context.Background(), tt.response)

			if !errors.Is(err, tt.err) {

				t.Fatalf("\ngot error: \n%s\nexpected error: \n%s\n", err, tt.err)
			}
		})
	}
}
