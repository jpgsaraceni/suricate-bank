package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/joho/godotenv"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

var ErrSignJWT = errors.New("failed to sign jwt")

type Jwt struct {
	token string
}

func (j Jwt) Value() string {
	return j.token
}

type jwtClaimsSchema struct {
	AccountId string `json:"account_id"`
	jwt.RegisteredClaims
}

func Sign(accountId account.AccountId) (Jwt, error) {

	claims := jwtClaimsSchema{
		accountId.String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 30)),
			Issuer:    "suricate bank",
		},
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := unsignedToken.SignedString(loadSecret())

	if err != nil {

		return Jwt{}, fmt.Errorf("%w: %s", ErrSignJWT, err)
	}

	return Jwt{token: signedToken}, nil
}

func loadSecret() []byte {
	godotenv.Load()
	return []byte(os.Getenv("JWT_SECRET"))
}
