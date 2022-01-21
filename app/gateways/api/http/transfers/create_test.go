package transfersroute

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
)

func TestCreate(t *testing.T) {
	t.Parallel()

	type httpIO struct {
		r *http.Request
		w http.ResponseWriter
	}

	type testCase struct {
		name            string
		usecase         transferuc.Usecase
		httpIO          httpIO
		expectedStatus  int
		expectedPayload interface{}
	}

	var (
		testTransferId = transfer.TransferId(uuid.New())
		testTime       = time.Now()
		testMoney10, _ = money.NewMoney(10)
	)

	testAccount1 := account.Account{
		Id: account.AccountId(uuid.New()),
		// Name:      "nice name",
		// Cpf:       cpf.Random(),
		Balance: testMoney10,
		// CreatedAt: time.Now(),
	}

	testAccount2 := account.Account{
		Id: account.AccountId(uuid.New()),
		// Name:      "nice name",
		// Cpf:       cpf.Random(),
		// Balance: testMoney10,
		// CreatedAt: time.Now(),
	}

	//TODO: fix requestPayload
	requestPayload := fmt.Sprintf(`{"account_destination_id":%s,"amount":5}`, testAccount2.Id.String())
	originIdToken, _ := token.Sign(testAccount1.Id)
	requestHeader := fmt.Sprintf(`"Bearer %s"`, originIdToken.Value())

	testCases := []testCase{
		{
			name: "successfully transfer",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						bytes.NewReader([]byte(requestPayload)))
					request.Header.Set("Authorization", requestHeader)
					return request
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: transferuc.MockUsecase{
				OnCreate: func(ctx context.Context, amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {
					return transfer.Transfer{
						Id:                   testTransferId,
						AccountOriginId:      testAccount1.Id,
						AccountDestinationId: testAccount2.Id,
						Amount:               testMoney10,
						CreatedAt:            testTime,
					}, nil
				},
			},
			expectedStatus: 201,
			expectedPayload: map[string]interface{}{
				"transfer_id":            testTransferId.String(),
				"account_origin_id":      testAccount1.Id.String(),
				"account_destination_id": testAccount2.Id.String(),
				"amount":                 5,
				"created_at":             testTime,
			},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.usecase)

			h.Create(tt.httpIO.w, tt.httpIO.r)

			recorder, ok := tt.httpIO.w.(*httptest.ResponseRecorder)
			if !ok {
				t.Errorf("Error getting ResponseRecorder")
			}

			if statusCode := recorder.Code; statusCode != tt.expectedStatus {
				t.Errorf("got status code %d expected %d", statusCode, tt.expectedStatus)
			}

			var got map[string]interface{}
			err := json.NewDecoder(recorder.Body).Decode(&got)

			if err != nil {
				t.Fatalf("failed to decode response body: %s", err)
			}

			if !reflect.DeepEqual(got, tt.expectedPayload) {
				t.Fatalf("got response body: %s, expected %s", got, tt.expectedPayload)
			}
		})
	}
}
