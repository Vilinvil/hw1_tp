package calc

import (
	"testing"

	"hw1_tp/pkg/calc"
)

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

			received, err := calc.Calc(testCase.input, calc.SplitTokensInt, calc.IsInt, calc.IsBasicOperations,
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
			expectedErr: "in Calc(): Error is: in convertPrn(): get incorrect input without tokens",
		},
		{
			name:        "test input with uncorrected rune",
			input:       "2  + 10 + f",
			expectedErr: "in Calc(): Error is: in converToPRN(): Error is: get not supported input. Token is: f",
		},
		{
			name:        "test input with unsupported operations",
			input:       "2 + 10 - 2^3",
			expectedErr: "in Calc(): Error is: in converToPRN(): Error is: get not supported input. Token is: 2^3",
		},
		{
			name:        "test input with float numbers",
			input:       "2 + 10 - 2.88",
			expectedErr: "in Calc(): Error is: in converToPRN(): Error is: get not supported input. Token is: 2.88",
		},
		{
			name:        "test input division on zero",
			input:       "2 + 10 / 0",
			expectedErr: "in Calc(): Error is: can`t division on zero",
		},
		{
			name:        "test input with unclosed parenthesis",
			input:       "(2 + 10 * 0",
			expectedErr: "in Calc(): unexpected '('. Error is: error with parenthesis",
		},
		{
			name:  "test input with parenthesis closing more than open parenthesis",
			input: "(2 + 10 * 0))",
			expectedErr: "in Calc(): Error is: in converToPRN():" +
				" Error is: in handleClosingParenthesis(): Error is: error with parenthesis",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := calc.Calc(testCase.input, calc.SplitTokensInt, calc.IsInt, calc.IsBasicOperations,
				calc.GetMapPriority())
			if err == nil {
				t.Errorf("in calc.Calc() not RETURN ERR\n But EXPECTED: %v", testCase.expectedErr)
			} else if err.Error() != testCase.expectedErr {
				t.Errorf("in calc.Calc() RETURN ERR: %v\n But EXPECTED: %v", err.Error(), testCase.expectedErr)
			}
		})
	}
}
