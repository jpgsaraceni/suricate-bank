package usecase

import (
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

	var testUUID, _ = uuid.NewUUID()

	testCases := []testCase{
		{
			name: "get 0 balance",
			repository: account.MockRepository{
				OnGetBalance: func(id account.AccountId) (int, error) {

					return 0, nil
				},
			},
			id:   account.AccountId(testUUID),
			want: 0,
		},
		{
			name: "get 100 balance",
			repository: account.MockRepository{
				OnGetBalance: func(id account.AccountId) (int, error) {

					return 100, nil
				},
			},
			id:   account.AccountId(testUUID),
			want: 100,
		},
		{
			name: "repository throws error",
			repository: account.MockRepository{
				OnGetBalance: func(id account.AccountId) (int, error) {

					return 0, ErrGetBalanceRepository
				},
			},
			id:   account.AccountId(uuid.Nil),
			want: 0,
			err:  ErrGetBalanceRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository}

			newAccount, err := uc.GetBalance(tt.id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
