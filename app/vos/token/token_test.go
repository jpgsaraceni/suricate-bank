package token

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func TestSign(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name string
		id   account.AccountId
		want string
		err  error
	}

	var testId = account.AccountId(uuid.New())

	testCases := []testCase{
		{
			name: "successfully generate token",
			id:   testId,
			want: testId.String(),
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			generatedToken, err := Sign(tt.id)

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %s expected error %s", err, tt.err)
			}

			parsedToken, _ := jwt.ParseWithClaims(generatedToken.Value(), &jwtClaimsSchema{}, func(token *jwt.Token) (interface{}, error) {
				return loadSecret(), nil
			})

			claims, ok := parsedToken.Claims.(*jwtClaimsSchema)

			if ok && parsedToken.Valid {
				log.Printf("%v %v", claims.AccountId, claims.RegisteredClaims.ExpiresAt)
			} else {
				t.Fatal("failed to parse token")
			}

			if claims.AccountId != tt.want {
				t.Errorf("got %s expected %s", claims.AccountId, tt.want)
			}
		})
	}
}

func TestVerify(t *testing.T) {
	t.Parallel()

	type testCase struct {
		name  string
		token string
		want  string
		err   error
	}

	var testId = account.AccountId(uuid.New())

	testCases := []testCase{
		{
			name: "successfully verify token",
			token: func() string {
				claims := jwtClaimsSchema{
					testId.String(),
					jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
						Issuer:    "suricate bank",
					},
				}

				unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				signedToken, _ := unsignedToken.SignedString(loadSecret())

				return signedToken
			}(),
			want: testId.String(),
		},
		{
			name: "fail to verify invalid signature",
			token: func() string {
				claims := jwtClaimsSchema{
					testId.String(),
					jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
						Issuer:    "suricate bank",
					},
				}

				unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				signedToken, _ := unsignedToken.SignedString([]byte("wrong signature"))

				return signedToken
			}(),
			want: uuid.Nil.String(),
			err:  ErrJwtSignature,
		},
		{
			name: "fail to verify token missing account_id",
			token: func() string {
				claims := jwtClaimsSchema{
					"",
					jwt.RegisteredClaims{
						ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
						Issuer:    "suricate bank",
					},
				}

				unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

				signedToken, _ := unsignedToken.SignedString(loadSecret())

				return signedToken
			}(),
			want: uuid.Nil.String(),
			err:  ErrParseUuid,
		},
	}

	for _, tt := range testCases {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			j := Jwt{token: tt.token}

			accountId, err := j.Verify()

			if !errors.Is(err, tt.err) {
				t.Fatalf("got error %s expected error %s", err, tt.err)
			}

			if accountId.String() != tt.want {
				t.Errorf("got %s expected %s", accountId.String(), tt.want)
			}
		})
	}
}
