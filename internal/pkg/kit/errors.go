package kit

type HttpErrorBody struct {
	Code    string `json:"code"`
	Message string `json:"message"`

	// Optionals
	Params map[string]string `json:"params,omitempty"`
}

type HttpError interface {
	HttpCode() int
}

type CodedError interface {
	Code() string
	Error() string
}

type BadRequestError struct {
	Params  map[string]string
	Message string
}

func (e *BadRequestError) HttpCode() int {
	return 400
}

func (e *BadRequestError) Code() string {
	return "bad_request"
}

func (e *BadRequestError) Error() string {
	if e.Message != "" {
		return e.Message
	}

	return "bad request"
}
