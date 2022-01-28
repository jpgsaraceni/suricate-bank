package accountuc

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestDebit(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		repository  account.Repository
		amount      money.Money
		testAccount account.Account
		want        account.Account
		err         error
	}

	var (
		testMoney100, _ = money.NewMoney(100)
		testMoney10, _  = money.NewMoney(10)
		testMoney110, _ = money.NewMoney(110)
		testMoney0, _   = money.NewMoney(0)
	)

	var testAccountId = account.AccountId(uuid.New())

	testCases := []testCase{
		{
			name: "successfully debit 100 from account with 100 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney100,
			},
			repository: account.MockRepository{
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney100,
			want: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
		},
		{
			name: "successfully debit 10 from account with 110 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney110,
			},
			repository: account.MockRepository{
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney10,
			want: account.Account{
				Id:      testAccountId,
				Balance: testMoney100,
			},
		},
		{
			name: "fail to debit 0 from account with 10 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney10,
			},
			amount: testMoney0,
			want: account.Account{
				Id:      testAccountId,
				Balance: testMoney10,
			},
			err: ErrAmount,
		},
		{
			name: "fail to debit 0 from account with 0 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			amount: testMoney0,
			want: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			err: ErrAmount,
		},
		{
			name: "fail to debit inexistent account",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return account.ErrIdNotFound
				},
			},
			amount: testMoney100,
			want:   account.Account{},
			err:    account.ErrIdNotFound,
		},
		{
			name: "fail to debit 10 from account with 0 initial balance",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return account.ErrInsufficientFunds
				},
			},
			amount: testMoney10,
			want: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			err: account.ErrInsufficientFunds,
		},
		{
			name: "repository error",
			testAccount: account.Account{
				Id:      testAccountId,
				Balance: testMoney100,
			},
			repository: account.MockRepository{
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return errors.New("")
				},
			},
			amount: testMoney100,
			want: account.Account{
				Id:      testAccountId,
				Balance: testMoney0,
			},
			err: ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := usecase{tt.repository}

			err := uc.Debit(context.Background(), tt.testAccount.Id, tt.amount)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}
		})
	}
}
