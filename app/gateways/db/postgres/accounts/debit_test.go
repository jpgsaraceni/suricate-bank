package accountspg

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestDebit(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type args struct {
		accountID account.ID
		amount    money.Money
	}

	type testCase struct {
		name            string
		runBefore       func() error
		args            args
		expectedBalance int
		err             error
	}

	var (
		testIDDebit10initial20 = account.ID(uuid.New())
		testIDDebit30initial30 = account.ID(uuid.New())
		testIDDebit20initial10 = account.ID(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully debit 10 from account with 20 balance",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIDDebit10initial20,
					cpf.Random().Value(),
					20,
				)
			},
			args: args{
				accountID: testIDDebit10initial20,
				amount:    testMoney10,
			},
			expectedBalance: 10,
		},
		{
			name: "successfully debit 30 from account with 30 balance",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIDDebit30initial30,
					cpf.Random().Value(),
					30,
				)
			},
			args: args{
				accountID: testIDDebit30initial30,
				amount:    testMoney30,
			},
			expectedBalance: 0,
		},
		{
			name: "fail to debit 20 from account with 10 balance",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIDDebit20initial10,
					cpf.Random().Value(),
					10,
				)
			},
			args: args{
				accountID: testIDDebit20initial10,
				amount:    testMoney20,
			},
			expectedBalance: 10,
			err:             account.ErrInsufficientFunds,
		},
		{
			name: "fail to debit inexistent account",
			args: args{
				accountID: account.ID(uuid.New()),
				amount:    testMoney10,
			},
			err: account.ErrIDNotFound,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				err := tt.runBefore()
				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			if err := testRepo.DebitAccount(testContext, tt.args.accountID, tt.args.amount); !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if gotBalance, _ := testRepo.GetBalance(testContext, tt.args.accountID); gotBalance != tt.expectedBalance {
				t.Fatalf("got %d expected %d", gotBalance, tt.expectedBalance)
			}
		})
	}
}
