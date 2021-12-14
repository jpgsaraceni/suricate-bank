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
			err:  errCpf,
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
				Secret:    "123456",
				Balance:   0,
				CreatedAt: time.Now(),
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAccount(tt.args.name, tt.args.cpf, tt.args.secret)

			if tt.err != err {
				t.Errorf("got error %v espected error %v", err, tt.err)

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
