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

func TestFetch(t *testing.T) {
	t.Parallel()

	testPool, tearDown := postgrestest.GetTestPool()
	testRepo := NewRepository(testPool)

	t.Cleanup(tearDown)

	type testCase struct {
		name             string
		runBefore        func() error
		expectedAccounts []account.Account
		err              error
	}

	var (
		testIdInitial0   = account.AccountId(uuid.New())
		testIdInitial10  = account.AccountId(uuid.New())
		testCpfInitial0  = cpf.Random()
		testCpfInitial10 = cpf.Random()
	)

	testCases := []testCase{
		{
			name: "successfully fetch 2 accounts",
			runBefore: func() error {
				return createTestAccountBatch(
					testPool,
					[]account.AccountId{
						testIdInitial0,
						testIdInitial10,
					},
					[]string{
						testCpfInitial0.Value(),
						testCpfInitial10.Value(),
					},
					[]int{
						0,
						10,
					},
				)
			},
			expectedAccounts: []account.Account{
				{
					Id:        testIdInitial0,
					Name:      "nice name",
					Cpf:       testCpfInitial0,
					Secret:    testHash,
					CreatedAt: testTime,
				},
				{
					Id:        testIdInitial10,
					Name:      "nice name",
					Cpf:       testCpfInitial10,
					Secret:    testHash,
					Balance:   testMoney10,
					CreatedAt: testTime,
				},
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {

			if tt.runBefore != nil {
				err := tt.runBefore()

				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			gotAccounts, err := testRepo.Fetch(testContext)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if !reflect.DeepEqual(gotAccounts, tt.expectedAccounts) {
				t.Fatalf("got\n %v \nexpected\n %v", gotAccounts, tt.expectedAccounts)
			}
		})
	}
}
