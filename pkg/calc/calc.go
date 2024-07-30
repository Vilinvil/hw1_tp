package calc

import (
	"bufio"
	"fmt"
	"strconv"

	"hw1_tp/pkg/stack"
)

func handleOperator(token string, stackTokens *stack.Stack[float64]) error {
	if stackTokens.Len() < countOperandsForOperator {
		return fmt.Errorf(errTemplate, ErrNotEnoughOperands)
	}

	num1 := stackTokens.Top()
	stackTokens.Pop()

	num2 := stackTokens.Top()
	stackTokens.Pop()

	switch token {
	case "+":
		stackTokens.Push(num1 + num2)
	case "-":
		stackTokens.Push(num2 - num1)
	case "*":
		stackTokens.Push(num1 * num2)
	case "/":
		if num1 == 0 {
			return fmt.Errorf(errTemplate, ErrDivisionZero)
		}

		stackTokens.Push(num2 / num1)
	}

	return nil
}

func handleToken(token string, stackTokens *stack.Stack[float64], isNumber IsNumber, isOperator IsOperator) error {
	switch {
	case isNumber(token):
		num, err := strconv.ParseFloat(token, 64)
		if err != nil {
			return fmt.Errorf(errTemplate, err)
		}

		stackTokens.Push(num)
	case isOperator(token):
		err := handleOperator(token, stackTokens)
		if err != nil {
			return fmt.Errorf(errTemplate, err)
		}
	default:
		return fmt.Errorf("token is: %s. Error is: %w", token, ErrWrongInput)
	}

	return nil
}

func Calc(input string, splitFunc bufio.SplitFunc, isNumber IsNumber, isOperator IsOperator,
	mapPriority MapPriority,
) (float64, error) {
	prnSl, err := convertToRPN(splitIntoTokens(input, splitFunc), isNumber, isOperator, mapPriority)
	if err != nil {
		return 0, fmt.Errorf(errTemplate, err)
	}

	stackTokens := stack.Stack[float64]{}

	for _, token := range prnSl {
		err = handleToken(token, &stackTokens, isNumber, isOperator)
		if err != nil {
			return 0, fmt.Errorf(errTemplate, err)
		}
	}

	return stackTokens.Top(), nil
}
