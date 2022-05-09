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
	"github.com/jpgsaraceni/suricate-bank/config"
)

func TestAuthenticate(t *testing.T) {
	t.Parallel()

	cfg := config.Config{
		JwtConfig: config.JwtConfig{
			JWTSecret: "whatever",
		},
	}

	type args struct {
		cpf    cpf.Cpf
		secret string
	}

	type testCase struct {
		name       string
		repository account.Repository
		args       args
		want       account.ID
		err        error
	}

	var (
		testID        = account.ID(uuid.New())
		testCpf       = cpf.Random()
		testSecret, _ = hash.NewHash("123456")
	)

	testCases := []testCase{
		{
			name: "successfully authenticate request",
			repository: account.MockRepository{
				OnGetByCpf: func(ctx context.Context, cpf cpf.Cpf) (account.Account, error) {
					return account.Account{
						ID:     testID,
						Cpf:    testCpf,
						Secret: testSecret,
					}, nil
				},
			},
			args: args{
				cpf:    testCpf,
				secret: "123456",
			},
			want: testID,
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
			want: account.ID{},
			err:  ErrCredentials,
		},
		{
			name: "fail to authenticate request with wrong password",
			repository: account.MockRepository{
				OnGetByCpf: func(ctx context.Context, cpf cpf.Cpf) (account.Account, error) {
					return account.Account{
						ID:     testID,
						Cpf:    testCpf,
						Secret: testSecret,
					}, nil
				},
			},
			args: args{
				cpf:    testCpf,
				secret: "12345",
			},
			want: account.ID{},
			err:  ErrCredentials,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := service{tt.repository}

			gotToken, err := s.Authenticate(context.Background(), cfg, tt.args.cpf.Value(), tt.args.secret)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}

			gotID, _ := token.Verify(cfg, gotToken)

			if !reflect.DeepEqual(gotID, tt.want) {
				t.Errorf("got %v expected %v", gotID, tt.want)
			}
		})
	}
}
