package accountuc

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository account.Repository
		want       []account.Account
		err        error
	}

	var testUUID1, _ = uuid.NewUUID()
	var testUUID2, _ = uuid.NewUUID()
	var testUUID3, _ = uuid.NewUUID()
	var testUUID4, _ = uuid.NewUUID()
	var testCpf = func(input string) cpf.Cpf {
		newCpf, _ := cpf.NewCpf(input)
		return newCpf
	}

	var errRepository = errors.New("repository error")

	testCases := []testCase{
		{
			name: "successfully fetch 1 account",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{
						{
							Id:   account.AccountId(testUUID1),
							Name: "Account 1",
							Cpf:  testCpf("220.614.460-35"),
						},
					}, nil
				},
			},
			want: []account.Account{
				{
					Id:   account.AccountId(testUUID1),
					Name: "Account 1",
					Cpf:  testCpf("220.614.460-35"),
				},
			},
		},
		{
			name: "successfully fetch 4 accounts",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{
						{
							Id:   account.AccountId(testUUID1),
							Name: "Account 1",
							Cpf:  testCpf("220.614.460-35"),
						},
						{
							Id:   account.AccountId(testUUID2),
							Name: "Account 2",
							Cpf:  testCpf("232.598.190-88"),
						},
						{
							Id:   account.AccountId(testUUID3),
							Name: "Account 3",
							Cpf:  testCpf("816.413.860-61"),
						},
						{
							Id:   account.AccountId(testUUID4),
							Name: "Account 4",
							Cpf:  testCpf("924.498.310-96"),
						},
					}, nil
				},
			},
			want: []account.Account{
				{
					Id:   account.AccountId(testUUID1),
					Name: "Account 1",
					Cpf:  testCpf("220.614.460-35"),
				},
				{
					Id:   account.AccountId(testUUID2),
					Name: "Account 2",
					Cpf:  testCpf("232.598.190-88"),
				},
				{
					Id:   account.AccountId(testUUID3),
					Name: "Account 3",
					Cpf:  testCpf("816.413.860-61"),
				},
				{
					Id:   account.AccountId(testUUID4),
					Name: "Account 4",
					Cpf:  testCpf("924.498.310-96"),
				},
			},
		},
		{
			name: "no existent accounts error",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{}, nil
				},
			},
			want: []account.Account{},
			err:  ErrNoAccountsToFetch,
		},
		{
			name: "repository throws error",
			repository: account.MockRepository{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {

					return []account.Account{}, errRepository
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
