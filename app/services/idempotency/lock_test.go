package idempotency

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/config"
)

func TestLock(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		key  string
		repo MockRepository
		err  error
	}

	testCases := []testCase{
		{
			name: "create key in repository",
			key:  uuid.NewString(),
			repo: MockRepository{
				OnLock: func(_ context.Context, _ config.Config, _ string) error {
					return nil
				},
			},
		},
		{
			name: "repository error",
			key:  uuid.NewString(),
			repo: MockRepository{
				OnLock: func(_ context.Context, _ config.Config, _ string) error {
					return errors.New("bad error")
				},
			},
			err: ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
		})

		s := NewService(tt.repo)

		err := s.Lock(context.Background(), config.Config{}, tt.key)
		if !errors.Is(err, tt.err) {
			t.Fatalf("got error %s expected error %s", err, tt.err)
		}
	}
}
