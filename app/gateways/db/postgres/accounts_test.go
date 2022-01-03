package postgres

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
)

var (
	testHash, _ = hash.NewHash("nicesecret")
	// testMoney10, _ = money.NewMoney(10)
	// testMoney30, _ = money.NewMoney(30)
	testContext = context.Background()
)

func TestCreate(t *testing.T) {
	t.Parallel()

	repeatedId := account.AccountId(uuid.New())

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
				truncate()
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

func TestFetch(t *testing.T) {
	t.Parallel()

	var (
		accountId1 = account.AccountId(uuid.New())
		accountId2 = account.AccountId(uuid.New())
		accountId3 = account.AccountId(uuid.New())
	)

	type testCase struct {
		name      string
		runBefore func(repo *accountspg.Repository) error
		err       error
	}

	testCases := []testCase{
		{
			name: "successfully fetch 2 accounts",
			runBefore: func(repo *accountspg.Repository) error {
				truncate()
				err := repo.Create(
					testContext,
					&account.Account{
						Id:     accountId1,
						Name:   "Nice name",
						Cpf:    cpf.Random(),
						Secret: testHash,
					},
				)

				if err != nil {

					return err
				}

				err = repo.Create(
					testContext,
					&account.Account{
						Id:     accountId2,
						Name:   "Another nice name",
						Cpf:    cpf.Random(),
						Secret: testHash,
					},
				)

				return err
			},
		},
		{
			name: "successfully fetch 1 account",
			runBefore: func(repo *accountspg.Repository) error {
				truncate()
				return repo.Create(
					testContext,
					&account.Account{
						Id:     accountId3,
						Name:   "Nice name",
						Cpf:    cpf.Random(),
						Secret: testHash,
					},
				)
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

			_, err := repo.Fetch(testContext)

			//TODO: compare account to expected

			if !errors.Is(err, tt.err) {
				t.Fatalf("expected error: %s got error: %s", tt.err, err)
			}
		})
	}
}

func TestGetById(t *testing.T) {
	t.Parallel()

	accountId := account.AccountId(uuid.New())
	accountCpf := cpf.Random()

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
				truncate()
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

func truncate() error {
	_, err := dbPool.Exec(testContext, "TRUNCATE accounts")

	if err != nil {

		return err
	}

	return nil
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
	// 		Id:        account.AccountId(uuid.New()),
	// 		Cpf:       cpf.Random(),
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
