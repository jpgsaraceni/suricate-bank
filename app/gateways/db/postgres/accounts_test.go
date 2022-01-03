package postgres

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

var (
	testId         = account.AccountId(uuid.New())
	testCpf, _     = cpf.NewCpf("22061446035")
	testCpf2, _    = cpf.NewCpf("04559118000")
	testHash, _    = hash.NewHash("nicesecret")
	testHash2, _   = hash.NewHash("anothernicesecret")
	testMoney10, _ = money.NewMoney(10)
	testMoney30, _ = money.NewMoney(30)
	testContext    = context.Background()
)

func truncate() error {
	_, err := dbPool.Exec(testContext, "TRUNCATE accounts")

	if err != nil {

		return err
	}

	return nil
}

func TestCreate(t *testing.T) {
	t.Parallel()

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
					Id:     testId,
					Name:   "Nice name",
					Cpf:    testCpf,
					Secret: testHash,
				},
			},
		},
		{
			name: "fail to create account with repeated id",
			runBefore: func(repo *accountspg.Repository) error {
				return repo.Create(
					testContext,
					&account.Account{
						Id:     testId,
						Name:   "Another nice name",
						Cpf:    testCpf2,
						Secret: testHash2,
					},
				)
			},
			args: args{
				ctx: testContext,
				account: &account.Account{
					Id:     testId,
					Name:   "Nice name",
					Cpf:    testCpf,
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

			defer truncate()

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

func TestAccount(t *testing.T) {
	// t.Parallel()
	// type testCase struct {
	// 	name    string
	// 	account account.Account
	// 	err     error
	// }

	// tt := testCase{
	// 	name: "test account",
	// 	account: account.Account{
	// 		Id:        testId,
	// 		Cpf:       testCpf,
	// 		Name:      "Nice name",
	// 		Secret:    testHash,
	// 		CreatedAt: time.Now(),
	// 	},
	// }

	// repo := accountspg.NewRepository(dbPool)
	// if err := repo.Create(testContext, &tt.account); !errors.Is(err, tt.err) {
	// 	t.Fatalf("expected error: %s got error: %s", tt.err, err)
	// }
	// accounts, err := repo.Fetch(testContext)

	// if !errors.Is(err, tt.err) {
	// 	t.Fatalf("expected error: %s got error: %s", tt.err, err)
	// }

	// account, err := repo.GetById(testContext, accounts[0].Id)

	// if !errors.Is(err, tt.err) {
	// 	t.Fatalf("expected error: %s got error: %s", tt.err, err)
	// }

	// if err := repo.CreditAccount(testContext, account.Id, testMoney30); !errors.Is(err, tt.err) {
	// 	t.Fatalf("expected error: %s got error: %s", tt.err, err)
	// }

	// balance, err := repo.GetBalance(testContext, account.Id)

	// if !errors.Is(err, tt.err) {
	// 	t.Fatalf("expected error: %s got error: %s", tt.err, err)
	// }

	// if balance != 30 {
	// 	t.Fatalf("expected balance: 30 got balance: %d", balance)
	// }

	// if err := repo.DebitAccount(testContext, account.Id, testMoney10); !errors.Is(err, tt.err) {
	// 	t.Fatalf("expected error: %s got error: %s", tt.err, err)
	// }

	// balance, err = repo.GetBalance(testContext, account.Id)

	// if !errors.Is(err, tt.err) {
	// 	t.Fatalf("expected error: %s got error: %s", tt.err, err)
	// }

	// if balance != 20 {
	// 	t.Fatalf("expected balance: 20 got balance: %d", balance)
	// }
}
