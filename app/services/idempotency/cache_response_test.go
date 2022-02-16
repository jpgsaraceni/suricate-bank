package idempotency

import (
	"context"
	"errors"
	"testing"

	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func TestCacheResponse(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name       string
		repository MockRepository
		request    schema.CachedResponse
		err        error
	}

	testKey := uuid.NewString()

	testCases := []testCase{
		{
			name: "cache a response",
			repository: MockRepository{
				OnGetCachedResponse: func(context.Context, string) (schema.CachedResponse, error) {
					return schema.CachedResponse{}, nil
				},
				OnCacheResponse: func(context.Context, schema.CachedResponse) error {
					return nil
				},
			},
		},
		{
			name: "fail to cache a response that is already cached",
			repository: MockRepository{
				OnGetCachedResponse: func(context.Context, string) (schema.CachedResponse, error) {
					return schema.CachedResponse{
						Key:            testKey,
						ResponseStatus: 200,
						ResponseBody:   []byte("awesome marshaled json"),
					}, nil
				},
				OnCacheResponse: func(context.Context, schema.CachedResponse) error {
					return nil
				},
			},
			err: ErrResponseExists,
		},
		{
			name: "fail due to repository error on getting cached response",
			repository: MockRepository{
				OnGetCachedResponse: func(context.Context, string) (schema.CachedResponse, error) {
					return schema.CachedResponse{}, errors.New("i am an error :(")
				},
			},
			err: ErrRepository,
		},
		{
			name: "fail due to repository error on caching response",
			repository: MockRepository{
				OnGetCachedResponse: func(context.Context, string) (schema.CachedResponse, error) {
					return schema.CachedResponse{}, nil
				},
				OnCacheResponse: func(context.Context, schema.CachedResponse) error {
					return errors.New("i am another error :( :(")
				},
			},
			err: ErrRepository,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			s := NewService(tt.repository)

			err := s.CacheResponse(context.Background(), tt.request)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}
		})
	}
}
