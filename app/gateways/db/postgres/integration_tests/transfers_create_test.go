package postgres_test

import (
	"errors"
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

func TestCreateTransfer(t *testing.T) {

	type testCase struct {
		name      string
		runBefore func(repo *accountspg.Repository) error
		transfer  *transfer.Transfer
		err       error
	}

	var (
		accountId1     = account.AccountId(uuid.New())
		accountId2     = account.AccountId(uuid.New())
		testHash, _    = hash.NewHash("nicesecret")
		testMoney10, _ = money.NewMoney(10)
		testMoney20, _ = money.NewMoney(20)
	)

	testCases := []testCase{
		{
			name: "successfully create transfer",
			runBefore: func(repo *accountspg.Repository) error {
				err := repo.Create(
					testContext,
					&account.Account{
						Id:      accountId1,
						Name:    "Nice name",
						Cpf:     cpf.Random(),
						Secret:  testHash,
						Balance: testMoney20,
					},
				)

				if err != nil {

					return err
				}

				return repo.Create(
					testContext,
					&account.Account{
						Id:     accountId2,
						Name:   "Another nice name",
						Cpf:    cpf.Random(),
						Secret: testHash,
					},
				)
			},
			transfer: &transfer.Transfer{
				Id:                   transfer.TransferId(uuid.New()),
				AccountOriginId:      accountId1,
				AccountDestinationId: accountId2,
				Amount:               testMoney10,
				CreatedAt:            time.Now(),
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Parallel()

		accountsRepo := accountspg.NewRepository(dbPool)
		transfersRepo := transferspg.NewRepository(dbPool)

		if tt.runBefore != nil {
			err := tt.runBefore(accountsRepo)

			if err != nil {
				t.Fatalf("runBefore() failed: %s", err)
			}
		}

		if err := transfersRepo.Create(testContext, tt.transfer); !errors.Is(err, tt.err) {

			t.Fatalf("got error: %s expected error: %s", err, tt.err)
		}
	}
}
