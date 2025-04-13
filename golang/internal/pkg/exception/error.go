package exception

type Error struct {
	Code    int
	Message string
	Detail  string
}

func NewError(code int, message string, detail string) *Error {
	return &Error{
		Code:    code,
		Message: message,
		Detail:  detail,
	}
}

func (e *Error) Error() string {
	return e.Message
}
