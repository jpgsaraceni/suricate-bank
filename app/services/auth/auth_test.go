package auth

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	type args struct {
		cpf    cpf.Cpf
		secret string
	}

	type testCase struct {
		name       string
		repository account.Repository
		args       args
		want       account.AccountId
		err        error
	}

	var (
		testId        = account.AccountId(uuid.New())
		testCpf       = cpf.Random()
		testSecret, _ = hash.NewHash("123456")
	)

	testCases := []testCase{
		{
			name: "successfully authenticate request",
			repository: account.MockRepository{
				OnGetByCpf: func(ctx context.Context, cpf cpf.Cpf) (account.Account, error) {
					return account.Account{
						Id:     testId,
						Cpf:    testCpf,
						Secret: testSecret,
					}, nil
				},
			},
			args: args{
				cpf:    testCpf,
				secret: "123456",
			},
			want: testId,
		},
		{
			name: "fail to authenticate request with inexistent cpf",
			repository: account.MockRepository{
				OnGetByCpf: func(ctx context.Context, cpf cpf.Cpf) (account.Account, error) {
					return account.Account{}, ErrCredentials
				},
			},
			args: args{
				cpf:    testCpf,
				secret: "123456",
			},
			want: account.AccountId{},
			err:  ErrCredentials,
		},
		{
			name: "fail to authenticate request with wrong password",
			repository: account.MockRepository{
				OnGetByCpf: func(ctx context.Context, cpf cpf.Cpf) (account.Account, error) {
					return account.Account{
						Id:     testId,
						Cpf:    testCpf,
						Secret: testSecret,
					}, nil
				},
			},
			args: args{
				cpf:    testCpf,
				secret: "12345",
			},
			want: account.AccountId{},
			err:  ErrCredentials,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := service{tt.repository}

			gotToken, err := s.Authenticate(context.Background(), tt.args.cpf.Value(), tt.args.secret)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			gotId, _ := token.Verify(gotToken)

			if !reflect.DeepEqual(gotId, tt.want) {
				t.Errorf("got %v expected %v", gotId, tt.want)
			}
		})
	}
}
