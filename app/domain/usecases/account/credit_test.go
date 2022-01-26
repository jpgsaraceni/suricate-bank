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

	var testAccountId = account.AccountId(uuid.New())

	testCases := []testCase{
		{
			name: "successfully credit 100 to account with 0 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney100,
		},
		{
			name: "successfully credit 100 to account with 10 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney10,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney100,
		},
		{
			name: "fail to credit 0 to account with 10 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney10,
			},
			amount: testMoney0,
			err:    ErrAmount,
		},
		{
			name: "fail to credit 0 to account with 0 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			amount: testMoney0,
			err:    ErrAmount,
		},
		{
			name: "fail to credit inexistent account",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return ErrIdNotFound
				},
			},
			amount: testMoney100,
			err:    ErrIdNotFound,
		},
		{
			name: "repository error",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnCreditAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
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

			err := uc.Credit(context.Background(), tt.testAccount.Id, tt.amount)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}
		})
	}
}
