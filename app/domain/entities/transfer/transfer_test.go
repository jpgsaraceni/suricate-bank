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
		originId      account.AccountId
		destinationId account.AccountId
		amount        money.Money
	}

	type testCase struct {
		name string
		args args
		want Transfer
		err  error
	}

	var testUUID1, _ = uuid.NewUUID()
	var testUUID2, _ = uuid.NewUUID()
	var testUUID3, _ = uuid.NewUUID()

	var testMoney100, _ = money.NewMoney(100)
	// var testMoney200, _ = money.NewMoney(200)
	// var testMoney0, _ = money.NewMoney(0)

	var testTime = time.Now()

	testCases := []testCase{
		{
			name: "makes transfer",
			args: args{
				originId:      account.AccountId(testUUID1),
				destinationId: account.AccountId(testUUID2),
				amount:        testMoney100,
			},
			want: Transfer{
				Id:                   TransferId(testUUID3),
				AccountOriginId:      account.AccountId(testUUID1),
				AccountDestinationId: account.AccountId(testUUID2),
				Amount:               testMoney100,
				CreatedAt:            testTime,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewTransfer(tt.args.amount, tt.args.originId, tt.args.destinationId)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			got.CreatedAt = testTime
			got.Id = TransferId(testUUID3)

			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %v expected %v", got, tt.want)
			}
		})
	}
}
