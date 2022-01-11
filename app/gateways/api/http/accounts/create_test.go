package accountsroute

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	accountuc "github.com/jpgsaraceni/suricate-bank/app/domain/usecases/account"
)

func TestCreate(t *testing.T) { // TODO: assert payloads
	t.Parallel()

	type httpIO struct {
		r *http.Request
		w http.ResponseWriter
	}

	type testCase struct {
		name           string
		usecase        accountuc.Usecase
		httpIO         httpIO
		expectedStatus int
	}

	var (
		testId = account.AccountId(uuid.New())
	)

	testCases := []testCase{
		{
			name: "successfully create account",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{
						Id: testId,
					}, nil
				},
			},
			expectedStatus: 201,
		},
		{
			name: "fail to create account missing name",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						bytes.NewReader([]byte(`{"Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{
						Id: testId,
					}, nil
				},
			},
			expectedStatus: 400,
		},
		{
			name: "fail to create account with short name",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"N", "Cpf":"220.614.460-35", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{
						Id: testId,
					}, nil
				},
			},
			expectedStatus: 400,
		},
		{
			name: "fail to create account missing cpf",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{
						Id: testId,
					}, nil
				},
			},
			expectedStatus: 400,
		},
		{
			name: "fail to create account with invalid cpf",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
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
			expectedStatus: 400,
		},
		{
			name: "fail to create account missing secret",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35"`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{
						Id: testId,
					}, nil
				},
			},
			expectedStatus: 400,
		},
		{
			name: "fail to create account short secret",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
						"/accounts",
						bytes.NewReader([]byte(`{"Name":"Nice Name", "Cpf":"220.614.460-35", "Secret": "12345"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			usecase: accountuc.MockUsecase{
				OnCreate: func(ctx context.Context, name, cpf, secret string) (account.Account, error) {
					return account.Account{
						Id: testId,
					}, nil
				},
			},
			expectedStatus: 400,
		},
		{
			name: "fail to create account due to usecase error",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodGet,
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
			expectedStatus: 500,
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
		})
	}
}
