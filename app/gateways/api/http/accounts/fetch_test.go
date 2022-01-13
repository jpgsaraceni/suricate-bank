package accountsroute

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
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
		expectedPayload responses.Payload
	}

	var (
		testId1        = uuid.New()
		testId2        = uuid.New()
		testId3        = uuid.New()
		testCpf1       = cpf.Random()
		testCpf2       = cpf.Random()
		testCpf3       = cpf.Random()
		testSecret, _  = hash.NewHash("123")
		testMoney10, _ = money.NewMoney(10)
		testMoney0, _  = money.NewMoney(0)
		testTime       = time.Now()
	)

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
							Id:        account.AccountId(testId1),
							Name:      "acceptable name",
							Cpf:       testCpf1,
							Secret:    testSecret,
							Balance:   testMoney10,
							CreatedAt: testTime,
						},
						{
							Id:        account.AccountId(testId2),
							Name:      "acceptable name",
							Cpf:       testCpf2,
							Secret:    testSecret,
							CreatedAt: testTime,
						},
						{
							Id:        account.AccountId(testId3),
							Name:      "acceptable name",
							Cpf:       testCpf3,
							Secret:    testSecret,
							Balance:   testMoney0,
							CreatedAt: testTime,
						},
					}, nil
				},
			},
			expectedStatus: 200,
			expectedPayload: FetchPayloadTest([]account.Account{
				{
					Id:        account.AccountId(testId1),
					Name:      "acceptable name",
					Cpf:       testCpf1,
					Secret:    testSecret,
					Balance:   testMoney10,
					CreatedAt: testTime,
				},
				{
					Id:        account.AccountId(testId2),
					Name:      "acceptable name",
					Cpf:       testCpf2,
					Secret:    testSecret,
					CreatedAt: testTime,
				},
				{
					Id:        account.AccountId(testId3),
					Name:      "acceptable name",
					Cpf:       testCpf3,
					Secret:    testSecret,
					Balance:   testMoney0,
					CreatedAt: testTime,
				},
			}),
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
					return []account.Account{}, accountuc.ErrNoAccountsToFetch
				},
			},
			expectedStatus:  200,
			expectedPayload: responses.NoAccounts,
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
			expectedPayload: responses.ErrInternalServerError,
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

func FetchPayloadTest(accountList []account.Account) responses.Payload {
	j, _ := json.Marshal(accountList)

	return responses.Ok(responses.FetchedAccountsPayload(j)).Payload
}
