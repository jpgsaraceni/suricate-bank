package token

import (
	"errors"
	"log"
	"testing"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

func TestGenerateJWT(t *testing.T) {
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
				t.Errorf("got error %s expected error %s", err, tt.err)
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
