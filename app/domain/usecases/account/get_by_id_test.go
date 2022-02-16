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
		id         account.ID
		want       account.Account
		err        error
	}

	testAccountID := account.ID(uuid.New())

	testCases := []testCase{
		{
			name: "get account",
			repository: account.MockRepository{
				OnGetByID: func(ctx context.Context, id account.ID) (account.Account, error) {
					return account.Account{
						ID: id,
					}, nil
				},
			},
			id: testAccountID,
			want: account.Account{
				ID: testAccountID,
			},
		},
		{
			name: "fail to get account",
			repository: account.MockRepository{
				OnGetByID: func(ctx context.Context, id account.ID) (account.Account, error) {
					return account.Account{}, errors.New("")
				},
			},
			id:   testAccountID,
			want: account.Account{},
			err:  ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository}

			newAccount, err := uc.GetByID(context.Background(), tt.id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
