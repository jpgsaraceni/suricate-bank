package accountsroute

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	accountspg "github.com/jpgsaraceni/suricate-bank/app/gateways/db/postgres/accounts"
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
		expectedPayload responses.Payload
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
					return account.Account{
						Id: account.AccountId(uuid.New()),
					}, nil
				},
			},
			expectedStatus:  201,
			expectedPayload: responses.AccountCreated,
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
			expectedPayload: responses.ErrInvalidRequestPayload.Payload,
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
			expectedPayload: responses.ErrMissingFields.Payload,
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
			expectedPayload: responses.ErrShortName.Payload,
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
			expectedPayload: responses.ErrMissingFields.Payload,
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
			expectedPayload: responses.ErrLengthCpf.Payload,
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
			expectedPayload: responses.ErrInvalidCpf.Payload,
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
			expectedPayload: responses.ErrMissingFields.Payload,
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
			expectedPayload: responses.ErrShortSecret.Payload,
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
					return account.Account{}, accountspg.ErrCpfAlreadyExists
				},
			},
			expectedStatus:  400,
			expectedPayload: responses.ErrCpfAlreadyExists.Payload,
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
					return account.Account{}, accountuc.ErrCreateAccount
				},
			},
			expectedStatus:  500,
			expectedPayload: responses.ErrInternalServerError,
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

			var got responses.Payload
			err := json.NewDecoder(recorder.Body).Decode(&got)

			if err != nil {
				t.Fatalf("failed to decode response body: %s", err)
			}

			if got != tt.expectedPayload {
				t.Fatalf("got response body: %s, expected %s", got, tt.expectedPayload)
			}
		})
	}
}
