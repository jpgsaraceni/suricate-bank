package accountspg

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
)

func TestFetch(t *testing.T) {
	// can't run this test in parallel because Fetch would possibly return accounts
	// created in parallel tests after the truncate function call

	var (
		accountId1  = account.AccountId(uuid.New())
		accountId2  = account.AccountId(uuid.New())
		cpf1        = cpf.Random()
		cpf2        = cpf.Random()
		testTime    = time.Now().Round(time.Hour)
		testHash, _ = hash.NewHash("nicesecret")
	)

	type testCase struct {
		name             string
		expectedAccounts []account.Account
		runBefore        func(repo *Repository) error
		err              error
	}

	testCases := []testCase{
		{
			name: "successfully fetch 2 accounts",
			runBefore: func(repo *Repository) error {
				postgrestest.Truncate(dbPool, "accounts")
				err := repo.Create(
					testContext,
					&account.Account{
						Id:        accountId1,
						Name:      "Nice name",
						Cpf:       cpf1,
						Secret:    testHash,
						CreatedAt: testTime,
					},
				)

				if err != nil {

					return err
				}

				err = repo.Create(
					testContext,
					&account.Account{
						Id:        accountId2,
						Name:      "Another nice name",
						Cpf:       cpf2,
						Secret:    testHash,
						CreatedAt: testTime,
					},
				)

				return err
			},
			expectedAccounts: []account.Account{
				{
					Id:        accountId1,
					Name:      "Nice name",
					Cpf:       cpf1,
					Secret:    testHash,
					CreatedAt: testTime,
				},
				{
					Id:        accountId2,
					Name:      "Another nice name",
					Cpf:       cpf2,
					Secret:    testHash,
					CreatedAt: testTime,
				},
			},
		},
		{
			name: "successfully fetch 1 account",
			runBefore: func(repo *Repository) error {
				postgrestest.Truncate(dbPool, "accounts")
				return repo.Create(
					testContext,
					&account.Account{
						Id:        accountId1,
						Name:      "Nice name",
						Cpf:       cpf1,
						Secret:    testHash,
						CreatedAt: testTime,
					},
				)
			},
			expectedAccounts: []account.Account{
				{
					Id:        accountId1,
					Name:      "Nice name",
					Cpf:       cpf1,
					Secret:    testHash,
					CreatedAt: testTime,
				},
			},
		},
		{
			name: "fail to fetch 0 accounts",
			runBefore: func(repo *Repository) error {
				postgrestest.Truncate(dbPool, "accounts")
				return nil
			},
			expectedAccounts: nil,
			err:              ErrEmptyFetch,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			repo := NewRepository(dbPool)

			if tt.runBefore != nil {
				err := tt.runBefore(repo)

				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			gotAccounts, err := repo.Fetch(testContext)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if !reflect.DeepEqual(gotAccounts, tt.expectedAccounts) {
				t.Fatalf("got %v expected %v", gotAccounts, tt.expectedAccounts)
			}
		})
	}
}
