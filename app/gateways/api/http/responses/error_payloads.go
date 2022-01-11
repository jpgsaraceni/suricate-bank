package responses

var (
	ErrInternalServerError = "Internal server error"
	ErrMissingFields       = "Missing fields: name, cpf and/or secret"
	ErrLengthCpf           = "CPF must contain 11 numeric digits or 14 digits including 3 symbols"
	ErrShortName           = "Name must have at least 3 digits"
	ErrShortSecret         = "Secret must have at least 6 digits"
	ErrInvalidCpf          = "cpf is invalid"
)
