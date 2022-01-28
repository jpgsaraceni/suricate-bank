package accountuc

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name            string
		repository      account.Repository
		accountInstance account.Account
		want            account.Account
		err             error
	}

	testAccount := account.Account{
		Id:        account.AccountId(uuid.New()),
		Name:      "cool name",
		Cpf:       cpf.Random(),
		CreatedAt: time.Now(),
	}

	testCases := []testCase{
		{
			name: "successfully create account",
			repository: account.MockRepository{
				OnCreate: func(ctx context.Context, accountInstance account.Account) (account.Account, error) {
					return testAccount, nil
				},
			},
			want: testAccount,
		},
		{
			name: "fail to create new account with cpf that already exists",
			repository: account.MockRepository{
				OnCreate: func(ctx context.Context, accountInstance account.Account) (account.Account, error) {
					return account.Account{}, account.ErrDuplicateCpf
				},
			},
			want: account.Account{},
			err:  account.ErrDuplicateCpf,
		},
		{
			name: "repository throws error",
			repository: account.MockRepository{
				OnCreate: func(ctx context.Context, accountInstance account.Account) (account.Account, error) {
					return account.Account{}, errors.New("")
				},
			},
			want: account.Account{},
			err:  ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository}

			newAccount, err := uc.Create(context.Background(), tt.accountInstance)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			tt.want.Secret = newAccount.Secret

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
