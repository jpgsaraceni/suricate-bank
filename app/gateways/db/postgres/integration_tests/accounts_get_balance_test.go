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

func TestGetBalance(t *testing.T) {
	t.Parallel()

	var (
		accountId1   = account.AccountId(uuid.New())
		accountId2   = account.AccountId(uuid.New())
		testHash, _  = hash.NewHash("nicesecret")
		testMoney, _ = money.NewMoney(10)
	)

	type testCase struct {
		name            string
		runBefore       func(repo *accountspg.Repository) error
		accountId       account.AccountId
		expectedBalance int
		err             error
	}

	testCases := []testCase{
		{
			name: "successfully get 0 balance",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:     accountId1,
						Name:   "Nice name",
						Cpf:    cpf.Random(),
						Secret: testHash,
					},
				)
			},
			accountId:       accountId1,
			expectedBalance: 0,
		},
		{
			name: "successfully get 10 balance",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:      accountId2,
						Name:    "Nice name",
						Cpf:     cpf.Random(),
						Secret:  testHash,
						Balance: testMoney,
					},
				)
			},
			accountId:       accountId2,
			expectedBalance: 10,
		},
		{
			name: "fail to get balance from inexistent account",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:      account.AccountId(uuid.New()),
						Name:    "Nice name",
						Cpf:     cpf.Random(),
						Secret:  testHash,
						Balance: testMoney,
					},
				)
			},
			accountId:       account.AccountId(uuid.New()),
			expectedBalance: 0,
			err:             accountspg.ErrQuery,
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

			gotBalance, err := repo.GetBalance(testContext, tt.accountId)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if gotBalance != tt.expectedBalance {
				t.Fatalf("got %d expected %d", tt.expectedBalance, gotBalance)
			}
		})
	}
}
