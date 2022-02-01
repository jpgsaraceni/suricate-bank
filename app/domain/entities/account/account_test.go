package account

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()

	type arg struct {
		name   string
		cpf    string
		secret string
	}

	type testCase struct {
		name string
		args arg
		want Account
		err  error
	}

	var wantCpf = func(input string) cpf.Cpf {
		newCpf, _ := cpf.NewCpf(input)
		return newCpf
	}
	var testMoney, _ = money.NewMoney(1000)

	testCases := []testCase{
		{
			name: "doesn't create account because cpf is invalid",
			args: arg{
				name:   "Me",
				cpf:    "000.000.000-00",
				secret: "123456",
			},
			want: Account{},
			err:  ErrInvalidCpf,
		},
		{
			name: "creates new account",
			args: arg{
				name:   "Me",
				cpf:    "220.614.460-35",
				secret: "123456",
			},
			want: Account{
				Name:      "Me",
				Cpf:       wantCpf("220.614.460-35"),
				Balance:   testMoney,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "creates new account with 100 character secret",
			args: arg{
				name:   "Me",
				cpf:    "220.614.460-35",
				secret: "1hgfkljdbngdlgQT$34534621yGhtry5426%$53#$@%6bsgsdgtrtywy$#@%643grfGfdsvgarst%6t23se@745gdfsghT$E#24t",
			},
			want: Account{
				Name:      "Me",
				Cpf:       wantCpf("22061446035"),
				Balance:   testMoney,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "creates new account with 100 character name",
			args: arg{
				name:   "fajkldsjkl jalksfjasdlkads ajdsklghkjlahrwgfirpequfhaksljdgh ropuq trwj ewjfdsg opgthr wdhwfhsadgjkl",
				cpf:    "220.614.460-35",
				secret: "123456",
			},
			want: Account{
				Name:      "fajkldsjkl jalksfjasdlkads ajdsklghkjlahrwgfirpequfhaksljdgh ropuq trwj ewjfdsg opgthr wdhwfhsadgjkl",
				Cpf:       wantCpf("22061446035"),
				Balance:   testMoney,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "creates new account with 1 character secret",
			args: arg{
				name:   "Me",
				cpf:    "220.614.460-35",
				secret: "a",
			},
			want: Account{
				Name:      "Me",
				Cpf:       wantCpf("22061446035"),
				Balance:   testMoney,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "creates new account with 1 character name",
			args: arg{
				name:   "a",
				cpf:    "220.614.460-35",
				secret: "123456",
			},
			want: Account{
				Name:      "a",
				Cpf:       wantCpf("22061446035"),
				Balance:   testMoney,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "fails to create new account with empty name",
			args: arg{
				name:   "",
				cpf:    "220.614.460-35",
				secret: "123456",
			},
			want: Account{},
			err:  ErrEmptyName,
		},
		{
			name: "fails to create new account with empty secret",
			args: arg{
				name:   "Me",
				cpf:    "220.614.460-35",
				secret: "",
			},
			want: Account{},
			err:  ErrEmptySecret,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAccount(tt.args.name, tt.args.cpf, tt.args.secret)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			if got.Id != AccountId(uuid.Nil) {
				got.Id = tt.want.Id
			}

			got.Secret = tt.want.Secret
			got.CreatedAt = tt.want.CreatedAt

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v expected %v", got, tt.want)
			}
		})
	}
}
