package transferuc

import (
	"context"
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type args struct {
		amount        money.Money
		originId      account.AccountId
		destinationId account.AccountId
	}

	type testCase struct {
		name       string
		repository transfer.Repository
		debiter    Debiter
		crediter   Crediter
		args       args
		want       transfer.Transfer
		err        error
	}

	var testUUID1, _ = uuid.NewUUID()
	var testUUID2, _ = uuid.NewUUID()
	var testUUID3, _ = uuid.NewUUID()

	var testTransferId = transfer.TransferId(testUUID3)

	var testMoney100, _ = money.NewMoney(100)
	var testMoney0, _ = money.NewMoney(0)

	var testTime = time.Now()

	var errRepository = errors.New("repository error")

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
			args: args{
				amount:        testMoney100,
				originId:      account.AccountId(testUUID1),
				destinationId: account.AccountId(testUUID2),
			},
			want: transfer.Transfer{
				Id:                   transfer.TransferId(testUUID3),
				Amount:               testMoney100,
				AccountOriginId:      account.AccountId(testUUID1),
				AccountDestinationId: account.AccountId(testUUID2),
				CreatedAt:            testTime,
			},
		},
		{
			name: "fail transfer to same account",
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
			args: args{
				amount:        testMoney100,
				originId:      account.AccountId(testUUID1),
				destinationId: account.AccountId(testUUID1),
			},
			want: transfer.Transfer{},
			err:  ErrSameAccounts,
		},
		{
			name: "fail to debit from origin",
			repository: transfer.MockRepository{
				OnCreate: func(ctx context.Context, transfer *transfer.Transfer) error {
					transfer.Id = testTransferId
					transfer.CreatedAt = testTime
					return nil
				},
			},
			debiter: MockDebiter{
				OnDebit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return errRepository
				},
			},
			crediter: MockCrediter{
				OnCredit: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			args: args{
				amount:        testMoney100,
				originId:      account.AccountId(testUUID1),
				destinationId: account.AccountId(testUUID2),
			},
			want: transfer.Transfer{},
			err:  ErrDebitOrigin,
		},
		{
			name: "fail to credit to destination",
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
					return errRepository
				},
			},
			args: args{
				amount:        testMoney100,
				originId:      account.AccountId(testUUID1),
				destinationId: account.AccountId(testUUID2),
			},
			want: transfer.Transfer{},
			err:  ErrCreditDestination,
		},
		{
			name: "fail to create transfer amount 0",
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
			args: args{
				amount:        testMoney0,
				originId:      account.AccountId(testUUID1),
				destinationId: account.AccountId(testUUID2),
			},
			want: transfer.Transfer{},
			err:  transfer.ErrAmountNotPositive,
		},
		{
			name: "repository error creating transfer",
			repository: transfer.MockRepository{
				OnCreate: func(ctx context.Context, transfer *transfer.Transfer) error {
					transfer.Id = testTransferId
					transfer.CreatedAt = testTime
					return errRepository
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
			args: args{
				amount:        testMoney100,
				originId:      account.AccountId(testUUID1),
				destinationId: account.AccountId(testUUID2),
			},
			want: transfer.Transfer{},
			err:  ErrCreateTransfer,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository, tt.crediter, tt.debiter}

			newTransfer, err := uc.Create(context.Background(), tt.args.amount, tt.args.originId, tt.args.destinationId)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected %v", err, tt.err)
			}

			if !reflect.DeepEqual(newTransfer, transfer.Transfer{}) {
				newTransfer.Id = transfer.TransferId(testUUID3)
				newTransfer.CreatedAt = testTime
			}

			if !reflect.DeepEqual(newTransfer, tt.want) {
				t.Errorf("got %v expected %v", newTransfer, tt.want)
			}
		})
	}
}
