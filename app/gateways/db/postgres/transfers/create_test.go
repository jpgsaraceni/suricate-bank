package transferspg

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v4/pgxpool"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/postgrestest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name      string
		runBefore func(testPool *pgxpool.Pool) error
		transfer  *transfer.Transfer
		err       error
	}

	var (
		testAccountIdInitial0  = account.AccountId(uuid.New())
		testAccountIdInitial10 = account.AccountId(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully create transfer",
			runBefore: func(testPool *pgxpool.Pool) error {
				return accountspg.CreateTestAccountBatch(
					testPool,
					[]account.AccountId{
						testAccountIdInitial0,
						testAccountIdInitial10,
					},
					[]string{
						cpf.Random().Value(),
						cpf.Random().Value(),
					},
					[]int{
						0,
						10,
					},
				)
			},
			transfer: &transfer.Transfer{
				Id:                   transfer.TransferId(uuid.New()),
				AccountOriginId:      testAccountIdInitial10,
				AccountDestinationId: testAccountIdInitial0,
				Amount:               testMoney10,
				CreatedAt:            testTime,
			},
		},
		{
			name: "fail transfer from inexistent account",
			runBefore: func(testPool *pgxpool.Pool) error {
				return accountspg.CreateTestAccountBatch(
					testPool,
					[]account.AccountId{
						testAccountIdInitial0,
						testAccountIdInitial10,
					},
					[]string{
						cpf.Random().Value(),
						cpf.Random().Value(),
					},
					[]int{
						0,
						10,
					},
				)
			},
			transfer: &transfer.Transfer{
				Id:                   transfer.TransferId(uuid.New()),
				AccountOriginId:      account.AccountId(uuid.New()),
				AccountDestinationId: testAccountIdInitial0,
				Amount:               testMoney10,
				CreatedAt:            testTime,
			},
			err: ErrCreateTransfer,
		},
		{
			name: "fail transfer to inexistent account",
			runBefore: func(testPool *pgxpool.Pool) error {
				return accountspg.CreateTestAccountBatch(
					testPool,
					[]account.AccountId{
						testAccountIdInitial0,
						testAccountIdInitial10,
					},
					[]string{
						cpf.Random().Value(),
						cpf.Random().Value(),
					},
					[]int{
						0,
						10,
					},
				)
			},
			transfer: &transfer.Transfer{
				Id:                   transfer.TransferId(uuid.New()),
				AccountOriginId:      testAccountIdInitial10,
				AccountDestinationId: account.AccountId(uuid.New()),
				Amount:               testMoney10,
				CreatedAt:            testTime,
			},
			err: ErrCreateTransfer,
		},
		{
			name: "fail transfer to same account",
			runBefore: func(testPool *pgxpool.Pool) error {
				return accountspg.CreateTestAccountBatch(
					testPool,
					[]account.AccountId{
						testAccountIdInitial0,
						testAccountIdInitial10,
					},
					[]string{
						cpf.Random().Value(),
						cpf.Random().Value(),
					},
					[]int{
						0,
						10,
					},
				)
			},
			transfer: &transfer.Transfer{
				Id:                   transfer.TransferId(uuid.New()),
				AccountOriginId:      testAccountIdInitial10,
				AccountDestinationId: testAccountIdInitial10,
				Amount:               testMoney10,
				CreatedAt:            testTime,
			},
			err: ErrCreateTransfer,
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

			if err := testRepo.Create(testContext, tt.transfer); !errors.Is(err, tt.err) {

				t.Fatalf("got error: %s expected error: %s", err, tt.err)
			}
		})
	}
}
