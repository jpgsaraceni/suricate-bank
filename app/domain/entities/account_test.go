package entities

import "testing"

func TestNewAccount(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name        string
		cpf         Cpf
		secret      string
		expectedErr error
	}

	testCases := []testCase{
		{},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewAccount(tt.name, tt.cpf, tt.secret)

			if tt.expectedErr != err {
				t.Errorf("got error %v espected error %v", err, tt.expectedErr)
			}
		})
	}
}
