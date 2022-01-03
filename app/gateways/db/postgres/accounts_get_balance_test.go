package postgres

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	var (
		accountId   = account.AccountId(uuid.New())
		accountCpf  = cpf.Random()
		testHash, _ = hash.NewHash("nicesecret")
	)

	type testCase struct {
		name            string
		runBefore       func(repo *accountspg.Repository) error
		expectedBalance int
		err             error
	}

	testCases := []testCase{
		{
			name: "successfully get 0 balance",
			runBefore: func(repo *accountspg.Repository) error {
				truncateAccounts()
				return repo.Create(
					testContext,
					&account.Account{
						Id:     accountId,
						Name:   "Nice name",
						Cpf:    accountCpf,
						Secret: testHash,
					},
				)
			},
			expectedBalance: 0,
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

			gotBalance, err := repo.GetBalance(testContext, accountId)

			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error: %s got error: %s", tt.err, err)
			}

			if gotBalance != tt.expectedBalance {
				t.Fatalf("expected %d got %d", tt.expectedBalance, gotBalance)
			}
		})
	}
}
