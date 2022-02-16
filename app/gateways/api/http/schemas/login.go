package schemas

type LoginRequest struct {
	Cpf    string `json:"cpf" example:"04559118000"`
	Secret string `json:"secret" example:"great-password"`
}

type LoginResponse struct {
	Token string `json:"token"`
}

func LoginToResponse(token string) LoginResponse {
	return LoginResponse{
		Token: token,
	}
}
