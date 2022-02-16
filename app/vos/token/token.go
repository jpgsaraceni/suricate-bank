package token

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
)

const defaultTimeout = 30

var (
	ErrSignJWT         = errors.New("failed to sign jwt")
	ErrMissingFieldJWT = errors.New("jwt missing account_id field")
	ErrJwtSignature    = errors.New("invalid token signature")
	ErrParseUUID       = errors.New("failed to parse account id to uuid")
	ErrInvalidClaims   = errors.New("invalid jwt claims")
)

type Jwt struct {
	token string
}

func (j Jwt) Value() string {
	return j.token
}

type jwtClaimsSchema struct {
	AccountID string `json:"account_id"`
	jwt.RegisteredClaims
}

func Sign(accountID account.ID) (Jwt, error) {
	timeoutString := os.Getenv("JWT_TIMEOUT")
	timeout, err := strconv.Atoi(timeoutString)
	if err != nil {
		timeout = defaultTimeout
	}

	claims := jwtClaimsSchema{
		accountID.String(),
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * time.Duration(timeout))),
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

func Verify(tokenString string) (account.ID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaimsSchema{}, func(token *jwt.Token) (interface{}, error) {
		return loadSecret(), nil
	})
	if err != nil {
		return account.ID{}, ErrJwtSignature
	}

	claims, ok := token.Claims.(*jwtClaimsSchema)

	if !ok || !token.Valid {
		return account.ID{}, ErrInvalidClaims
	}

	accountID, err := account.ParseAccountID(claims.AccountID)
	if err != nil {
		return account.ID{}, fmt.Errorf("%w: %s", ErrParseUUID, err)
	}

	return accountID, nil
}

func loadSecret() []byte {
	secret := []byte(os.Getenv("JWT_SECRET"))

	return secret
}
