package accountsroute

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
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/vos/cpf"
	"github.com/jpgsaraceni/suricate-bank/app/vos/hash"
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
		usecase         accountuc.Usecase
		httpIO          httpIO
		expectedStatus  int
		expectedPayload interface{}
	}

	testAccount1 := account.Account{
		Id:        account.AccountId(uuid.New()),
		Name:      "nice name",
		Cpf:       cpf.Random(),
		Secret:    hash.Parse("123456"),
		Balance:   money.Money{},
		CreatedAt: time.Now(),
	}
	testAccount2 := account.Account{
		Id:        account.AccountId(uuid.New()),
		Name:      "nice name",
		Cpf:       cpf.Random(),
		Secret:    hash.Parse("123456"),
		Balance:   money.Money{},
		CreatedAt: time.Now(),
	}
	testAccount3 := account.Account{
		Id:        account.AccountId(uuid.New()),
		Name:      "nice name",
		Cpf:       cpf.Random(),
		Secret:    hash.Parse("123456"),
		Balance:   money.Money{},
		CreatedAt: time.Now(),
	}

	testCases := []testCase{
		{
			name: "successfully fetch accounts",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {
					return []account.Account{
						{
							Id:        testAccount1.Id,
							Name:      testAccount1.Name,
							Cpf:       testAccount1.Cpf,
							Secret:    testAccount1.Secret,
							Balance:   testAccount1.Balance,
							CreatedAt: testAccount1.CreatedAt,
						},
						{
							Id:        testAccount2.Id,
							Name:      testAccount2.Name,
							Cpf:       testAccount2.Cpf,
							Secret:    testAccount2.Secret,
							Balance:   testAccount2.Balance,
							CreatedAt: testAccount2.CreatedAt,
						},
						{
							Id:        testAccount3.Id,
							Name:      testAccount3.Name,
							Cpf:       testAccount3.Cpf,
							Secret:    testAccount3.Secret,
							Balance:   testAccount3.Balance,
							CreatedAt: testAccount3.CreatedAt,
						},
					}, nil
				},
			},
			expectedStatus: 200,
			expectedPayload: map[string]interface{}{
				"accounts": []interface{}{
					map[string]interface{}{
						"account_id": testAccount1.Id.String(),
						"name":       testAccount1.Name,
						"cpf":        testAccount1.Cpf.Masked(),
						"balance":    testAccount1.Balance.BRL(),
						"created_at": testAccount1.CreatedAt.Format(time.RFC3339Nano),
					},
					map[string]interface{}{
						"account_id": testAccount2.Id.String(),
						"name":       testAccount2.Name,
						"cpf":        testAccount2.Cpf.Masked(),
						"balance":    testAccount2.Balance.BRL(),
						"created_at": testAccount2.CreatedAt.Format(time.RFC3339Nano),
					},
					map[string]interface{}{
						"account_id": testAccount3.Id.String(),
						"name":       testAccount3.Name,
						"cpf":        testAccount3.Cpf.Masked(),
						"balance":    testAccount3.Balance.BRL(),
						"created_at": testAccount3.CreatedAt.Format(time.RFC3339Nano),
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
						"/accounts",
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {
					return []account.Account{}, nil
				},
			},
			expectedStatus: 200,
			expectedPayload: map[string]interface{}{
				"accounts": []interface{}{},
			},
		},
		{
			name: "fail due to repository error",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						nil,
					)
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnFetch: func(ctx context.Context) ([]account.Account, error) {
					return []account.Account{}, accountuc.ErrFetchAccounts
				},
			},
			expectedStatus:  500,
			expectedPayload: map[string]interface{}{"title": responses.ErrInternalServerError.Message},
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
