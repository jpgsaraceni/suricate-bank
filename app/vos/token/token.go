package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"

	"github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"
	"github.com/jpgsaraceni/suricate-bank/config"
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

func Sign(cfg config.Config, accountID account.ID) (Jwt, error) {
	timeout := cfg.JWTTimeout
	if timeout == 0 {
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

	signedToken, err := unsignedToken.SignedString(loadSecret(cfg))
	if err != nil {
		return Jwt{}, fmt.Errorf("%w: %s", ErrSignJWT, err)
	}

	return Jwt{token: signedToken}, nil
}

func Verify(cfg config.Config, tokenString string) (account.ID, error) {
	token, err := jwt.ParseWithClaims(tokenString, &jwtClaimsSchema{}, func(token *jwt.Token) (interface{}, error) {
		return loadSecret(cfg), nil
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

func loadSecret(cfg config.Config) []byte {
	secret := []byte(cfg.JWTSecret)

	return secret
}
