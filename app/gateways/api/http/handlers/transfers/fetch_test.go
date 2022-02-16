package transfersroute

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/transfer"
	transferuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/transfer"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/money"
)

func TestFetch(t *testing.T) {
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
		testMoney10, _ = money.NewMoney(10)
		testMoney20, _ = money.NewMoney(50)
	)

	testAccount1 := account.Account{
		ID: account.ID(uuid.New()),
	}
	testAccount2 := account.Account{
		ID: account.ID(uuid.New()),
	}
	testAccount3 := account.Account{
		ID: account.ID(uuid.New()),
	}
	testAccount4 := account.Account{
		ID: account.ID(uuid.New()),
	}

	testTransfer1 := transfer.Transfer{
		ID:                   transfer.ID(uuid.New()),
		AccountOriginID:      testAccount1.ID,
		AccountDestinationID: testAccount2.ID,
		Amount:               testMoney20,
		CreatedAt:            time.Now(),
	}
	testTransfer2 := transfer.Transfer{
		ID:                   transfer.ID(uuid.New()),
		AccountOriginID:      testAccount1.ID,
		AccountDestinationID: testAccount3.ID,
		Amount:               testMoney10,
		CreatedAt:            time.Now(),
	}
	testTransfer3 := transfer.Transfer{
		ID:                   transfer.ID(uuid.New()),
		AccountOriginID:      testAccount4.ID,
		AccountDestinationID: testAccount3.ID,
		Amount:               testMoney10,
		CreatedAt:            time.Now(),
	}

	testCases := []testCase{
		{
			name: "successfully fetch transfers",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/transfers",
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: transferuc.MockUsecase{
				OnFetch: func(ctx context.Context) ([]transfer.Transfer, error) {
					return []transfer.Transfer{
						testTransfer1,
						testTransfer2,
						testTransfer3,
					}, nil
				},
			},
			expectedStatus: 200,
			expectedPayload: map[string]interface{}{
				"transfers": []interface{}{
					map[string]interface{}{
						"transfer_id":            testTransfer1.ID.String(),
						"account_origin_id":      testAccount1.ID.String(),
						"account_destination_id": testAccount2.ID.String(),
						"amount":                 testTransfer1.Amount.BRL(),
						"created_at":             testTransfer1.CreatedAt.Format(time.RFC3339Nano),
					},
					map[string]interface{}{
						"transfer_id":            testTransfer2.ID.String(),
						"account_origin_id":      testAccount1.ID.String(),
						"account_destination_id": testAccount3.ID.String(),
						"amount":                 testTransfer2.Amount.BRL(),
						"created_at":             testTransfer2.CreatedAt.Format(time.RFC3339Nano),
					},
					map[string]interface{}{
						"transfer_id":            testTransfer3.ID.String(),
						"account_origin_id":      testAccount4.ID.String(),
						"account_destination_id": testAccount3.ID.String(),
						"amount":                 testTransfer3.Amount.BRL(),
						"created_at":             testTransfer3.CreatedAt.Format(time.RFC3339Nano),
					},
				},
			},
		},
		{
			name: "successfully fetch 0 accounts",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/transfers",
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: transferuc.MockUsecase{
				OnFetch: func(ctx context.Context) ([]transfer.Transfer, error) {
					return []transfer.Transfer{}, nil
				},
			},
			expectedStatus: 200,
			expectedPayload: map[string]interface{}{
				"transfers": []interface{}{},
			},
		},
		{
			name: "fail due to repository error",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/transfers",
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: transferuc.MockUsecase{
				OnFetch: func(ctx context.Context) ([]transfer.Transfer, error) {
					return []transfer.Transfer{}, transferuc.ErrRepository
				},
			},
			expectedStatus:  500,
			expectedPayload: map[string]interface{}{"title": responses.ErrInternalServerError.Payload.Message},
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			h := NewHandler(tt.usecase)

			h.Fetch(tt.httpIO.w, tt.httpIO.r)

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
				t.Fatalf("\ngot response body:\n %s\n expected response body:\n %s", got["accounts"], tt.expectedPayload)
			}
		})
	}
}
