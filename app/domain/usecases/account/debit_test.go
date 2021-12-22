package accountuc

import (
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

	testCases := []testCase{
		{
			name: "successfully debit 100 from account with 100 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney100,
			},
			repository: account.MockRepository{
				OnGetById: func(id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney100,
					}, nil
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
				OnGetById: func(id account.AccountId) (account.Account, error) {
					return account.Account{
						Id:      account.AccountId(testUUID),
						Balance: testMoney110,
					}, nil
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
				OnGetById: func(id account.AccountId) (account.Account, error) {
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
			err: money.ErrChangeByZero,
		},
		{
			name: "fail to debit 0 from account with 0 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnGetById: func(id account.AccountId) (account.Account, error) {
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
			err: money.ErrChangeByZero,
		},
		{
			name: "fail to debit inexistent account",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID2),
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnGetById: func(id account.AccountId) (account.Account, error) {
					return account.Account{}, ErrAccountNotFound
				},
			},
			amount: testMoney100,
			want:   account.Account{},
			err:    ErrAccountNotFound,
		},
		{
			name: "fail to debit 10 from account with 0 initial balance",
			testAccount: account.Account{
				Id:      account.AccountId(testUUID),
				Balance: testMoney0,
			},
			repository: account.MockRepository{
				OnGetById: func(id account.AccountId) (account.Account, error) {
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
			err: money.ErrInsuficientFunds,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			uc := Usecase{tt.repository}

			err := uc.Debit(tt.testAccount.Id, tt.amount)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}
		})
	}
}
