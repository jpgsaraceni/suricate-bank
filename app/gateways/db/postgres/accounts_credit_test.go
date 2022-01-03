package postgres

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
	t.Parallel()

	var (
		testId       = account.AccountId(uuid.New())
		testCpf      = cpf.Random()
		testHash, _  = hash.NewHash("nicesecret")
		testMoney, _ = money.NewMoney(10)
	)

	type testCase struct {
		name            string
		runBefore       func(repo *accountspg.Repository) error
		amount          money.Money
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
			amount:          testMoney,
			expectedBalance: 10,
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

			if err := repo.CreditAccount(testContext, testId, tt.amount); !errors.Is(err, tt.err) {
				t.Fatalf("expected error: %s got error: %s", tt.err, err)
			}

			if gotBalance, _ := repo.GetBalance(testContext, testId); gotBalance != tt.expectedBalance {
				t.Fatalf("expected %d got %d", tt.expectedBalance, gotBalance)
			}
		})
	}
}
