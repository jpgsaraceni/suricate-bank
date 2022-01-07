package money

import (
	"errors"
	"testing"
)

func TestNewMoney(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		amount int
		want   int
		err    error
	}

	testCases := []testCase{
		{
			name:   "create money containing 10",
			amount: 10,
			want:   10,
		},
		{
			name:   "create money containing 0",
			amount: 0,
			want:   0,
		},
		{
			name:   "fail to create money with negative value",
			amount: -1,
			err:    ErrNegative,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newMoney, err := NewMoney(tt.amount)

			if !errors.Is(tt.err, err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			if newMoney.Cents() != tt.want {
				t.Errorf("got %v expected %v", newMoney.Cents(), tt.want)
			}
		})
	}
}

func TestAdd(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		initialValue int
		amount       int
		want         int
		err          error
	}

	testCases := []testCase{
		{
			name:         "add 10 to money containing 0",
			initialValue: 0,
			amount:       10,
			want:         10,
		},
		{
			name:         "fail to add 0 to money containing 0",
			initialValue: 0,
			amount:       0,
			want:         0,
			err:          ErrChangeByZero,
		},
		{
			name:         "fail to add 0 to money containing 10",
			initialValue: 10,
			amount:       0,
			want:         10,
			err:          ErrChangeByZero,
		},
		{
			name:         "fail to add negative to money containing 0",
			initialValue: 0,
			amount:       -10,
			want:         0,
			err:          ErrNegative,
		},
		{
			name:         "fail to add negative to money containing 10",
			initialValue: 10,
			amount:       -10,
			want:         10,
			err:          ErrNegative,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newMoney, _ := NewMoney(tt.initialValue)
			err := newMoney.Add(tt.amount)

			if !errors.Is(tt.err, err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			if newMoney.Cents() != tt.want {
				t.Errorf("got %v expected %v", newMoney.Cents(), tt.want)
			}
		})
	}
}

func TestSubtract(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name         string
		initialValue int
		amount       int
		want         int
		err          error
	}

	testCases := []testCase{
		{
			name:         "subtract 10 from money containing 10",
			initialValue: 10,
			amount:       10,
			want:         0,
		},
		{
			name:         "subtract 10 from money containing 20",
			initialValue: 20,
			amount:       10,
			want:         10,
		},
		{
			name:         "fail to subtract 0 from money containing 0",
			initialValue: 0,
			amount:       0,
			want:         0,
			err:          ErrChangeByZero,
		},
		{
			name:         "fail to subtract 0 from money containing 10",
			initialValue: 10,
			amount:       0,
			want:         10,
			err:          ErrChangeByZero,
		},
		{
			name:         "fail to subtract negative from money containing 0",
			initialValue: 0,
			amount:       -10,
			want:         0,
			err:          ErrNegative,
		},
		{
			name:         "fail to subtract negative from money containing 10",
			initialValue: 10,
			amount:       -10,
			want:         10,
			err:          ErrNegative,
		},
		{
			name:         "fail to subtract 10 from money containing 5",
			initialValue: 5,
			amount:       10,
			want:         5,
			err:          ErrInsuficientFunds,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			newMoney, _ := NewMoney(tt.initialValue)
			err := newMoney.Subtract(tt.amount)

			if !errors.Is(tt.err, err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			if newMoney.Cents() != tt.want {
				t.Errorf("got %v expected %v", newMoney.Cents(), tt.want)
			}
		})
	}
}

func TestScan(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name  string
		value interface{}
		err   error
	}

	testCases := []testCase{
		{
			name:  "successfully scan",
			value: 10,
		},
		{
			name:  "fail to scan empty value",
			value: nil,
			err:   errScanEmpty,
		},
		// {
		// 	name:  "fail to scan invalid cpf",
		// 	value: "220.614.460-34",
		// 	err:   errInvalid,
		// },
		{
			name:  "fail to scan invalid type",
			value: "22061446035",
			err:   errScan,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var testMoney = Money{}

			if err := testMoney.Scan(tt.value); err != tt.err {
				t.Errorf("got error: %s expected error: %s", err, tt.err)
			}
		})
	}
}
