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

func TestCreate(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type testCase struct {
		name      string
		runBefore func() error
		account   account.Account
		err       error
	}

	var (
		testSecret, _ = hash.NewHash("123456")
		repeatedID    = account.ID(uuid.New())
		repeatedCpf   = cpf.Random()
	)

	testCases := []testCase{
		{
			name: "successfully create account",
			account: account.Account{
				ID:        account.ID(uuid.New()),
				Name:      "Nice name",
				Cpf:       cpf.Random(),
				Secret:    testSecret,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "successfully create account with initial balance",
			account: account.Account{
				ID:        account.ID(uuid.New()),
				Name:      "Nice name",
				Cpf:       cpf.Random(),
				Secret:    testSecret,
				Balance:   testMoney10,
				CreatedAt: time.Now(),
			},
		},
		{
			name: "fail to create 2 accounts with same id",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					repeatedID,
					cpf.Random().Value(),
					0,
				)
			},
			account: account.Account{
				ID:        repeatedID,
				Name:      "Nice name",
				Cpf:       cpf.Random(),
				Secret:    testSecret,
				CreatedAt: time.Now(),
			},
			err: ErrQuery,
		},
		{
			name: "fail to create 2 accounts with same cpf",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					account.ID(uuid.New()),
					repeatedCpf.Value(),
					0,
				)
			},
			account: account.Account{
				ID:        account.ID(uuid.New()),
				Name:      "Nice name",
				Cpf:       repeatedCpf,
				Secret:    testSecret,
				CreatedAt: time.Now(),
			},
			err: account.ErrDuplicateCpf,
		},
		{
			name: "successfully create 2 different accounts in sequence",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					account.ID(uuid.New()),
					cpf.Random().Value(),
					0,
				)
			},
			account: account.Account{
				ID:        account.ID(uuid.New()),
				Name:      "Nice name",
				Cpf:       cpf.Random(),
				Secret:    testSecret,
				CreatedAt: time.Now(),
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

			gotAccount, err := testRepo.Create(testContext, tt.account)
			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if err != nil && !reflect.DeepEqual(gotAccount, account.Account{}) {
				t.Fatalf("got %v expected empty account", gotAccount)
			}
		})
	}
}
