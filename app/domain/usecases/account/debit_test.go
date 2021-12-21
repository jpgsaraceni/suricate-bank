package accountuc

import (
	"errors"
	"reflect"
	"testing"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestDebit(t *testing.T) {
	t.Parallel()

	type args struct {
		account account.Account
		amount  money.Money
	}

	type testCase struct {
		name string
		args args
		want account.Account
		err  error
	}

	var (
		testMoney100, _ = money.NewMoney(100)
		testMoney10, _  = money.NewMoney(10)
		testMoney110, _ = money.NewMoney(110)
		testMoney0, _   = money.NewMoney(0)
	)

	testCases := []testCase{
		{
			name: "successfully debit 100 from account with 110 initial balance",
			args: args{
				account: account.Account{
					Balance: testMoney110,
				},
				amount: testMoney100,
			},
			want: account.Account{Balance: testMoney10},
		},
		{
			name: "successfully debit 100 from account with 100 initial balance",
			args: args{
				account: account.Account{
					Balance: testMoney100,
				},
				amount: testMoney100,
			},
			want: account.Account{Balance: testMoney0},
		},
		{
			name: "fail to debit 0 from account with 10 initial balance",
			args: args{
				account: account.Account{
					Balance: testMoney10,
				},
				amount: testMoney0,
			},
			want: account.Account{Balance: testMoney10},
			err:  money.ErrChangeByZero,
		},
		{
			name: "fail to debit 0 from account with 0 initial balance",
			args: args{
				account: account.Account{},
				amount:  testMoney0,
			},
			want: account.Account{Balance: testMoney0},
			err:  money.ErrChangeByZero,
		},
		{
			name: "fail to debit 100 from account with 0 initial balance",
			args: args{
				account: account.Account{},
				amount:  testMoney100,
			},
			want: account.Account{Balance: testMoney0},
			err:  money.ErrInsuficientFunds,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := Debit(&tt.args.account, tt.args.amount)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			if !reflect.DeepEqual(tt.args.account.Balance, tt.want.Balance) {
				t.Errorf("got %v expected %v", tt.args.account.Balance, tt.want.Balance)
			}
		})
	}
}
