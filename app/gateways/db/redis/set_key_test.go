package redis

import (
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis/redistest"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestSetKeyValue(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)

	t.Cleanup(tearDown)

	type testCase struct {
		name      string
		runBefore func()
		key       string
		res       responses.Response
		err       error
	}

	testAccount := func() account.Account {
		return account.Account{
			Id:        account.AccountId(uuid.New()),
			Name:      "nice name",
			Cpf:       cpf.Random(),
			Balance:   money.Money{},
			CreatedAt: time.Now(),
		}
	}
	testTransfer := func(amount int) transfer.Transfer {
		amountTransfered, _ := money.NewMoney(amount)
		return transfer.Transfer{
			Id:                   transfer.TransferId(uuid.New()),
			AccountOriginId:      testAccount().Id,
			AccountDestinationId: testAccount().Id,
			Amount:               amountTransfered,
			CreatedAt:            time.Now(),
		}
	}

	testAccounts := []account.Account{
		testAccount(),
		testAccount(),
		testAccount(),
		testAccount(),
	}

	testTransfers := []transfer.Transfer{
		testTransfer(10),
		testTransfer(5),
		testTransfer(100),
	}
	repeatedKey := uuid.NewString()

	testCases := []testCase{
		{
			name: "set a created account response",
			key:  uuid.NewString(),
			res: responses.Response{
				Status: 201,
				Payload: map[string]interface{}{
					"account_id": testAccounts[0].Id.String(),
					"name":       testAccounts[0].Name,
					"cpf":        testAccounts[0].Cpf.Masked(),
					"balance":    testAccounts[0].Balance.BRL(),
					"created_at": testAccounts[0].CreatedAt.Format(time.RFC3339Nano),
				},
			},
		},
		{
			name: "set a create account error response",
			key:  uuid.NewString(),
			res: responses.Response{
				Status:  400,
				Payload: map[string]interface{}{"title": responses.ErrInvalidCreateAccountPayload.Payload.Message},
			},
		},
		{
			name: "set a fetched accounts response",
			key:  uuid.NewString(),
			res: responses.Response{
				Status: 200,
				Payload: map[string]interface{}{
					"accounts": []interface{}{
						map[string]interface{}{
							"account_id": testAccounts[0].Id.String(),
							"name":       testAccounts[0].Name,
							"cpf":        testAccounts[0].Cpf.Masked(),
							"balance":    testAccounts[0].Balance.BRL(),
							"created_at": testAccounts[0].CreatedAt.Format(time.RFC3339Nano),
						},
						map[string]interface{}{
							"account_id": testAccounts[1].Id.String(),
							"name":       testAccounts[1].Name,
							"cpf":        testAccounts[1].Cpf.Masked(),
							"balance":    testAccounts[1].Balance.BRL(),
							"created_at": testAccounts[1].CreatedAt.Format(time.RFC3339Nano),
						},
						map[string]interface{}{
							"account_id": testAccounts[2].Id.String(),
							"name":       testAccounts[2].Name,
							"cpf":        testAccounts[2].Cpf.Masked(),
							"balance":    testAccounts[2].Balance.BRL(),
							"created_at": testAccounts[2].CreatedAt.Format(time.RFC3339Nano),
						},
					},
				},
			},
		},
		{
			name: "set a get account balance response",
			key:  uuid.NewString(),
			res: responses.Response{
				Status: 200,
				Payload: map[string]interface{}{
					"account_id": testAccounts[0].Id.String(),
					"balance":    "R$0,10",
				},
			},
		},
		{
			name: "set a created transfer response",
			key:  uuid.NewString(),
			res: responses.Response{
				Status: 201,
				Payload: map[string]interface{}{
					"account_id": testAccounts[0].Id.String(),
					"name":       testAccounts[0].Name,
					"cpf":        testAccounts[0].Cpf.Masked(),
					"balance":    testAccounts[0].Balance.BRL(),
					"created_at": testAccounts[0].CreatedAt.Format(time.RFC3339Nano),
				},
			},
		},
		{
			name: "set a fetched transfers response",
			key:  uuid.NewString(),
			res: responses.Response{
				Status: 200,
				Payload: map[string]interface{}{
					"transfers": []interface{}{
						map[string]interface{}{
							"transfer_id":            testTransfers[0].Id.String(),
							"account_origin_id":      testTransfers[0].AccountOriginId.String(),
							"account_destination_id": testTransfers[0].AccountDestinationId.String(),
							"amount":                 testTransfers[0].Amount.BRL(),
							"created_at":             testTransfers[0].CreatedAt.Format(time.RFC3339Nano),
						},
						map[string]interface{}{
							"transfer_id":            testTransfers[1].Id.String(),
							"account_origin_id":      testTransfers[1].AccountOriginId.String(),
							"account_destination_id": testTransfers[1].AccountDestinationId.String(),
							"amount":                 testTransfers[1].Amount.BRL(),
							"created_at":             testTransfers[1].CreatedAt.Format(time.RFC3339Nano),
						},
						map[string]interface{}{
							"transfer_id":            testTransfers[2].Id.String(),
							"account_origin_id":      testTransfers[2].AccountOriginId.String(),
							"account_destination_id": testTransfers[2].AccountDestinationId.String(),
							"amount":                 testTransfers[2].Amount.BRL(),
							"created_at":             testTransfers[2].CreatedAt.Format(time.RFC3339Nano),
						},
					},
				},
			},
		},
		{
			name: "fail to set existent key",
			runBefore: func() {
				testRepo.SetKeyValue(repeatedKey, responses.Response{Status: 200})
			},
			key: repeatedKey,
			res: responses.Response{
				Status: 201,
				Payload: map[string]interface{}{
					"account_id": testAccounts[0].Id.String(),
					"name":       testAccounts[0].Name,
					"cpf":        testAccounts[0].Cpf.Masked(),
					"balance":    testAccounts[0].Balance.BRL(),
					"created_at": testAccounts[0].CreatedAt.Format(time.RFC3339Nano),
				},
			},
			err: errKeyExists,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				tt.runBefore()
			}

			err := testRepo.SetKeyValue(tt.key, tt.res)

			if !errors.Is(err, tt.err) {

				t.Fatalf("\ngot error: \n%s\nexpected error: \n%s\n", err, tt.err)
			}
		})
	}
}
