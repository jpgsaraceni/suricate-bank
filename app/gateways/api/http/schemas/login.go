package schemas

type LoginRequest struct {
	Cpf    string `json:"cpf"`
	Secret string `json:"secret"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginToResponse(token string) LoginResponse {
	return LoginResponse{
		Token: token,
	}
}
