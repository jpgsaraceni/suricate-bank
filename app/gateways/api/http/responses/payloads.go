package responses

var (
	AccountCreated = Payload{Message: "Account successfully created"}

	// TODO: refatorar payload para ficar junto do erro
	ErrInternalServerError   = Payload{Message: "Internal server error"}
	ErrMissingFields         = Payload{Message: "Missing fields: name, cpf and/or secret"}
	ErrInvalidRequestPayload = Payload{Message: "Invalid payload. Expecting JSON containing keys name, cpf and secret"}
	ErrLengthCpf             = Payload{Message: "CPF must contain 11 numeric digits or 14 digits including 3 symbols"}
	ErrShortName             = Payload{Message: "Name must have at least 3 digits"}
	ErrShortSecret           = Payload{Message: "Secret must have at least 6 digits"}
	ErrInvalidCpf            = Payload{Message: "CPF is invalid"}
	ErrInvalidPathParameter  = Payload{Message: "Invalid URL"}
	ErrAccountNotFound       = Payload{Message: "Account not found"}
)
