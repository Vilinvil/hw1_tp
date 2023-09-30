package calc

import (
	"fmt"

	"hw1_tp/pkg/stack"
)

func convertToBinaryMinus(input []string) []string {
	var result []string
	result = append(result, getSlFromStr(startInsteadUnaryMinus)...)
	result = append(result, input...)
	result = append(result, ")")

	return result
}

func handleUnaryMinus(input []string) []string {
	var result []string

	idxCur := 0

	if len(input) != 0 && input[0] == "-" {
		innerResult := convertToBinaryMinus(input[:2])
		result = append(append(result, innerResult...), input[2:]...)
		idxCur += len(innerResult)
	} else {
		idxCur++
		result = input
	}

	for ; idxCur < len(result)-1; idxCur++ {
		idxPrev := idxCur - 1
		if result[idxCur] == "-" && (IsBasicOperator(result[idxPrev]) || result[idxPrev] == "(") {
			innerResult := convertToBinaryMinus(result[idxCur : idxCur+2])
			result = append(result[:idxCur], append(innerResult, result[idxCur+2:]...)...)
			idxCur += len(innerResult)
		}
	}

	return result
}

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
	var stackToken string

	var preResult []string

	for {
		if stackOperations.Len() == 0 {
			return nil, fmt.Errorf(errTemplate, ErrParenthesis)
		}

		stackToken = stackOperations.Top()
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
		return nil, fmt.Errorf(errTemplate, ErrWrongInput)
	}

	tokens = handleUnaryMinus(tokens)

	for _, token := range tokens {
		switch {
		case isNumber(token):
			result = append(result, token)
		case isOperator(token) || token == "(":
			result = append(result, handleOperationPRN(token, &stackOperations, mapPriority)...)
		case token == ")":
			preResult, err := handleClosingParenthesisPRN(&stackOperations)
			if err != nil {
				return nil, fmt.Errorf(errTemplate, err)
			}

			result = append(result, preResult...)
		default:
			return nil, fmt.Errorf("token is: %s. %w", token, ErrNotSupportedInput)
		}
	}

	for stackOperations.Len() != 0 {
		result = append(result, stackOperations.Top())
		stackOperations.Pop()
	}

	return result, nil
}
