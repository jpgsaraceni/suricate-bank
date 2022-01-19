package token

import "github.com/jpgsaraceni/suricate-bank/app/domain/entities/account"

// receive accountId and return token
// receive token and return accountId

type Jwt struct {
	token string
}

func GenerateJWT(accountId account.AccountId) Jwt {
	return Jwt{token: ""} // TODO: implement jwt
}

func (j Jwt) Value() string {
	return j.token
}
