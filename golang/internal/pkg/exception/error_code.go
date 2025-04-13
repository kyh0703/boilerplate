package exception

const (
	CodeServerInternal  = 146000
	CodeValidationError = 146001
	CodeStatusError     = 146002
)

var (
	ErrServerInternal  = NewError(CodeServerInternal, "Internal Server Error", "")
	ErrValidationError = NewError(CodeValidationError, "Validation Error", "")
	ErrStatusError     = NewError(CodeStatusError, "Status Error", "")
)
