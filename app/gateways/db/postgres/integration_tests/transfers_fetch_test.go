package postgres_test

import (
	"errors"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
	transferspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/transfers"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestFetchTransfers(t *testing.T) {
	// can't run this test in parallel because Fetch would possibly return accounts
	// created in parallel tests after the truncate function call

	var (
		accountId1     = account.AccountId(uuid.New())
		accountId2     = account.AccountId(uuid.New())
		accountId3     = account.AccountId(uuid.New())
		accountId4     = account.AccountId(uuid.New())
		cpf1           = cpf.Random()
		cpf2           = cpf.Random()
		cpf3           = cpf.Random()
		cpf4           = cpf.Random()
		transferId1    = transfer.TransferId(uuid.New())
		transferId2    = transfer.TransferId(uuid.New())
		testMoney10, _ = money.NewMoney(10)
		testMoney20, _ = money.NewMoney(20)
		testTime       = time.Now().Round(time.Hour)
		testHash, _    = hash.NewHash("nicesecret")
	)

	type testCase struct {
		name              string
		runBefore         func(*accountspg.Repository, *transferspg.Repository) error
		expectedTransfers []transfer.Transfer
		err               error
	}

	testCases := []testCase{
		{
			name: "successfully fetch 2 transfers",
			runBefore: func(accountsRepo *accountspg.Repository, transfersRepo *transferspg.Repository) error {
				truncateAccounts()
				err := accountsRepo.Create(
					testContext,
					&account.Account{
						Id:        accountId1,
						Name:      "Nice name",
						Cpf:       cpf1,
						Secret:    testHash,
						CreatedAt: testTime,
						Balance:   testMoney10,
					},
				)

				if err != nil {

					return err
				}

				err = accountsRepo.Create(
					testContext,
					&account.Account{
						Id:        accountId2,
						Name:      "Another nice name",
						Cpf:       cpf2,
						Secret:    testHash,
						CreatedAt: testTime,
					},
				)

				if err != nil {

					return err
				}

				err = accountsRepo.Create(
					testContext,
					&account.Account{
						Id:        accountId3,
						Name:      "Yet another nice name",
						Cpf:       cpf3,
						Secret:    testHash,
						CreatedAt: testTime,
						Balance:   testMoney20,
					},
				)

				if err != nil {

					return err
				}

				err = accountsRepo.Create(
					testContext,
					&account.Account{
						Id:        accountId4,
						Name:      "A not so nice name",
						Cpf:       cpf4,
						Secret:    testHash,
						CreatedAt: testTime,
					},
				)

				if err != nil {

					return err
				}

				err = transfersRepo.Create(
					testContext,
					&transfer.Transfer{
						Id:                   transferId1,
						AccountOriginId:      accountId1,
						AccountDestinationId: accountId2,
						Amount:               testMoney10,
						CreatedAt:            testTime,
					},
				)

				if err != nil {

					return err
				}

				err = transfersRepo.Create(
					testContext,
					&transfer.Transfer{
						Id:                   transferId2,
						AccountOriginId:      accountId3,
						AccountDestinationId: accountId4,
						Amount:               testMoney20,
						CreatedAt:            testTime,
					},
				)

				return err
			},
			expectedTransfers: []transfer.Transfer{
				{
					Id:                   transferId1,
					AccountOriginId:      accountId1,
					AccountDestinationId: accountId2,
					Amount:               testMoney10,
					CreatedAt:            testTime,
				},
				{
					Id:                   transferId2,
					AccountOriginId:      accountId3,
					AccountDestinationId: accountId4,
					Amount:               testMoney20,
					CreatedAt:            testTime,
				},
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Parallel()

		accountsRepo := accountspg.NewRepository(dbPool)
		transfersRepo := transferspg.NewRepository(dbPool)

		if tt.runBefore != nil {
			err := tt.runBefore(accountsRepo, transfersRepo)

			if err != nil {
				t.Fatalf("runBefore() failed: %s", err)
			}
		}

		gotTransfers, err := transfersRepo.Fetch(testContext)

		if !errors.Is(err, tt.err) {
			t.Fatalf("got error: %s expected error: %s", err, tt.err)
		}

		if !reflect.DeepEqual(gotTransfers, tt.expectedTransfers) {
			t.Fatalf("got %v expected %v", gotTransfers, tt.expectedTransfers)
		}
	}
}
