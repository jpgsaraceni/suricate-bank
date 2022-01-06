package responses

var (
	ErrInternalServerError = ErrorPayload{Message: "internal server error"}
	ErrMissingFields       = ErrorPayload{Message: "missing fields: name, cpf and/or secret"}
	ErrLengthCpf           = ErrorPayload{Message: "invalid cpf length"}
	ErrShortName           = ErrorPayload{Message: "name too short"}
	ErrShortSecret         = ErrorPayload{Message: "password too short"}
	ErrInvalidCpf          = ErrorPayload{Message: "cpf is invalid"}
)
