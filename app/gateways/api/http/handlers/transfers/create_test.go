package transfersroute

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/middlewares"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
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
		Id:      account.AccountId(uuid.New()),
		Balance: testMoney10,
	}

	testAccount2 := account.Account{
		Id: account.AccountId(uuid.New()),
	}

	var (
		requestPayload               = fmt.Sprintf(`{"account_destination_id":"%s","amount":5}`, testAccount2.Id.String())
		requestPayloadZeroAmount     = fmt.Sprintf(`{"account_destination_id":"%s","amount":0}`, testAccount2.Id.String())
		requestPayloadNegativeAmount = fmt.Sprintf(`{"account_destination_id":"%s","amount":-10}`, testAccount2.Id.String())
		requestPayloadRepeatedId     = fmt.Sprintf(`{"account_destination_id":"%s","amount":5}`, testAccount1.Id.String())
	)

	testCases := []testCase{
		{
			name: "successfully transfer",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(requestPayload),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))
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
				"amount":                 testMoney10.BRL(),
				"created_at":             testTime.Format(time.RFC3339Nano),
			},
		},
		{
			name: "respond 400 to invalid request payload",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						strings.NewReader(`Not what the server expects`))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus: 400,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrInvalidCreateTransferPayload.Payload.Message,
			},
		},
		{
			name: "respond 400 to request with empty destination account",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(`{"account_destination_id":"","amount":5}`),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))
					return request
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus: 400,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrMissingFieldsTransferPayload.Payload.Message,
			},
		},
		{
			name: "respond 400 to request with zero amount",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(requestPayloadZeroAmount),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))
					return request
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus: 400,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrMissingFieldsTransferPayload.Payload.Message,
			},
		},
		{
			name: "respond 400 to request with negative amount",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(requestPayloadNegativeAmount),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))
					return request
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus: 400,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrInvalidAmount.Payload.Message,
			},
		},
		{
			name: "respond 400 to request with same origin and destination ids",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(requestPayloadRepeatedId),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))
					return request
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus: 400,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrSameAccounts.Payload.Message,
			},
		},
		{
			name: "respond 422 to request with insufficient funds in origin account",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(requestPayload),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))

					return request
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: transferuc.MockUsecase{
				OnCreate: func(ctx context.Context, amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {
					return transfer.Transfer{}, account.ErrInsufficientFunds
				},
			},
			expectedStatus: 422,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrInsuficientFunds.Payload.Message,
			},
		},
		{
			name: "respond 404 to request with inexistent account id",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(requestPayload),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))

					return request
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: transferuc.MockUsecase{
				OnCreate: func(ctx context.Context, amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {
					return transfer.Transfer{}, account.ErrIdNotFound
				},
			},
			expectedStatus: 404,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrAccountNotFound.Payload.Message,
			},
		},
		{
			name: "respond 500 on usecase error",
			httpIO: httpIO{
				r: func() *http.Request {
					request := httptest.NewRequest(
						http.MethodPost,
						"/transfers",
						strings.NewReader(requestPayload),
					)
					request = request.WithContext(middlewares.WithOriginId(context.Background(), testAccount1.Id))

					return request
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: transferuc.MockUsecase{
				OnCreate: func(ctx context.Context, amount money.Money, originId, destinationId account.AccountId) (transfer.Transfer, error) {
					return transfer.Transfer{}, accountuc.ErrRepository
				},
			},
			expectedStatus: 500,
			expectedPayload: map[string]interface{}{
				"title": responses.ErrInternalServerError.Payload.Message,
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
