package calc

import (
	"bufio"
	"fmt"
	"strconv"

	"hw1_tp/pkg/stack"
)

func handleOperation(operation string, stackOperations *stack.Stack[string], mapPriority MapPriority) []string {
	if stackOperations.Len() == 0 || operation == "(" {
		stackOperations.Push(operation)

		return nil
	}

	stackToken := stackOperations.Top()
	tokenPriority := mapPriority[operation]
	stackTokenPriority := mapPriority[stackToken]

	var preResult []string

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

func handleClosingParenthesis(stackOperations *stack.Stack[string]) ([]string, error) {
	if stackOperations.Len() == 0 {
		return nil, fmt.Errorf("in handleClosingParenthesis(): Error is: %w", ErrParenthesis)
	}

	stackToken := stackOperations.Top()

	var preResult []string

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

func convertToRPN(tokens []string, isNumber IsNumber, isOperations IsOperations,
	mapPriority MapPriority) ([]string, error,
) {
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

func Calc(input string, splitFunc bufio.SplitFunc, isNumber IsNumber, isOperations IsOperations,
	mapPriority MapPriority,
) (float64, error) {
	prnSl, err := convertToRPN(splitIntoTokens(input, splitFunc), isNumber, isOperations, mapPriority)
	if err != nil {
		return 0, fmt.Errorf("in Calc(): Error is: %w", err)
	}

	stackTokens := stack.Stack[float64]{}

	for _, token := range prnSl {
		switch {
		case isNumber(token):
			num, err := strconv.ParseFloat(token, 64)
			if err != nil {
				return 0, fmt.Errorf("in Calc(): Error is: %w", err)
			}

			stackTokens.Push(num)
		case isOperations(token):
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
				stackTokens.Push(num2 / num1)
			}
		}
	}

	return stackTokens.Top(), nil
}
