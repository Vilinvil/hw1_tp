package calc

import "fmt"

var ErrParenthesis = fmt.Errorf("can`t use ')' before '('")

var ErrNotSupportedInput = fmt.Errorf("get not supported input")
