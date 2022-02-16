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
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository transfer.Repository
		debiter    Debiter
		crediter   Crediter
		want       transfer.Transfer
		err        error
	}

	testMoney100, _ := money.NewMoney(100)

	testTransfer := transfer.Transfer{
		ID:                   transfer.ID(uuid.New()),
		AccountOriginID:      account.ID(uuid.New()),
		AccountDestinationID: account.ID(uuid.New()),
		Amount:               testMoney100,
		CreatedAt:            time.Now(),
	}

	testCases := []testCase{
		{
			name: "create transfer",
			repository: transfer.MockRepository{
				OnCreate: func(ctx context.Context, transferInstance transfer.Transfer) (transfer.Transfer, error) {
					return testTransfer, nil
				},
			},
			debiter: MockDebiter{
				OnDebitAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return nil
				},
			},
			crediter: MockCrediter{
				OnCreditAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return nil
				},
			},
			want: testTransfer,
		},
		{
			name: "fail to debit from origin",
			debiter: MockDebiter{
				OnDebitAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return accountuc.ErrRepository
				},
			},
			err: accountuc.ErrRepository,
		},
		{
			name: "fail to credit to destination",
			debiter: MockDebiter{
				OnDebitAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return nil
				},
			},
			crediter: MockCrediter{
				OnCreditAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return accountuc.ErrRepository
				},
			},
			err: accountuc.ErrRepository,
		},
		{
			name: "repository error creating transfer",
			repository: transfer.MockRepository{
				OnCreate: func(ctx context.Context, transferInstance transfer.Transfer) (transfer.Transfer, error) {
					return transfer.Transfer{}, ErrRepository
				},
			},
			debiter: MockDebiter{
				OnDebitAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return nil
				},
			},
			crediter: MockCrediter{
				OnCreditAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return nil
				},
			},
			err: ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository, tt.crediter, tt.debiter}

			gotTransfer, err := uc.Create(context.Background(), tt.want)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected %v", err, tt.err)
			}

			if !reflect.DeepEqual(gotTransfer, tt.want) {
				t.Errorf("got transfer %v expected transfer %v", gotTransfer, tt.want)
			}
		})
	}
}
