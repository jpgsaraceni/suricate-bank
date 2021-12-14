package usecase

// import (
// 	"reflect"
// 	"testing"

// 	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
// )

// func TestCreateAccount(t *testing.T) {
// 	t.Parallel()

// 	type arg struct {
// 		name   string
// 		cpf    string
// 		secret string
// 	}

// 	type testCase struct {
// 		name string
// 		args arg
// 		want account.Account
// 		err  error
// 	}

// 	testCases := []testCase{
// 		{
// 			name: "name is too short",
// 			args: arg{
// 				name:   "me",
// 				cpf:    "123.456.789-01",
// 				secret: "123456",
// 			},
// 			err: errName,
// 		},
// 	}

// 	for _, tt := range testCases {
// 		tt := tt
// 		t.Run(tt.name, func(t *testing.T) {
// 			t.Parallel()
// 			got, err := Usecase.Create(tt.args.name, tt.args.cpf, tt.args.secret)

// 			if err != tt.err {
// 				t.Errorf("got error %s expected error %s", err, tt.err)
// 			}

// 			if !reflect.DeepEqual(got, tt.want) {
// 				t.Errorf("got %v expected %v", got, tt.want)
// 			}
// 		})
// 	}
// }
