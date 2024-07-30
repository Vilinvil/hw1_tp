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
		{
			name:          "test expression with float",
			input:         "2.9 * 2 - 45/2",
			expectedValue: -16.7,
		},
		{
			name:          "test expression with unary minus",
			input:         "-1*( - 2.9 * 2 - (-45)/2)",
			expectedValue: -16.7,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			received, err := calc.Calc(testCase.input, calc.SplitTokensInt, calc.IsFloat, calc.IsBasicOperator,
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
			expectedErr: "get incorrect input",
		},
		{
			name:        "test input with uncorrected rune",
			input:       "2  + 10 + f",
			expectedErr: "token is: f. get not supported input for help use -h or --help",
		},
		{
			name:        "test input with unsupported operations",
			input:       "2 + 10 - 2^3",
			expectedErr: "token is: 2^3. get not supported input for help use -h or --help",
		},
		{
			name:        "test input division on zero",
			input:       "2 + 10 / 0",
			expectedErr: "can`t division on zero",
		},
		{
			name:        "test input with unclosed parenthesis",
			input:       "(2 + 10 * 0",
			expectedErr: "token is: (. Error is: get incorrect input",
		},
		{
			name:        "test input with unopened parenthesis",
			input:       "2 + 10)",
			expectedErr: "error with parenthesis",
		},
		{
			name:        "test input with parenthesis closing more than open parenthesis",
			input:       "(2 + 10 * 0))",
			expectedErr: "error with parenthesis",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := calc.Calc(testCase.input, calc.SplitTokensInt, calc.IsFloat, calc.IsBasicOperator,
				calc.GetMapPriority())
			if err == nil {
				t.Errorf("in calc.Calc() not RETURN ERR\n But EXPECTED: %v", testCase.expectedErr)
			} else if err.Error() != testCase.expectedErr {
				t.Errorf("in calc.Calc() RETURN ERR: %v\n But EXPECTED: %v", err.Error(), testCase.expectedErr)
			}
		})
	}
}
