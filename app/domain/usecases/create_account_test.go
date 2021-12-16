package usecase

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
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
		name       string
		repository account.Repository
		args       args
		want       account.Account
		err        error
	}

	var testUUID, _ = uuid.NewUUID()
	var testAccountId = account.AccountId(testUUID)
	var testTime = time.Now()

	testCases := []testCase{
		{
			name: "successfully create account",
			repository: account.MockRepository{
				OnCreate: func(account *account.Account) error {
					account.Id = testAccountId
					account.CreatedAt = testTime
					account.Secret = "hashedpassphrase"
					return nil
				},
			},
			args: args{
				name:   "meee",
				cpf:    "220.614.460-35",
				secret: "reallygoodpassphrase",
			},
			want: account.Account{
				Name:      "meee",
				Cpf:       "22061446035",
				Id:        testAccountId,
				CreatedAt: testTime,
				Secret:    "hashedpassphrase",
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository}

			newAccount, err := uc.Create(tt.args.name, tt.args.cpf, tt.args.secret)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			if !reflect.DeepEqual(newAccount, tt.want) {
				t.Errorf("got %v expected %v", newAccount, tt.want)
			}
		})
	}
}
