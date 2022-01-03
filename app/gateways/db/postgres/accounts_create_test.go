package postgres_test

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	var (
		repeatedId  = account.AccountId(uuid.New())
		testHash, _ = hash.NewHash("nicesecret")
	)

	type args struct {
		ctx     context.Context
		account *account.Account
	}

	type testCase struct {
		name      string
		runBefore func(repo *accountspg.Repository) error
		args      args
		err       error
	}

	testCases := []testCase{
		{
			name: "successfully create account",
			args: args{
				ctx: testContext,
				account: &account.Account{
					Id:     account.AccountId(uuid.New()),
					Name:   "Nice name",
					Cpf:    cpf.Random(),
					Secret: testHash,
				},
			},
		},
		{
			name: "fail to create account with repeated id",
			runBefore: func(repo *accountspg.Repository) error {
				truncateAccounts()
				return repo.Create(
					testContext,
					&account.Account{
						Id:     repeatedId,
						Name:   "Another nice name",
						Cpf:    cpf.Random(),
						Secret: testHash,
					},
				)
			},
			args: args{
				ctx: testContext,
				account: &account.Account{
					Id:     repeatedId,
					Name:   "Nice name",
					Cpf:    cpf.Random(),
					Secret: testHash,
				},
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

			if err := repo.Create(testContext, tt.args.account); !errors.Is(err, tt.err) {

				t.Fatalf("expected error: %s got error: %s", tt.err, err)
			}
		})
	}
}
