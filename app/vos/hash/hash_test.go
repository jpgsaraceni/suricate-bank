package hash

import (
	"errors"
	"testing"

	"golang.org/x/crypto/bcrypt"
)

func TestNewHash(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name   string
		secret string
		err    error
	}

	testCases := []testCase{
		{
			name:   "successfully hash password",
			secret: "123456",
		},
		{
			name:   "successfully hash password with only 1 character",
			secret: "1",
		},
		{
			name:   "successfully hash password with 100 characters",
			secret: "1hgfkljdbngdlgQT$34534621yGhtry5426%$53#$@%6bsgsdgtrtywy$#@%643grfGfdsvgarst%6t23se@745gdfsghT$E#24t",
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			got, err := NewHash(tt.secret)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %v expected error %v", err, tt.err)
			}

			gotHashBytes := []byte(got.Value())
			secretBytes := []byte(tt.secret)

			if err := bcrypt.CompareHashAndPassword(gotHashBytes, secretBytes); err != nil {
				t.Errorf("bcrypt compareHashAndPassword failed with error %s", err)
			}
		})
	}
}

func TestCompare(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name           string
		storedSecret   string
		triedSecret    string
		expectsSuccess bool
		err            error
	}

	testCases := []testCase{
		{
			name:           "successfully compare hashed password",
			storedSecret:   "123456",
			triedSecret:    "123456",
			expectsSuccess: true,
		},
		{
			name:           "fail to compare hashed password",
			storedSecret:   "123456",
			triedSecret:    "123455",
			expectsSuccess: false,
		},
		{
			name:           "successfully compare hashed password with only 1 character",
			storedSecret:   "1",
			triedSecret:    "1",
			expectsSuccess: true,
		},
		{
			name:           "fail to compare hashed password with only 1 character",
			storedSecret:   "1",
			triedSecret:    "0",
			expectsSuccess: false,
		},
		{
			name:           "successfully compare hashed password with 100 characters",
			storedSecret:   "1hgfkljdbngdlgQT$34534621yGhtry5426%$53#$@%6bsgsdgtrtywy$#@%643grfGfdsvgarst%6t23se@745gdfsghT$E#24t",
			triedSecret:    "1hgfkljdbngdlgQT$34534621yGhtry5426%$53#$@%6bsgsdgtrtywy$#@%643grfGfdsvgarst%6t23se@745gdfsghT$E#24t",
			expectsSuccess: true,
		},
		{
			name:           "fail to compare hashed password with 100 characters",
			storedSecret:   "1hgfkljdbngdlgQT$34534621yGhtry5426%$53#$@%6bsgsdgtrtywy$#@%643grfGfdsvgarst%6t23se@745gdfsghT$E#24t",
			triedSecret:    "1hgfkljdbngdlgQT$34534621yGhtry5426%$53#$@%6bgsdgtrtywy$#@%643grfGfdsvgarst%6t23se@745gdfsghT$E#24t",
			expectsSuccess: false,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			hash, err := bcrypt.GenerateFromPassword([]byte(tt.storedSecret), hashCost)

			if err != tt.err {
				t.Errorf("expected error %s got error %s", err, tt.err)

				return
			}

			hashedSecret := Secret{
				hash: string(hash),
			}

			success := hashedSecret.Compare(tt.triedSecret)

			if success != tt.expectsSuccess {
				t.Errorf("got %t expected %t", success, tt.expectsSuccess)
			}
		})
	}
}
