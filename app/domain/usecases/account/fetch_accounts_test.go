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
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository account.Repository
		want       []account.Account
		err        error
	}

	var (
		testAccount1 = account.Account{
			Id:        account.AccountId(uuid.New()),
			Name:      "Nice name",
			Cpf:       cpf.Random(),
			Secret:    hash.Parse("123456"),
			CreatedAt: time.Now(),
		}
		testAccount2 = account.Account{
			Id:        account.AccountId(uuid.New()),
			Name:      "Nice name",
			Cpf:       cpf.Random(),
			Secret:    hash.Parse("123456"),
			CreatedAt: time.Now(),
		}
		testAccount3 = account.Account{
			Id:        account.AccountId(uuid.New()),
			Name:      "Nice name",
			Cpf:       cpf.Random(),
			Secret:    hash.Parse("123456"),
			CreatedAt: time.Now(),
		}
		testAccount4 = account.Account{
			Id:        account.AccountId(uuid.New()),
			Name:      "Nice name",
			Cpf:       cpf.Random(),
			Secret:    hash.Parse("123456"),
			CreatedAt: time.Now(),
		}
	)

	testCases := []testCase{
		{
			name: "successfully fetch 1 account",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{
						testAccount1,
					}, nil
				},
			},
			want: []account.Account{
				testAccount1,
			},
		},
		{
			name: "successfully fetch 4 accounts",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{
						testAccount1,
						testAccount2,
						testAccount3,
						testAccount4,
					}, nil
				},
			},
			want: []account.Account{
				testAccount1,
				testAccount2,
				testAccount3,
				testAccount4,
			},
		},
		{
			name: "successfully fetch 0 accounts",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{}, nil
				},
			},
			want: []account.Account{},
		},
		{
			name: "repository throws error",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{}, errors.New("")
				},
			},
			want: []account.Account{},
			err:  ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository}

			accountList, err := uc.Fetch(context.Background())

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(accountList, tt.want) {
				t.Errorf("got %v expected %v", accountList, tt.want)
			}
		})
	}
}
