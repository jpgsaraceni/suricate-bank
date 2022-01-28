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

func TestCredit(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type args struct {
		accountId account.AccountId
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
		testIdCredit10initial0  = account.AccountId(uuid.New())
		testIdCredit10initial10 = account.AccountId(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully credit 10 to account with 0 balance",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIdCredit10initial0,
					cpf.Random().Value(),
					0,
				)
			},
			args: args{
				accountId: testIdCredit10initial0,
				amount:    testMoney10,
			},
			expectedBalance: 10,
		},
		{
			name: "successfully credit 10 to account with 10 balance",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIdCredit10initial10,
					cpf.Random().Value(),
					10,
				)
			},
			args: args{
				accountId: testIdCredit10initial10,
				amount:    testMoney10,
			},
			expectedBalance: 20,
		},
		{
			name: "fail to credit inexistent account",
			args: args{
				accountId: account.AccountId(uuid.New()),
				amount:    testMoney10,
			},
			err: account.ErrIdNotFound,
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

			if err := testRepo.CreditAccount(testContext, tt.args.accountId, tt.args.amount); !errors.Is(err, tt.err) {
				t.Fatalf(" got error: %s expected error: %s", err, tt.err)
			}

			if gotBalance, _ := testRepo.GetBalance(testContext, tt.args.accountId); gotBalance != tt.expectedBalance {
				t.Fatalf("got %d expected %d ", gotBalance, tt.expectedBalance)
			}
		})
	}
}
