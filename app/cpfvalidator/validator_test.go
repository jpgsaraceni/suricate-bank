package cpfvalidator

import (
	"testing"
)

func TestIsValid(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name     string
		cpf      Cpf
		expected bool
		err      error
	}

	testCases := []testCase{
		{
			name:     "doesn't pass the algorithm",
			cpf:      Cpf("12345678901"),
			expected: false,
		},
		{
			name:     "passes the algorithm but is a known invalid",
			cpf:      Cpf("000.000.000-00"),
			expected: false,
		},
		{
			name:     "passes the algorthm",
			cpf:      Cpf("045.591.180-00"),
			expected: true,
		},
		{
			name:     "invalid format",
			cpf:      Cpf("123.456-789.01"),
			expected: false,
			err:      errFormat,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := tt.cpf.IsValid()

			if tt.err != nil && err == nil {
				t.Errorf("got %t expected an error", got)

				return
			}

			if got != tt.expected {
				t.Errorf("got %t expected %t", got, tt.expected)
			}
		})
	}
}

func TestMask(t *testing.T) {
	//TODO
}
