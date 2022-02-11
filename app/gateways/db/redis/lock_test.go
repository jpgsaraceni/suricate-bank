package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis/redistest"
)

func TestLock(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)

	t.Cleanup(tearDown)

	type testCase struct {
		name string
		key  string
		err  error
	}

	testCases := []testCase{
		{},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			err := testRepo.Lock(context.Background(), tt.key)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %s expected error %s", err, tt.err)
			}
		})
	}
}
