package account

import (
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
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

	testCases := []testCase{
		{
			name: "doesn't create account because cpf is invalid",
			args: arg{
				name:   "Me",
				cpf:    "000.000.000-00",
				secret: "123456",
			},
			want: Account{},
			err:  errInvalidCpf,
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
				Cpf:       "22061446035",
				Balance:   0,
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
				Cpf:       "22061446035",
				Balance:   0,
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
				Cpf:       "22061446035",
				Balance:   0,
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
				Cpf:       "22061446035",
				Balance:   0,
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
				Cpf:       "22061446035",
				Balance:   0,
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
			err:  errEmptyName,
		},
		{
			name: "fails to create new account with empty secret",
			args: arg{
				name:   "Me",
				cpf:    "220.614.460-35",
				secret: "",
			},
			want: Account{},
			err:  errEmptySecret,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAccount(tt.args.name, tt.args.cpf, tt.args.secret)

			if err != tt.err {
				t.Errorf("got error %v expected error %v", err, tt.err)

				return
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
