package postgres

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
)

func TestGetById(t *testing.T) {
	t.Parallel()

	var (
		accountId   = account.AccountId(uuid.New())
		accountCpf  = cpf.Random()
		testHash, _ = hash.NewHash("nicesecret")
	)

	type testCase struct {
		name            string
		runBefore       func(repo *accountspg.Repository) error
		expectedAccount account.Account
		err             error
	}

	testCases := []testCase{
		{
			name: "successfully get an account",
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
			expectedAccount: account.Account{
				Id:     accountId,
				Name:   "Nice name",
				Cpf:    accountCpf,
				Secret: testHash,
			},
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

			gotAccount, err := repo.GetById(testContext, accountId)

			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error: %s got error: %s", tt.err, err)
			}

			tt.expectedAccount.CreatedAt = gotAccount.CreatedAt

			if !reflect.DeepEqual(tt.expectedAccount, gotAccount) {
				t.Fatalf("expected %v got %v", tt.expectedAccount, gotAccount)
			}
		})
	}
}
