package transferspg

import (
	"errors"
	"reflect"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestFetch(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name              string
		runBefore         func(*pgxpool.Pool) error
		expectedTransfers []transfer.Transfer
		err               error
	}

	var (
		testIDInitial0  = account.ID(uuid.New())
		testIDInitial10 = account.ID(uuid.New())
		testIDInitial20 = account.ID(uuid.New())
		testIDInitial30 = account.ID(uuid.New())
		testTransferID1 = transfer.ID(uuid.New())
		testTransferID2 = transfer.ID(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully fetch 1 transfer",
			runBefore: func(testPool *pgxpool.Pool) error {
				err := accountspg.CreateTestAccountBatch(
					testPool,
					[]account.ID{
						testIDInitial10,
						testIDInitial20,
					},
					[]string{
						cpf.Random().Value(),
						cpf.Random().Value(),
					},
					[]int{
						10,
						20,
					},
				)
				if err != nil {
					return err
				}

				return createTestTransferBatch(
					testPool,
					[]transfer.ID{
						testTransferID1,
					},
					[]account.ID{
						testIDInitial20,
					},
					[]account.ID{
						testIDInitial10,
					},
					[]int{
						10,
					},
				)
			},
			expectedTransfers: []transfer.Transfer{
				{
					ID:                   testTransferID1,
					AccountOriginID:      testIDInitial20,
					AccountDestinationID: testIDInitial10,
					Amount:               testMoney10,
					CreatedAt:            testTime,
				},
			},
		},
		{
			name: "successfully fetch 2 transfers",
			runBefore: func(testPool *pgxpool.Pool) error {
				err := accountspg.CreateTestAccountBatch(
					testPool,
					[]account.ID{
						testIDInitial0,
						testIDInitial10,
						testIDInitial20,
						testIDInitial30,
					},
					[]string{
						cpf.Random().Value(),
						cpf.Random().Value(),
						cpf.Random().Value(),
						cpf.Random().Value(),
					},
					[]int{
						0,
						10,
						20,
						30,
					},
				)
				if err != nil {
					return err
				}

				return createTestTransferBatch(
					testPool,
					[]transfer.ID{
						testTransferID1,
						testTransferID2,
					},
					[]account.ID{
						testIDInitial20,
						testIDInitial30,
					},
					[]account.ID{
						testIDInitial10,
						testIDInitial0,
					},
					[]int{
						10,
						15,
					},
				)
			},
			expectedTransfers: []transfer.Transfer{
				{
					ID:                   testTransferID1,
					AccountOriginID:      testIDInitial20,
					AccountDestinationID: testIDInitial10,
					Amount:               testMoney10,
					CreatedAt:            testTime,
				},
				{
					ID:                   testTransferID2,
					AccountOriginID:      testIDInitial30,
					AccountDestinationID: testIDInitial0,
					Amount:               testMoney15,
					CreatedAt:            testTime,
				},
			},
		},
		{
			name:              "fetch 0 transfers",
			expectedTransfers: nil,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			testPool, tearDown := postgrestest.GetTestPool()
			testRepo := NewRepository(testPool)

			t.Cleanup(tearDown)

			if tt.runBefore != nil {
				err := tt.runBefore(testPool)
				if err != nil {
					t.Fatalf("runBefore() failed: %s", err)
				}
			}

			gotAccounts, err := testRepo.Fetch(testContext)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}

			if !reflect.DeepEqual(gotAccounts, tt.expectedTransfers) {
				t.Fatalf("got\n %v \nexpected\n %v", gotAccounts, tt.expectedTransfers)
			}
		})
	}
}
