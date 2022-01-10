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

	var testUUID, _ = uuid.NewUUID()
	var testUUID2, _ = uuid.NewUUID()

	var errRepository = errors.New("repository error")

	testCases := []testCase{
		{
			name: "successfully debit 100 from account with 100 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney100,
			},
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney100,
					}, nil
				},
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney100,
			want: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
		},
		{
			name: "successfully debit 10 from account with 110 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney110,
			},
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney110,
					}, nil
				},
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return nil
				},
			},
			amount: testMoney10,
			want: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney100,
			},
		},
		{
			name: "fail to debit 0 from account with 10 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney10,
			},
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney10,
					}, nil
				},
			},
			amount: testMoney0,
			want: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney10,
			},
			err: ErrAmount,
		},
		{
			name: "fail to debit 0 from account with 0 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney0,
					}, nil
				},
			},
			amount: testMoney0,
			want: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
			err: ErrAmount,
		},
		{
			name: "fail to debit inexistent account",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID2),
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{}, errRepository
				},
			},
			amount: testMoney100,
			want:   account.Account{},
			err:    ErrGetAccount,
		},
		{
			name: "fail to debit 10 from account with 0 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney0,
					}, nil
				},
			},
			amount: testMoney10,
			want: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
			err: ErrAmount,
		},
		{
			name: "repository error",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney100,
			},
			repository: account.MockRepository{
				OnGetById: func(ctx context.Context, id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney100,
					}, nil
				},
				OnDebitAccount: func(ctx context.Context, id account.AccountId, amount money.Money) error {
					return errRepository
				},
			},
			amount: testMoney100,
			want: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
			err: ErrDebitAccount,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository}

			err := uc.Debit(context.Background(), tt.testAccount.Id, tt.amount)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}
		})
	}
}
