package redis

import (
	"encoding/json"
	"errors"
	"reflect"
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

func TestGetCachedResponse(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)

	t.Cleanup(tearDown)

	type testCase struct {
		name      string
		runBefore func()
		key       string
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

	createdAccountJson, _ := json.Marshal(testAccounts[0])
	badRequestJson, _ := json.Marshal(responses.ErrorPayload{Message: "Great and meaningful error message"})

	testKey := uuid.NewString()
	testKey2 := uuid.NewString()

	testCases := []testCase{
		{
			name: "get a created account response",
			key:  testKey,
			runBefore: func() {
				testRepo.CacheResponse(schema.CachedResponse{
					Key:            testKey,
					ResponseStatus: 201,
					ResponseBody:   createdAccountJson,
				})
			},
			response: schema.CachedResponse{
				Key:            testKey,
				ResponseStatus: 201,
				ResponseBody:   createdAccountJson,
			},
		},
		{
			name: "get an error response",
			key:  testKey2,
			runBefore: func() {
				testRepo.CacheResponse(schema.CachedResponse{
					Key:            testKey2,
					ResponseStatus: 400,
					ResponseBody:   badRequestJson,
				})
			},
			response: schema.CachedResponse{
				Key:            testKey2,
				ResponseStatus: 400,
				ResponseBody:   badRequestJson,
			},
		},
		{
			name:     "respond empty when trying to get key that does not exist",
			key:      uuid.NewString(),
			response: schema.CachedResponse{},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				tt.runBefore()
			}

			got, err := testRepo.GetCachedResponse(tt.key)

			if !errors.Is(err, tt.err) {

				t.Fatalf("\ngot error: \n%s\nexpected error: \n%s\n", err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.response) {

				t.Fatalf("got %v expected %v", got, tt.response)
			}
		})
	}
}
