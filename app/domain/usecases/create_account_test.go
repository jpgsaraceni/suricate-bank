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

	uc := Usecase{
		account.MockRepository{
			OnCreate: func(account *account.Account) error {
				return nil
			},
		}}

	type testCase struct {
		name string
		uc   Usecase
		args args
		want account.Account
		err  error
	}

	testCases := []testCase{
		{
			name: "successfully create account",
			uc:   uc,
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

			newAccount, err := tt.uc.Create(tt.args.name, tt.args.cpf, tt.args.secret)

			if err != tt.err {
				t.Errorf("got %s expected %s", err, tt.err)
			}

			err = tt.uc.Repository.Create(&newAccount)

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
