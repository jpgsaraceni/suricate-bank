package idempotency

import (
	"context"
	"encoding/json"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestGetCachedResponse(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository MockRepository
		key        string
		want       schema.CachedResponse
		err        error
	}

	testKey := uuid.NewString()
	testAccount := account.Account{
		ID:        account.ID(uuid.New()),
		Name:      "nice name",
		Cpf:       cpf.Random(),
		Balance:   money.Money{},
		CreatedAt: time.Now(),
	}
	createdAccountJSON, _ := json.Marshal(testAccount)

	testCases := []testCase{
		{
			name: "get a created account response",
			key:  testKey,
			repository: MockRepository{
				OnGetCachedResponse: func(context.Context, string) (schema.CachedResponse, error) {
					return schema.CachedResponse{
						Key:            testKey,
						ResponseStatus: 200,
						ResponseBody:   createdAccountJSON,
					}, nil
				},
			},
			want: schema.CachedResponse{
				Key:            testKey,
				ResponseStatus: 200,
				ResponseBody:   createdAccountJSON,
			},
		},
		{
			name: "repository error",
			key:  testKey,
			repository: MockRepository{
				OnGetCachedResponse: func(context.Context, string) (schema.CachedResponse, error) {
					return schema.CachedResponse{}, errors.New("repository uh-oh")
				},
			},
			err: ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewService(tt.repository)

			got, err := s.GetCachedResponse(context.Background(), tt.key)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %s expected error %s", err, tt.err)
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Fatalf("got %v wanted %v", got, tt.want)
			}
		})
	}
}
