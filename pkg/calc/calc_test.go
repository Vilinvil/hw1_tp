package calc_test

import (
	"math"
	"testing"

	"hw1_tp/pkg/calc"
)

const tolerance = 1e-8

func equalFloat(num1, num2 float64) bool {
	return math.Abs(num1-num2) < tolerance
}

func TestCalcSuccessful(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name          string
		input         string
		expectedValue float64
	}

	testCases := [...]TestCase{
		{
			name:          "test basic functional",
			input:         "233*344/(5*666)-10",
			expectedValue: 14.06966966,
		},
		{
			name:          "test basic functional with a lot off space rune",
			input:         "233 	* 344 / (5 * 666) -  10",
			expectedValue: 14.06966966,
		},
		{
			name:          "test long expression",
			input:         "(233 * 344 / (5 * 666) - 10) - (((23) * 4 )+ 589)",
			expectedValue: -666.93033033,
		},
		{
			name:          "test zero expression",
			input:         "(233 * 344 / (5 * 666) - 10) * 0",
			expectedValue: 0,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			received, err := calc.Calc(testCase.input, calc.SplitTokensInt, calc.IsInt, calc.IsBasicOperator,
				calc.GetMapPriority())
			if err != nil {
				t.Errorf("in calc.Calc(): unexpexted ERROR: %v", err)
			}

			if !equalFloat(testCase.expectedValue, received) {
				t.Errorf("in calc.Calc(): not equal EXPECTED: %v\n and RECEIVED: %v", testCase.expectedValue, received)
			}
		})
	}
}

//nolint:funlen
func TestCalcErrors(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		input       string
		expectedErr string
	}

	testCases := [...]TestCase{
		{
			name:        "test empty input",
			input:       "",
			expectedErr: "in Calc(): in convertPrn(): Error is: get incorrect input",
		},
		{
			name:        "test input with uncorrected rune",
			input:       "2  + 10 + f",
			expectedErr: "in Calc(): in convertToPRN(): Error is: get not supported input. Token is: f",
		},
		{
			name:        "test input with unsupported operations",
			input:       "2 + 10 - 2^3",
			expectedErr: "in Calc(): in convertToPRN(): Error is: get not supported input. Token is: 2^3",
		},
		{
			name:        "test input with float numbers",
			input:       "2 + 10 - 2.88",
			expectedErr: "in Calc(): in convertToPRN(): Error is: get not supported input. Token is: 2.88",
		},
		{
			name:        "test input division on zero",
			input:       "2 + 10 / 0",
			expectedErr: "in Calc(): in handleToken(): in handleOperator(): Error is: can`t division on zero",
		},
		{
			name:        "test input with unclosed parenthesis",
			input:       "(2 + 10 * 0",
			expectedErr: "in Calc(): in handleToken(): unexpected token. Token is: (. Error is: get incorrect input",
		},
		{
			name:        "test input with parenthesis closing more than open parenthesis",
			input:       "(2 + 10 * 0))",
			expectedErr: "in Calc(): in convertToPRN(): in handleClosingParenthesisPRN(): Error is: error with parenthesis",
		},
		{
			name:  "test input with not enough operands for operator",
			input: " -2 ",
			expectedErr: "in Calc(): in handleToken():" +
				" in handleOperator(): get not enough operands for operator. Error is: get incorrect input",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := calc.Calc(testCase.input, calc.SplitTokensInt, calc.IsInt, calc.IsBasicOperator,
				calc.GetMapPriority())
			if err == nil {
				t.Errorf("in calc.Calc() not RETURN ERR\n But EXPECTED: %v", testCase.expectedErr)
			} else if err.Error() != testCase.expectedErr {
				t.Errorf("in calc.Calc() RETURN ERR: %v\n But EXPECTED: %v", err.Error(), testCase.expectedErr)
			}
		})
	}
}
