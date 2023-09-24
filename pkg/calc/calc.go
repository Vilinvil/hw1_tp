package calc

import (
	"bufio"
	"fmt"
	"hw1_tp/pkg/stack"
)

func init() {
	input := "233 * 344 / (5 * 666) + 10 "

	fmt.Println(Calc(input, splitTokensInt, isInt, isBasicOperations, getMapPriority()))
}

func handleOperation(operation string, stackOperations *stack.Stack[string], mapPriority mapPriority) (preResult []string) {
	if stackOperations.Len() == 0 || operation == "(" {
		stackOperations.Push(operation)

		return nil
	}

	stackToken := stackOperations.Top()
	tokenPriority := mapPriority[operation]
	stackTokenPriority := mapPriority[stackToken]

	for tokenPriority <= stackTokenPriority {
		preResult = append(preResult, stackToken)

		stackOperations.Pop()

		if stackOperations.Len() != 0 {
			stackToken = stackOperations.Top()
			stackTokenPriority = mapPriority[stackToken]
		} else {
			break
		}
	}

	stackOperations.Push(operation)

	return preResult
}

func handleClosingParenthesis(stackOperations *stack.Stack[string]) (preResult []string, err error) {
	if stackOperations.Len() == 0 {
		return nil, fmt.Errorf("in handleClosingParenthesis(): Error is: %w", ErrParenthesis)
	}

	stackToken := stackOperations.Top()

	for ; ; stackToken = stackOperations.Top() {
		if stackOperations.Len() == 0 {
			return nil, fmt.Errorf("in handleClosingParenthesis(): Error is: %w", ErrParenthesis)
		}

		stackOperations.Pop()

		if stackToken == "(" {
			break
		}

		preResult = append(preResult, stackToken)
	}

	return preResult, nil
}

func convertToRPN(tokens []string, isNumber isNumber, isOperations isOperations,
	mapPriority mapPriority) ([]string, error) {
	var stackOperations stack.Stack[string]

	var result []string

	for _, token := range tokens {
		switch {
		case isNumber(token):
			result = append(result, token)
		case isOperations(token):
			result = append(result, handleOperation(token, &stackOperations, mapPriority)...)
		case token == ")":
			preResult, err := handleClosingParenthesis(&stackOperations)
			if err != nil {
				return nil, fmt.Errorf("in converToPRN(): Error is: %w", err)
			}

			result = append(result, preResult...)
		default:
			return nil, fmt.Errorf("in converToPRN(): Error is: %w. Token is: %s", ErrNotSupportedInput, token)
		}
	}

	for stackOperations.Len() != 0 {
		result = append(result, stackOperations.Top())
		stackOperations.Pop()
	}

	return result, nil
}

func Calc(input string, splitFunc bufio.SplitFunc, isNumber isNumber, isOperations isOperations,
	mapPriority mapPriority) (result Number, err error) {
	prnSl, err := convertToRPN(splitIntoTokens(input, splitFunc), isNumber, isOperations, mapPriority)
	if err != nil {
		return nil, fmt.Errorf("in Calc(): Error is: %w", err)
	}

	stackTokens := stack.Stack[Number]{}

	for _, token := range prnSl {
		switch {
		case isNumber(token):
			var num Number

			num, err := Number.FromString(num, token)
			if err != nil {
				return nil, fmt.Errorf("in Calc(): Error is: %w", err)
			}

			stackTokens.Push(num)
		case isOperations(token):
			num1 := stackTokens.Top()
			stackTokens.Pop()

			num2 := stackTokens.Top()
			stackTokens.Pop()

			switch token {
			case "+":
				stackTokens.Push(num1.Sum(num2))
			case "-":
				stackTokens.Push(num1.Sub(num2))
			case "*":
				stackTokens.Push(num1.Mul(num2))
			case "/":
				stackTokens.Push(num1.Div(num2))
			}
		}
	}

	return stackTokens.Top(), nil
}
