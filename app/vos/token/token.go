package token

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

var (
	ErrSignJWT         = errors.New("failed to sign jwt")
	ErrMissingFieldJWT = errors.New("jwt missing account_id field")
	ErrJwtSignature    = errors.New("invalid token signature")
	ErrParseUuid       = errors.New("failed to parse account id to uuid")
	ErrInvalidClaims   = errors.New("invalid jwt claims")
)

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

func Verify(tokenString string) (account.AccountId, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaimsSchema{}, func(token *jwt.Token) (interface{}, error) {
		return loadSecret(), nil
	})

	if err != nil {

		return account.AccountId{}, ErrJwtSignature
	}

	claims, ok := token.Claims.(*jwtClaimsSchema)

	if !ok || !token.Valid {

		return account.AccountId{}, ErrInvalidClaims
	}

	accountId, err := account.ParseAccountId(claims.AccountId)

	if err != nil {

		return account.AccountId{}, fmt.Errorf("%w: %s", ErrParseUuid, err)
	}

	return accountId, nil
}

func loadSecret() []byte {
	secret := []byte(os.Getenv("JWT_SECRET"))

	return secret
}
