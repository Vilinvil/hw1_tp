package calc

import "fmt"

var ErrParenthesis = fmt.Errorf("error with parenthesis")

var ErrNotSupportedInput = fmt.Errorf("get not supported input")

var ErrWrongInput = fmt.Errorf("get incorrect input")

var ErrDivisionZero = fmt.Errorf("can`t division on zero")
