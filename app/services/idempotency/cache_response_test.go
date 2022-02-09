package idempotency

import (
	"errors"
	"testing"

	"github.com/google/uuid"
	"github.com/jpgsaraceni/suricate-bank/app/services/idempotency/schema"
)

func TestCacheResponse(t *testing.T) {
	type testCase struct {
		name       string
		key        string
		repository MockRepository
		request    schema.CachedResponse
		err        error
	}

	testKey := uuid.NewString()

	testCases := []testCase{
		{
			name: "cache a response",
			key:  testKey,
			repository: MockRepository{
				OnGetCachedResponse: func(key string) (schema.CachedResponse, error) {
					return schema.CachedResponse{}, nil
				},
				OnCacheResponse: func(key string, request schema.CachedResponse) error {
					return nil
				},
			},
		},
		{
			name: "fail to cache a response that is already cached",
			key:  testKey,
			repository: MockRepository{
				OnGetCachedResponse: func(key string) (schema.CachedResponse, error) {
					return schema.CachedResponse{
						Key:            testKey,
						ResponseStatus: 200,
						ResponseBody:   []byte("awesome marshaled json"),
					}, nil
				},
				OnCacheResponse: func(key string, request schema.CachedResponse) error {
					return nil
				},
			},
			err: ErrResponseExists,
		},
		{
			name: "fail due to repository error on getting cached response",
			key:  testKey,
			repository: MockRepository{
				OnGetCachedResponse: func(key string) (schema.CachedResponse, error) {
					return schema.CachedResponse{}, errors.New("i am an error :(")
				},
			},
			err: ErrRepository,
		},
		{
			name: "fail due to repository error on caching response",
			key:  testKey,
			repository: MockRepository{
				OnGetCachedResponse: func(key string) (schema.CachedResponse, error) {
					return schema.CachedResponse{}, nil
				},
				OnCacheResponse: func(key string, request schema.CachedResponse) error {
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

			err := s.CacheResponse(tt.key, tt.request)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got %s expected %s", err, tt.err)
			}
		})
	}
}
