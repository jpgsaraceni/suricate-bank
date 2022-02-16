package transfer

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestNewTransfer(t *testing.T) {
	t.Parallel()

	type args struct {
		originID      account.ID
		destinationID account.ID
		amount        money.Money
	}

	type testCase struct {
		name string
		args args
		want Transfer
		err  error
	}

	testUUID1, _ := uuid.NewUUID()
	testUUID2, _ := uuid.NewUUID()
	testUUID3, _ := uuid.NewUUID()

	testMoney100, _ := money.NewMoney(100)
	testMoney0, _ := money.NewMoney(0)

	testTime := time.Now()

	testCases := []testCase{
		{
			name: "makes transfer",
			args: args{
				originID:      account.ID(testUUID1),
				destinationID: account.ID(testUUID2),
				amount:        testMoney100,
			},
			want: Transfer{
				ID:                   ID(testUUID3),
				AccountOriginID:      account.ID(testUUID1),
				AccountDestinationID: account.ID(testUUID2),
				Amount:               testMoney100,
				CreatedAt:            testTime,
			},
		},
		{
			name: "fails transfer when origin and destination are the same",
			args: args{
				originID:      account.ID(testUUID1),
				destinationID: account.ID(testUUID1),
				amount:        testMoney100,
			},
			want: Transfer{},
			err:  ErrSameAccounts,
		},
		{
			name: "fails transfer when amount is zero",
			args: args{
				originID:      account.ID(testUUID1),
				destinationID: account.ID(testUUID2),
				amount:        testMoney0,
			},
			want: Transfer{},
			err:  ErrAmountNotPositive,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewTransfer(tt.args.amount, tt.args.originID, tt.args.destinationID)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			if !reflect.DeepEqual(got, Transfer{}) {
				got.ID = ID(testUUID3)
				got.CreatedAt = testTime
			}

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v expected %v", got, tt.want)
			}
		})
	}
}
