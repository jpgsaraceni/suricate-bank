package transferuc

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name             string
		repository       transfer.Repository
		debiter          Debiter
		crediter         Crediter
		transferInstance transfer.Transfer
		err              error
	}

	var testUUID1, _ = uuid.NewUUID()
	var testUUID2, _ = uuid.NewUUID()
	var testUUID3, _ = uuid.NewUUID()

	var testTransferId = transfer.TransferId(testUUID3)

	var testMoney100, _ = money.NewMoney(100)

	var testTime = time.Now()

	testCases := []testCase{
		{
			name: "create transfer",
			repository: transfer.MockRepository{
				OnCreate: func(ctx context.Context, transfer *transfer.Transfer) error {
					transfer.Id = testTransferId
					transfer.CreatedAt = testTime
					return nil
				},
			},
			debiter: MockDebiter{
				OnDebit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			crediter: MockCrediter{
				OnCredit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			transferInstance: transfer.Transfer{
				Id:                   transfer.TransferId(testUUID3),
				Amount:               testMoney100,
				AccountOriginId:      account.AccountId(testUUID1),
				AccountDestinationId: account.AccountId(testUUID2),
				CreatedAt:            testTime,
			},
		},
		// {
		// 	name: "fail transfer to same account",
		// 	args: args{
		// 		amount:        testMoney100,
		// 		originId:      account.AccountId(testUUID1),
		// 		destinationId: account.AccountId(testUUID1),
		// 	},
		// 	want: transfer.Transfer{},
		// 	err:  transfer.ErrSameAccounts,
		// },
		// {
		// 	name: "fail to debit from origin",
		// 	debiter: MockDebiter{
		// 		OnDebit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
		// 			return accountuc.ErrRepository
		// 		},
		// 	},
		// 	args: args{
		// 		amount:        testMoney100,
		// 		originId:      account.AccountId(testUUID1),
		// 		destinationId: account.AccountId(testUUID2),
		// 	},
		// 	want: transfer.Transfer{},
		// 	err:  accountuc.ErrRepository,
		// },
		// {
		// 	name: "fail to credit to destination",
		// 	debiter: MockDebiter{
		// 		OnDebit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
		// 			return nil
		// 		},
		// 	},
		// 	crediter: MockCrediter{
		// 		OnCredit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
		// 			return accountuc.ErrRepository
		// 		},
		// 	},
		// 	args: args{
		// 		amount:        testMoney100,
		// 		originId:      account.AccountId(testUUID1),
		// 		destinationId: account.AccountId(testUUID2),
		// 	},
		// 	want: transfer.Transfer{},
		// 	err:  accountuc.ErrRepository,
		// },
		// {
		// 	name: "fail to create transfer amount 0",
		// 	args: args{
		// 		amount:        testMoney0,
		// 		originId:      account.AccountId(testUUID1),
		// 		destinationId: account.AccountId(testUUID2),
		// 	},
		// 	want: transfer.Transfer{},
		// 	err:  transfer.ErrAmountNotPositive,
		// },
		// {
		// 	name: "repository error creating transfer",
		// 	repository: transfer.MockRepository{
		// 		OnCreate: func(ctx context.Context, transfer *transfer.Transfer) error {
		// 			transfer.Id = testTransferId
		// 			transfer.CreatedAt = testTime
		// 			return ErrRepository
		// 		},
		// 	},
		// 	debiter: MockDebiter{
		// 		OnDebit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
		// 			return nil
		// 		},
		// 	},
		// 	crediter: MockCrediter{
		// 		OnCredit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
		// 			return nil
		// 		},
		// 	},
		// 	args: args{
		// 		amount:        testMoney100,
		// 		originId:      account.AccountId(testUUID1),
		// 		destinationId: account.AccountId(testUUID2),
		// 	},
		// 	want: transfer.Transfer{},
		// 	err:  ErrRepository,
		// },
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository, tt.crediter, tt.debiter}

			err := uc.Create(context.Background(), tt.transferInstance)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected %v", err, tt.err)
			}
		})
	}
}
