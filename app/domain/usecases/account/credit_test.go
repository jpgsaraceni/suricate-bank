package accountuc

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestCredit(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		repository  account.Repository
		amount      money.Money
		testAccount account.Account
		err         error
	}

	var (
		testMoney100, _ = money.NewMoney(100)
		testMoney10, _  = money.NewMoney(10)
		testMoney0, _   = money.NewMoney(0)
	)

	testAccountID := account.ID(uuid.New())

	testCases := []testCase{
		{
			name: "successfully credit 100 to account with 0 initial balance",
			testAccount: account.Account{
				ID:      testAccountID,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney100,
		},
		{
			name: "successfully credit 100 to account with 10 initial balance",
			testAccount: account.Account{
				ID:      testAccountID,
				Balance: testMoney10,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney100,
		},
		{
			name: "fail to credit 0 to account with 10 initial balance",
			testAccount: account.Account{
				ID:      testAccountID,
				Balance: testMoney10,
			},
			amount: testMoney0,
			err:    ErrAmount,
		},
		{
			name: "fail to credit 0 to account with 0 initial balance",
			testAccount: account.Account{
				ID:      testAccountID,
				Balance: testMoney0,
			},
			amount: testMoney0,
			err:    ErrAmount,
		},
		{
			name: "fail to credit inexistent account",
			testAccount: account.Account{
				ID:      testAccountID,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return account.ErrIDNotFound
				},
			},
			amount: testMoney100,
			err:    account.ErrIDNotFound,
		},
		{
			name: "repository error",
			testAccount: account.Account{
				ID:      testAccountID,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.ID, amount money.Money) error {
					return errors.New("")
				},
			},
			amount: testMoney100,
			err:    ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository}

			err := uc.Credit(context.Background(), tt.testAccount.ID, tt.amount)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}
		})
	}
}
