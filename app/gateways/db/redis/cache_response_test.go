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
			ID:        account.ID(uuid.New()),
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

	testKey := uuid.NewString()
	testKey2 := uuid.NewString()
	testKey3 := uuid.NewString()

	createdAccountJSON, _ := json.Marshal(testAccounts[0])
	fetchedAccountsJSON, _ := json.Marshal(testAccounts)
	createAccountErrorJSON, _ := json.Marshal(responses.ErrorPayload{Message: "Super descriptive error message"})

	testCases := []testCase{
		{
			name: "set a created account response",
			runBefore: func() {
				conn := testRepo.pool.Get()
				_, err := conn.Do("SET", testKey, "")
				if err != nil {
					t.Fatalf("failed to set key: %s", err)
				}
				conn.Close()
			},
			response: schema.CachedResponse{
				Key:            testKey,
				ResponseStatus: 201,
				ResponseBody:   createdAccountJSON,
			},
		},
		{
			name: "set a create account error response",
			runBefore: func() {
				conn := testRepo.pool.Get()
				_, err := conn.Do("SET", testKey2, "")
				if err != nil {
					t.Fatalf("failed to set key: %s", err)
				}
				conn.Close()
			},
			response: schema.CachedResponse{
				Key:            testKey2,
				ResponseStatus: 400,
				ResponseBody:   createAccountErrorJSON,
			},
		},
		{
			name: "set a fetched accounts slice response",
			runBefore: func() {
				conn := testRepo.pool.Get()
				_, err := conn.Do("SET", testKey3, "")
				if err != nil {
					t.Fatalf("failed to set key: %s", err)
				}
				conn.Close()
			},
			response: schema.CachedResponse{
				Key:            testKey3,
				ResponseStatus: 200,
				ResponseBody:   fetchedAccountsJSON,
			},
		},
		{
			name: "fail to set inexistent key",
			response: schema.CachedResponse{
				Key:            uuid.NewString(),
				ResponseStatus: 201,
			},
			err: errKeyNotFound,
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
