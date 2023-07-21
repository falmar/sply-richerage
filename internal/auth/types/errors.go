package types

type ErrUnauthorized struct {
	Username string
	Message  string
}

func (e *ErrUnauthorized) HttpCode() int {
	return 401
}

func (e *ErrUnauthorized) Code() string {
	return "unauthorized"
}

func (e *ErrUnauthorized) Error() string {
	if e.Username != "" {
		return "unauthorized: " + e.Username
	}
	if e.Message != "" {
		return e.Message
	}

	return "unauthorized"
}

type ErrCredentialsMismatch struct{}

func (e *ErrCredentialsMismatch) HttpCode() int {
	return 401
}

func (e *ErrCredentialsMismatch) Code() string {
	return "credentials_mismatch"
}

func (e *ErrCredentialsMismatch) Error() string {
	return "invalid credentials"
}
