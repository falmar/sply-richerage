package hasher

type ErrInvalidToken struct {
	Message string
}

func (e *ErrInvalidToken) HttpCode() int {
	return 401
}
func (e *ErrInvalidToken) Code() string {
	return "invalid_token"
}
func (e *ErrInvalidToken) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "invalid token"
}

type ErrExpiredToken struct{}

func (e *ErrExpiredToken) HttpCode() int {
	return 401
}
func (e *ErrExpiredToken) Code() string {
	return "expired_token"
}
func (e *ErrExpiredToken) Error() string {
	return "token expired"
}

type ErrEmptyData struct{}

func (e *ErrEmptyData) HttpCode() int {
	return 500
}

func (e *ErrEmptyData) Code() string {
	return "empty_data"
}

func (e *ErrEmptyData) Error() string {
	return "empty data"
}
