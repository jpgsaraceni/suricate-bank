package cpf

import (
	"errors"
	"testing"
)

func TestIsValid(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		cpf            string
		expectedValue  string
		expectedMasked string
		err            error
	}

	testCases := []testCase{
		{
			name: "doesn't pass the algorithm on second verifying digit",
			cpf:  "12345678901",
			err:  errInvalid,
		},
		{
			name: "doesn't pass the algorithm on first verifying digit",
			cpf:  "12345678910",
			err:  errInvalid,
		},
		{
			name: "passes the algorithm but is a known invalid",
			cpf:  "000.000.000-00",
			err:  errInvalid,
		},
		{
			name: "passes the algorithm but is a known invalid",
			cpf:  "11111111111",
			err:  errInvalid,
		},
		{
			name:           "masked valid cpf",
			cpf:            "220.614.460-35",
			expectedValue:  "22061446035",
			expectedMasked: "220.614.460-35",
		},
		{
			name:           "unmasked valid cpf",
			cpf:            "22061446035",
			expectedValue:  "22061446035",
			expectedMasked: "220.614.460-35",
		},
		{
			name:           "masked valid cpf beginning and ending with 0",
			cpf:            "045.591.180-00",
			expectedValue:  "04559118000",
			expectedMasked: "045.591.180-00",
		},
		{
			name:           "masked valid cpf beginning and ending with 0",
			cpf:            "04559118000",
			expectedValue:  "04559118000",
			expectedMasked: "045.591.180-00",
		},
		{
			name: "valid cpf digits invalid because symbols in wrong place",
			cpf:  "045.591-180.00",
			err:  errInvalid,
		},
		{
			name: "valid cpf digits invalid because missing some symbols",
			cpf:  "045.591.18000",
			err:  errInvalid,
		},
		{
			name: "valid cnpj not valid as cpf",
			cpf:  "34.728.944/0001-00",
			err:  errInvalid,
		},
		{
			name: "valid cnpj without symbols not valid as cpf",
			cpf:  "34728944000100",
			err:  errInvalid,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewCpf(tt.cpf)

			if errors.Is(err, tt.err) {
				t.Fatalf("got error %s expected error %s", got, tt.err)
			}

			if got.Value() != tt.expectedValue {
				t.Errorf("got %s expected %s", got, tt.expectedValue)
			}

			if got.Masked() != tt.expectedMasked {
				t.Errorf("got %s expected %s", got, tt.expectedMasked)
			}
		})
	}
}
