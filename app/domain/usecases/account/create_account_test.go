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

	var (
		testAccountId = account.AccountId(uuid.New())
		testCpf       = cpf.Random()
		testTime      = time.Now()
	)

	testCases := []testCase{
		{
			name: "successfully create account",
			repository: account.MockRepository{
				OnCreate: func(ctx context.Context, createdAccount *account.Account) error {
					createdAccount.Id = testAccountId
					createdAccount.CreatedAt = testTime
					return nil
				},
			},
			args: args{
				name:   "meee",
				cpf:    testCpf.Masked(),
				secret: "123456",
			},
			want: account.Account{
				Name:      "meee",
				Cpf:       testCpf,
				Id:        testAccountId,
				CreatedAt: testTime,
			},
		},
		{
			name: "fail to create account because password is too short",
			args: args{
				name:   "meee",
				cpf:    "220.614.460-35",
				secret: "123",
			},
			want: account.Account{},
			err:  ErrShortSecret,
		},
		{
			name: "fail to create account because name is too short",
			args: args{
				name:   "me",
				cpf:    "220.614.460-35",
				secret: "123",
			},
			want: account.Account{},
			err:  ErrNameLength,
		},
		{
			name: "fail to create account because NewAccount returned error",
			args: args{
				name:   "meee",
				cpf:    "220.614.4",
				secret: "123456",
			},
			want: account.Account{},
			err:  account.ErrInvalidCpf,
		},
		{
			name: "fail to create new account with cpf that already exists",
			repository: account.MockRepository{
				OnCreate: func(ctx context.Context, createdAccount *account.Account) error {
					return account.ErrDuplicateCpf
				},
			},
			args: args{
				name:   "meee",
				cpf:    "220.614.460-35",
				secret: "reallygoodpassphrase",
			},
			want: account.Account{},
			err:  account.ErrDuplicateCpf,
		},
		{
			name: "creates new account but Repository throws error",
			repository: account.MockRepository{
				OnCreate: func(ctx context.Context, createdAccount *account.Account) error {
					return errors.New("")
				},
			},
			args: args{
				name:   "meee",
				cpf:    "220.614.460-35",
				secret: "reallygoodpassphrase",
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

			newAccount, err := uc.Create(context.Background(), tt.args.name, tt.args.cpf, tt.args.secret)

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
