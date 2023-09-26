package calc

import (
	"bufio"
	"fmt"
	"strconv"

	"hw1_tp/pkg/stack"
)

func handleOperationPRN(operation string, stackOperations *stack.Stack[string], mapPriority MapPriority) []string {
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

func handleClosingParenthesisPRN(stackOperations *stack.Stack[string]) ([]string, error) {
	if stackOperations.Len() == 0 {
		return nil, fmt.Errorf("in handleClosingParenthesisPRN(): Error is: %w", ErrParenthesis)
	}

	stackToken := stackOperations.Top()

	var preResult []string

	for ; ; stackToken = stackOperations.Top() {
		if stackOperations.Len() == 0 {
			return nil, fmt.Errorf("in handleClosingParenthesisPRN(): Error is: %w", ErrParenthesis)
		}

		stackOperations.Pop()

		if stackToken == "(" {
			break
		}

		preResult = append(preResult, stackToken)
	}

	return preResult, nil
}

func convertToRPN(tokens []string, isNumber IsNumber, isOperator IsOperator,
	mapPriority MapPriority) ([]string, error,
) {
	var stackOperations stack.Stack[string]

	var result []string

	if len(tokens) == 0 {
		return nil, fmt.Errorf("in convertPrn(): Error is: %w", ErrWrongInput)
	}

	for _, token := range tokens {
		switch {
		case isNumber(token):
			result = append(result, token)
		case isOperator(token) || token == "(":
			result = append(result, handleOperationPRN(token, &stackOperations, mapPriority)...)
		case token == ")":
			preResult, err := handleClosingParenthesisPRN(&stackOperations)
			if err != nil {
				return nil, fmt.Errorf("in convertToPRN(): %w", err)
			}

			result = append(result, preResult...)
		default:
			return nil, fmt.Errorf("in convertToPRN(): Error is: %w. Token is: %s", ErrNotSupportedInput, token)
		}
	}

	for stackOperations.Len() != 0 {
		result = append(result, stackOperations.Top())
		stackOperations.Pop()
	}

	return result, nil
}

func handleOperator(token string, stackTokens *stack.Stack[float64]) error {
	if stackTokens.Len() < countOperandsForOperator {
		return fmt.Errorf("in handleOperator(): get not enough operands for operator. Error is: %w", ErrWrongInput)
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
			return fmt.Errorf("in handleOperator(): Error is: %w", ErrDivisionZero)
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
			return fmt.Errorf("in handleToken(): Error is: %w", err)
		}

		stackTokens.Push(num)
	case isOperator(token):
		err := handleOperator(token, stackTokens)
		if err != nil {
			return fmt.Errorf("in handleToken(): %w", err)
		}
	default:
		return fmt.Errorf("in handleToken(): unexpected token. Token is: %s. Error is: %w", token, ErrWrongInput)
	}

	return nil
}

func Calc(input string, splitFunc bufio.SplitFunc, isNumber IsNumber, isOperator IsOperator,
	mapPriority MapPriority,
) (float64, error) {
	prnSl, err := convertToRPN(splitIntoTokens(input, splitFunc), isNumber, isOperator, mapPriority)
	if err != nil {
		return 0, fmt.Errorf("in Calc(): %w", err)
	}

	stackTokens := stack.Stack[float64]{}

	for _, token := range prnSl {
		err = handleToken(token, &stackTokens, isNumber, isOperator)
		if err != nil {
			return 0, fmt.Errorf("in Calc(): %w", err)
		}
	}

	return stackTokens.Top(), nil
}
