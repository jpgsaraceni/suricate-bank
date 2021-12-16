package usecase

import (
	"reflect"
	"testing"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		name   string
		cpf    string
		secret string
	}

	type testCase struct {
		name string
		uc   account.Repository
		args args
		want account.Account
		err  error
	}

	testCases := []testCase{
		{
			name: "successfully create account",
			uc: account.MockRepository{
				OnCreate: func(account *account.Account) error {
					return nil
				},
			},
			args: args{
				name:   "meee",
				cpf:    "220.614.460-35",
				secret: "reallygoodpassphrase",
			},
			want: account.Account{
				Name: "meee",
				Cpf:  "22061446035",
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.uc}

			newAccount, err := uc.Create(tt.args.name, tt.args.cpf, tt.args.secret)

			if err != tt.err {
				t.Errorf("got %s expected %s", err, tt.err)
			}

			newAccount.Id = tt.want.Id
			newAccount.Secret = tt.want.Secret
			newAccount.CreatedAt = tt.want.CreatedAt

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
