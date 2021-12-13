package entities

import (
	"reflect"
	"testing"

	"github.com/google/uuid"
)

func TestNewAccount(t *testing.T) {
	t.Parallel()

	type arg struct {
		name   string
		cpf    Cpf
		secret string
	}

	type testCase struct {
		name        string
		args        arg
		want        Account
		expectedErr error
	}

	var TestId = AccountId(uuid.New())

	testCases := []testCase{
		{
			name: "doesn't create account because cpf is invalid",
			args: arg{
				name:   "Me",
				cpf:    Cpf("000.000.000-00"),
				secret: "123456",
			},
			want: Account{
				Id: TestId,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAccount(tt.args.name, tt.args.cpf, tt.args.secret)

			if tt.expectedErr != err {
				t.Errorf("got error %v espected error %v", err, tt.expectedErr)
			}

			if got.Id != AccountId(uuid.Nil) {
				got.Id = tt.want.Id
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v expected %v", got, tt.want)
			}
		})
	}
}
