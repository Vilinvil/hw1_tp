package uniq

import "fmt"

const (
	errTemplate = "%w"

	helpMessage = `uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]

With no options uniq remove duplicate lines and show first occurrence

Options`
	helpNeedMessage = "For help use argument -h or --help"
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
	ErrTooMuchFilesArgs = NewError("too much files arguments. %s", helpNeedMessage)
	ErrTogetherFlagsCDU = NewError("you can`t use c, d, u flags together. %s", helpNeedMessage)
)
