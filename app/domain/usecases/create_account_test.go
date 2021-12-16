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

	var mockCreateSuccess = func(account *account.Account) error {
		account.Id = testAccountId
		account.CreatedAt = testTime
		account.Secret = "hashedpassphrase"
		return nil
	}

	testCases := []testCase{
		{
			name: "successfully create account",
			repository: account.MockRepository{
				OnCreate: mockCreateSuccess,
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
		{
			name: "fail to create account because password is too short",
			repository: account.MockRepository{
				OnCreate: mockCreateSuccess,
			},
			args: args{
				name:   "meee",
				cpf:    "220.614.460-35",
				secret: "123",
			},
			want: account.Account{},
			err:  errShortSecret,
		},
		{
			name: "fail to create account because name is too short",
			repository: account.MockRepository{
				OnCreate: mockCreateSuccess,
			},
			args: args{
				name:   "me",
				cpf:    "220.614.460-35",
				secret: "123",
			},
			want: account.Account{},
			err:  errNameLength,
		},
		{
			name: "fail to create account because NewAccount returned error",
			repository: account.MockRepository{
				OnCreate: mockCreateSuccess,
			},
			args: args{
				name:   "meee",
				cpf:    "220.614.4",
				secret: "123456",
			},
			want: account.Account{},
			err:  errCreate,
		},
		{
			name: "creates new account but Repository throws error",
			repository: account.MockRepository{
				OnCreate: func(account *account.Account) error {
					account.Id = testAccountId
					account.CreatedAt = testTime
					account.Secret = "hashedpassphrase"
					return errRepository
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
			err: errRepository,
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
