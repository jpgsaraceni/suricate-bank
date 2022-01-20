package accountspg

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestGetByCpf(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type testCase struct {
		name            string
		runBefore       func() error
		cpf             cpf.Cpf
		expectedAccount account.Account
		err             error
	}

	var (
		testId  = account.AccountId(uuid.New())
		testCpf = cpf.Random()
	)

	testCases := []testCase{
		{
			name: "successfully get an account",
			runBefore: func() error {
				return createTestAccount(
					testPool,
					testId,
					testCpf.Value(),
					0,
				)
			},
			cpf: testCpf,
			expectedAccount: account.Account{
				Id:        testId,
				Name:      "nice name",
				Cpf:       testCpf,
				Secret:    testHash,
				CreatedAt: testTime,
			},
		},
		{
			name:            "fail to get an inexixtent account",
			cpf:             cpf.Random(),
			expectedAccount: account.Account{},
			err:             ErrCpfNotFound,
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

			gotAccount, err := testRepo.GetByCpf(testContext, tt.cpf)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if !reflect.DeepEqual(gotAccount, tt.expectedAccount) {
				t.Fatalf("got %v expected %v", gotAccount, tt.expectedAccount)
			}
		})
	}
}
