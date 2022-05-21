package redis

import (
	"context"
	"errors"
	"testing"

	"github.com/jpgsaraceni/suricate-bank/app/gateways/db/redis/redistest"
	"github.com/jpgsaraceni/suricate-bank/config"
)

func TestLock(t *testing.T) {
	t.Parallel()

	testConn, tearDown := redistest.GetTestPool()
	testRepo := NewRepository(testConn)

	t.Cleanup(tearDown)

	type testCase struct {
		name          string
		runBefore     func()
		cfg           config.Config
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
			cfg: config.Config{
				RedisConfig: config.RedisConfig{
					IdempotencyKeyTTL: 84600,
				},
			},
			key: repeatedKey,
			err: errKeyExists,
		},
		{
			name: "set a key",
			cfg: config.Config{
				RedisConfig: config.RedisConfig{
					IdempotencyKeyTTL: 84600,
				},
			},
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

			err := testRepo.Lock(context.Background(), tt.cfg, tt.key)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %s expected error %s", err, tt.err)
			}

			if tt.shouldSucceed {
				conn := testRepo.pool.Get()
				reply, err := conn.Do("TTL", tt.key)
				if err != nil {
					t.Fatalf("runBefore failed: %s", err)
				}
				if replyInt, ok := reply.(int64); !ok || replyInt < 84000 {
					t.Fatalf("expect TTL > 84000, got %d", replyInt)
				}
				conn.Close()
			}
		})
	}
}
