package calc_test

import (
	"testing"

	"hw1_tp/pkg/calc"
)

func TestRunSuccessful(t *testing.T) {
	t.Parallel()

	args := []string{"3 - 3 + 0"}

	received, err := calc.Run(args)
	if err != nil {
		t.Errorf("unexpexted ERROR: %v", err)
	}

	if expected := 0.0; !equalFloat(expected, received) {
		t.Errorf("not equal EXPECTED: %v\n and RECEIVED: %v", expected, received)
	}
}

func TestRunErrors(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		input       []string
		expectedErr string
	}

	testCases := [...]TestCase{
		{
			name:        "test empty input",
			input:       []string{""},
			expectedErr: "get incorrect input",
		},
		{
			name:        "test a lot of args",
			input:       []string{"", ""},
			expectedErr: "count args == 2. unexpected count arguments for help use -h or --help",
		},
		{
			name:  "test help",
			input: []string{"-h"},
			expectedErr: `calc "your expression" 

This command supports operators: * / + - unary - ( ) 
Also it supports float numbers.

Example:
 calc 8 * (7) - -2.1
Result: 58.1`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := calc.Run(testCase.input)
			if err == nil {
				t.Errorf("NOT RETURN ERR\n But EXPECTED: %v", testCase.expectedErr)
			} else if err.Error() != testCase.expectedErr {
				t.Errorf("RETURN ERR: %v\n But EXPECTED: %v", err.Error(), testCase.expectedErr)
			}
		})
	}
}
