package usecase

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository account.Repository
		want       []account.Account
		err        error
	}

	var testUUID, _ = uuid.NewUUID()
	var testCpf = func(input string) cpf.Cpf {
		newCpf, _ := cpf.NewCpf(input)
		return newCpf
	}

	testCases := []testCase{
		{
			name: "successfully fetch 1 account",
			repository: account.MockRepository{
				OnFetch: func() ([]account.Account, error) {

					return []account.Account{
						{
							Id:      account.AccountId(testUUID),
							Name:    "Account 1",
							Cpf:     testCpf("220.614.460-35"),
							Balance: 10,
						},
					}, nil
				},
			},
			want: []account.Account{
				{
					Id:      account.AccountId(testUUID),
					Name:    "Account 1",
					Cpf:     testCpf("220.614.460-35"),
					Balance: 10,
				},
			},
		},
		{
			name: "no existent accounts error",
			repository: account.MockRepository{
				OnFetch: func() ([]account.Account, error) {

					return []account.Account{}, nil
				},
			},
			want: []account.Account{},
			err:  ErrNoAccountsToFetch,
		},
		{
			name: "repository throws error",
			repository: account.MockRepository{
				OnFetch: func() ([]account.Account, error) {

					return []account.Account{}, ErrFetchAccounts
				},
			},
			want: []account.Account{},
			err:  ErrFetchAccounts,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository}

			newAccount, err := uc.Fetch()

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
