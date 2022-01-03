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

func TestCredit(t *testing.T) {

	var (
		testId       = account.AccountId(uuid.New())
		testCpf      = cpf.Random()
		testHash, _  = hash.NewHash("nicesecret")
		testMoney, _ = money.NewMoney(10)
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
			name: "successfully credit 10 to account with 0 balance",
			runBefore: func(repo *accountspg.Repository) error {
				truncateAccounts()
				return repo.Create(
					testContext,
					&account.Account{
						Id:     testId,
						Name:   "Nice name",
						Cpf:    testCpf,
						Secret: testHash,
					},
				)
			},
			args: args{
				accountId: testId,
				amount:    testMoney,
			},
			expectedBalance: 10,
		},
		{
			name: "successfully credit 10 to account with 10 balance",
			runBefore: func(repo *accountspg.Repository) error {
				truncateAccounts()
				return repo.Create(
					testContext,
					&account.Account{
						Id:      testId,
						Name:    "Nice name",
						Cpf:     testCpf,
						Secret:  testHash,
						Balance: testMoney,
					},
				)
			},
			args: args{
				accountId: testId,
				amount:    testMoney,
			},
			expectedBalance: 20,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			repo := accountspg.NewRepository(dbPool)

			if tt.runBefore != nil {
				err := tt.runBefore(repo)

				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			if err := repo.CreditAccount(testContext, tt.args.accountId, tt.args.amount); !errors.Is(err, tt.err) {
				t.Fatalf(" got error: %s expected error: %s", err, tt.err)
			}

			if gotBalance, _ := repo.GetBalance(testContext, tt.args.accountId); gotBalance != tt.expectedBalance {
				t.Fatalf("got %d expected %d ", gotBalance, tt.expectedBalance)
			}
		})
	}
}
