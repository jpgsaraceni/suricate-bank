package accountsroute

import (
	"bytes"
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
		usecase         accountuc.Usecase
		httpIO          httpIO
		expectedStatus  int
		expectedPayload interface{}
	}

	testAccount := account.Account{
		Id:        account.AccountId(uuid.New()),
		Name:      "nice name",
		Cpf:       cpf.Random(),
		Balance:   money.Money{},
		CreatedAt: time.Now(),
	}

	testCases := []testCase{
		{
			name: "successfully create account",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return testAccount, nil
				},
			},
			expectedStatus: 201,
			expectedPayload: map[string]interface{}{
				"account_id": testAccount.Id.String(),
				"name":       testAccount.Name,
				"cpf":        testAccount.Cpf.Masked(),
				"balance":    testAccount.Balance.BRL(),
				"created_at": testAccount.CreatedAt.Format(time.RFC3339Nano),
			},
		},
		{
			name: "respond 400 to invalid request payload",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`Not what the server expects`)))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrInvalidCreateAccountPayload.Payload.Message},
		},
		{
			name: "respond 400 to request missing name",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrMissingFieldsAccountPayload.Payload.Message},
		},
		{
			name: "respond 400 to request with short name",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"N", "Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrShortName.Payload.Message},
		},
		{
			name: "respond 400 to request with missing cpf",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrMissingFieldsAccountPayload.Payload.Message},
		},
		{
			name: "respond 400 to request with invalid cpf length",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"123456789000", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrLengthCpf.Payload.Message},
		},
		{
			name: "respond 400 when cpf validation fails",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-34", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{}, account.ErrInvalidCpf
				},
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrInvalidCpf.Payload.Message},
		},
		{
			name: "respond 400 to request missing secret",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrMissingFieldsAccountPayload.Payload.Message},
		},
		{
			name: "respond 400 to request with short secret",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35", "Secret": "12345"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrShortSecret.Payload.Message},
		},
		{
			name: "respond 400 to request containing cpf that already exists",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{}, accountuc.ErrDuplicateCpf
				},
			},
			expectedStatus:  400,
			expectedPayload: map[string]interface{}{"title": responses.ErrCpfAlreadyExists.Payload.Message},
		},
		{
			name: "respond 500 due to usecase error",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{}, accountuc.ErrRepository
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
