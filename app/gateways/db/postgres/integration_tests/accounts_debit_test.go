package postgres_test

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestDebit(t *testing.T) {
	t.Parallel()

	var (
		accountId1     = account.AccountId(uuid.New())
		accountId2     = account.AccountId(uuid.New())
		testHash, _    = hash.NewHash("nicesecret")
		testMoney20, _ = money.NewMoney(20)
		testMoney10, _ = money.NewMoney(10)
	)

	type args struct {
		accountId account.AccountId
		amount    money.Money
	}

	type testCase struct {
		name            string
		runBefore       func(repo *accountspg.Repository) error
		args            args
		expectedBalance int
		err             error
	}

	testCases := []testCase{
		{
			name: "successfully debit 10 from account with 20 balance",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:      accountId1,
						Name:    "Nice name",
						Cpf:     cpf.Random(),
						Secret:  testHash,
						Balance: testMoney20,
					},
				)
			},
			args: args{
				accountId: accountId1,
				amount:    testMoney10,
			},
			expectedBalance: 10,
		},
		{
			name: "successfully debit 20 from account with 20 balance",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:      accountId2,
						Name:    "Nice name",
						Cpf:     cpf.Random(),
						Secret:  testHash,
						Balance: testMoney20,
					},
				)
			},
			args: args{
				accountId: accountId2,
				amount:    testMoney20,
			},
			expectedBalance: 0,
		},
		{
			name: "fail to debit inexistent account",
			args: args{
				accountId: account.AccountId(uuid.New()),
				amount:    testMoney10,
			},
			err: accountspg.ErrQuery,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			repo := accountspg.NewRepository(dbPool)

			if tt.runBefore != nil {
				err := tt.runBefore(repo)

				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			if err := repo.DebitAccount(testContext, tt.args.accountId, tt.args.amount); !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if gotBalance, _ := repo.GetBalance(testContext, tt.args.accountId); gotBalance != tt.expectedBalance {
				t.Fatalf("got %d expected %d", gotBalance, tt.expectedBalance)
			}
		})
	}
}
