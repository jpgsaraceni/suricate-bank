package schemas

import "github.com/jpgsaraceni/suricate-bank/app/vos/token"

type LoginRequest struct {
	Cpf    string `json:"cpf"`
	Secret string `json:"secret"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginToResponse(token token.Jwt) LoginResponse {
	return LoginResponse{
		Token: token.Value(),
	}
}

// TODO: create login schemas
