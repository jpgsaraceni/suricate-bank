package loginroute

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/app/gateways/api/http/responses"
	"github.com/jpgsaraceni/suricate-bank/app/services/auth"
	"github.com/jpgsaraceni/suricate-bank/app/vos/token"
)

func TestLogin(t *testing.T) {
	t.Parallel()

	type httpIO struct {
		r *http.Request
		w http.ResponseWriter
	}

	type testCase struct {
		name            string
		service         auth.Service
		httpIO          httpIO
		expectedStatus  int
		expectedPayload interface{}
	}

	var testId = account.AccountId(uuid.New())

	testCases := []testCase{
		{
			name: "successfully login",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Cpf":"22061446035", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			service: auth.MockService{
				OnAuthenticate: func(ctx context.Context, cpfInput, secret string) (string, error) {
					jwt, _ := token.Sign(testId)
					return jwt.Value(), nil
				},
			},
			expectedStatus: 200,
			expectedPayload: map[string]interface{}{
				"token": func() string {
					jwt, _ := token.Sign(testId)
					return jwt.Value()
				}(),
			},
		},
		{
			name: "fail login invalid credentials",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Cpf":"2206144605", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			service: auth.MockService{
				OnAuthenticate: func(ctx context.Context, cpfInput, secret string) (string, error) {
					return "", auth.ErrCredentials
				},
			},
			expectedStatus:  401,
			expectedPayload: map[string]interface{}{"title": responses.ErrCredentials.Payload.Message},
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
			expectedPayload: map[string]interface{}{"title": responses.ErrInvalidLoginPayload.Payload.Message},
		},
		{
			name: "fail login service error",
			httpIO: httpIO{
				r: func() *http.Request {
					return httptest.NewRequest(
						http.MethodPost,
						"/accounts",
						bytes.NewReader([]byte(`{"Cpf":"22061446035", "Secret": "123456"}`)))
				}(),
				w: httptest.NewRecorder(),
			},
			service: auth.MockService{
				OnAuthenticate: func(ctx context.Context, cpfInput, secret string) (string, error) {
					return "", auth.ErrSignToken
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

			h := NewHandler(tt.service)

			h.Login(tt.httpIO.w, tt.httpIO.r)

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
