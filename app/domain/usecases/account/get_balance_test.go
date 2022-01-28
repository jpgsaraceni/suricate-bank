package accountuc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository account.Repository
		id         account.AccountId
		want       int
		err        error
	}

	var testAccountId = account.AccountId(uuid.New())

	testCases := []testCase{
		{
			name: "get 0 balance",
			repository: account.MockRepository{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {

					return 0, nil
				},
			},
			id:   testAccountId,
			want: 0,
		},
		{
			name: "get 100 balance",
			repository: account.MockRepository{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {

					return 100, nil
				},
			},
			id:   testAccountId,
			want: 100,
		},
		{
			name: "repository throws error id not found",
			repository: account.MockRepository{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {

					return 0, account.ErrIdNotFound
				},
			},
			id:   account.AccountId(uuid.Nil),
			want: 0,
			err:  account.ErrIdNotFound,
		},
		{
			name: "repository throws some other error",
			repository: account.MockRepository{
				OnGetBalance: func(ctx context.Context, id account.AccountId) (int, error) {

					return 0, errors.New("")
				},
			},
			id:   account.AccountId(uuid.Nil),
			want: 0,
			err:  ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository}

			newAccount, err := uc.GetBalance(context.Background(), tt.id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
