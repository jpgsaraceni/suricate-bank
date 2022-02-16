package accountspg

import (
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestGetBalance(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type testCase struct {
		name            string
		runBefore       func() error
		accountID       account.ID
		expectedBalance int
		err             error
	}

	var (
		testIDInitial0  = account.ID(uuid.New())
		testIDInitial10 = account.ID(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully get 0 balance",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIDInitial0,
					cpf.Random().Value(),
					0,
				)
			},
			accountID:       testIDInitial0,
			expectedBalance: 0,
		},
		{
			name: "successfully get 10 balance",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testIDInitial10,
					cpf.Random().Value(),
					10,
				)
			},
			accountID:       testIDInitial10,
			expectedBalance: 10,
		},
		{
			name:            "fail to get balance from inexistent account",
			accountID:       account.ID(uuid.New()),
			expectedBalance: 0,
			err:             account.ErrIDNotFound,
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

			gotBalance, err := testRepo.GetBalance(testContext, tt.accountID)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if gotBalance != tt.expectedBalance {
				t.Fatalf("got %d expected %d", gotBalance, tt.expectedBalance)
			}
		})
	}
}
