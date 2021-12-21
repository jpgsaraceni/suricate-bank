package accountuc

import (
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

	testCases := []testCase{
		{
			name: "get account",
			repository: account.MockRepository{
				OnGetById: func(id account.AccountId) (account.Account, error) {
					return account.Account{
						Id: account.AccountId(testUUID),
					}, nil
				},
			},
			id: account.AccountId(testUUID),
			want: account.Account{
				Id: account.AccountId(testUUID),
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository}

			newAccount, err := uc.GetById(tt.id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
