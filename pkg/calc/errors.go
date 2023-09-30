package calc

import "fmt"

const (
	errTemplate = "%w"

	helpMessage = `calc "your expression" 

This command supports operators: * / + - unary - ( ) 
Also it supports float numbers.

Example:
 calc 8 * (7) - -2.1
Result: 58.1`

	helpNeedMessage = " for help use -h or --help"
)

type Error struct {
	err string
}

func NewError(format string, args ...any) *Error {
	return &Error{fmt.Sprintf(format, args...)}
}

func (e *Error) Error() string {
	return e.err
}

var (
	ErrParenthesis = NewError("error with parenthesis")

	ErrNotSupportedInput = NewError("get not supported input%s", helpNeedMessage)

	ErrWrongInput = NewError("get incorrect input")

	ErrDivisionZero = NewError("can`t division on zero")

	ErrUnexpectedCountArgs = NewError("unexpected count arguments%s", helpNeedMessage)

	ErrNotEnoughOperands = NewError("get not enough operands for operator.")

	ErrHelpMessage = NewError(helpMessage)
)
