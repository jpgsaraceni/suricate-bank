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
		testIdInitial0  = account.AccountId(uuid.New())
		testIdInitial10 = account.AccountId(uuid.New())
		testIdInitial20 = account.AccountId(uuid.New())
		testIdInitial30 = account.AccountId(uuid.New())
		testTransferId1 = transfer.TransferId(uuid.New())
		testTransferId2 = transfer.TransferId(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully fetch 1 transfer",
			runBefore: func(testPool *pgxpool.Pool) error {
				err := accountspg.CreateTestAccountBatch(
					testPool,
					[]account.AccountId{
						testIdInitial10,
						testIdInitial20,
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
					[]transfer.TransferId{
						testTransferId1,
					},
					[]account.AccountId{
						testIdInitial20,
					},
					[]account.AccountId{
						testIdInitial10,
					},
					[]int{
						10,
					},
				)
			},
			expectedTransfers: []transfer.Transfer{
				{
					Id:                   testTransferId1,
					AccountOriginId:      testIdInitial20,
					AccountDestinationId: testIdInitial10,
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
					[]account.AccountId{
						testIdInitial0,
						testIdInitial10,
						testIdInitial20,
						testIdInitial30,
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
					[]transfer.TransferId{
						testTransferId1,
						testTransferId2,
					},
					[]account.AccountId{
						testIdInitial20,
						testIdInitial30,
					},
					[]account.AccountId{
						testIdInitial10,
						testIdInitial0,
					},
					[]int{
						10,
						15,
					},
				)
			},
			expectedTransfers: []transfer.Transfer{
				{
					Id:                   testTransferId1,
					AccountOriginId:      testIdInitial20,
					AccountDestinationId: testIdInitial10,
					Amount:               testMoney10,
					CreatedAt:            testTime,
				},
				{
					Id:                   testTransferId2,
					AccountOriginId:      testIdInitial30,
					AccountDestinationId: testIdInitial0,
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
