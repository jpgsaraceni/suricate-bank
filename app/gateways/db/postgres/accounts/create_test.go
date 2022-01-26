package accountspg

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type testCase struct {
		name      string
		runBefore func() error
		account   *account.Account
		err       error
	}

	var (
		repeatedId  = account.AccountId(uuid.New())
		repeatedCpf = cpf.Random()
	)

	testCases := []testCase{
		{
			name: "successfully create account",
			account: &account.Account{
				Id:     account.AccountId(uuid.New()),
				Name:   "Nice name",
				Cpf:    cpf.Random(),
				Secret: testHash,
			},
		},
		{
			name: "successfully create account with initial balance",
			account: &account.Account{
				Id:      account.AccountId(uuid.New()),
				Name:    "Nice name",
				Cpf:     cpf.Random(),
				Secret:  testHash,
				Balance: testMoney10,
			},
		},
		{
			name: "fail to create 2 accounts with same id",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					repeatedId,
					cpf.Random().Value(),
					0,
				)
			},
			account: &account.Account{
				Id:     repeatedId,
				Name:   "Nice name",
				Cpf:    cpf.Random(),
				Secret: testHash,
			},
			err: ErrQuery,
		},
		{
			name: "fail to create 2 accounts with same cpf",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					account.AccountId(uuid.New()),
					repeatedCpf.Value(),
					0,
				)
			},
			account: &account.Account{
				Id:     account.AccountId(uuid.New()),
				Name:   "Nice name",
				Cpf:    repeatedCpf,
				Secret: testHash,
			},
			err: accountuc.ErrDuplicateCpf,
		},
		{
			name: "successfully create 2 different accounts in sequence",
			runBefore: func() error {
				return testRepo.Create(
					testContext,
					&account.Account{
						Id:     account.AccountId(uuid.New()),
						Name:   "Nice name",
						Cpf:    cpf.Random(),
						Secret: testHash,
					},
				)
			},
			account: &account.Account{
				Id:     account.AccountId(uuid.New()),
				Name:   "Another nice name",
				Cpf:    cpf.Random(),
				Secret: testHash,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				err := tt.runBefore()

				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			if err := testRepo.Create(testContext, tt.account); !errors.Is(err, tt.err) {

				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}
		})
	}
}
