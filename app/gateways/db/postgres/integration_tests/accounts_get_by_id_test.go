package postgres_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

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
		cpf1        = cpf.Random()
		cpf2        = cpf.Random()
		testHash, _ = hash.NewHash("nicesecret")
		testTime    = time.Now().Round(time.Hour)
	)

	type testCase struct {
		name            string
		runBefore       func(repo *accountspg.Repository) error
		idArg           account.AccountId
		expectedAccount account.Account
		err             error
	}

	testCases := []testCase{
		{
			name: "successfully get an account",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:        accountId,
						Name:      "Nice name",
						Cpf:       cpf1,
						Secret:    testHash,
						CreatedAt: testTime,
					},
				)
			},
			idArg: accountId,
			expectedAccount: account.Account{
				Id:        accountId,
				Name:      "Nice name",
				Cpf:       cpf1,
				Secret:    testHash,
				CreatedAt: testTime,
			},
		},
		{
			name: "fail to get an inexixtent account",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:        account.AccountId(uuid.New()),
						Name:      "Nice name",
						Cpf:       cpf2,
						Secret:    testHash,
						CreatedAt: testTime,
					},
				)
			},
			idArg:           account.AccountId(uuid.New()),
			expectedAccount: account.Account{},
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

			gotAccount, err := repo.GetById(testContext, tt.idArg)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if !reflect.DeepEqual(gotAccount, tt.expectedAccount) {
				t.Fatalf("got %v expected %v", gotAccount, tt.expectedAccount)
			}
		})
	}
}
