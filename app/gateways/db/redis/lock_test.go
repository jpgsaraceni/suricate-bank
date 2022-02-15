package redis

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis/redistest"
)

func TestLock(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)
	os.Setenv("IDEMPOTENCY_TTL", "24")

	t.Cleanup(tearDown)

	type testCase struct {
		name          string
		runBefore     func()
		key           string
		shouldSucceed bool
		err           error
	}

	repeatedKey := "nicekey"

	testCases := []testCase{
		{
			name: "fail to set a key that already exists in redis",
			runBefore: func() {
				conn := testRepo.pool.Get()
				_, err := conn.Do("SET", repeatedKey, "")
				if err != nil {
					t.Fatalf("runBefore failed: %s", err)
				}
				conn.Close()
			},
			key: repeatedKey,
			err: errKeyExists,
		},
		{
			name:          "set a key",
			key:           "greatkey",
			shouldSucceed: true,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			if tt.runBefore != nil {
				tt.runBefore()
			}

			err := testRepo.Lock(context.Background(), tt.key)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %s expected error %s", err, tt.err)
			}

			if tt.shouldSucceed {
				conn := testRepo.pool.Get()
				reply, err := conn.Do("TTL", tt.key)
				if err != nil {
					t.Fatalf("runBefore failed: %s", err)
				}
				if replyInt := reply.(int64); replyInt < 23*60*60 {
					t.Fatalf("expect TTL > 23 hours, got %d", replyInt)
				}
				conn.Close()
			}
		})
	}
}
