package accountuc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func TestGetById(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository account.Repository
		id         account.AccountId
		want       account.Account
		err        error
	}

	var testUUID, _ = uuid.NewUUID()

	var errRepository = errors.New("repository error")

	testCases := []testCase{
		{
			name: "get account",
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{
						Id: id,
					}, nil
				},
			},
			id: account.AccountId(testUUID),
			want: account.Account{
				Id: account.AccountId(testUUID),
			},
		},
		{
			name: "fail to get account",
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{}, errRepository
				},
			},
			id:   account.AccountId(testUUID),
			want: account.Account{},
			err:  ErrGetAccount,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository}

			newAccount, err := uc.GetById(context.Background(), tt.id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
