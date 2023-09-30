package uniq_test

import (
	"testing"

	"hw1_tp/pkg/uniq"
)

func TestRunErrors(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		input       []string
		expectedErr string
	}

	testCases := [...]TestCase{
		{
			name:        "test not existing file",
			input:       []string{"not_exist.file"},
			expectedErr: "open not_exist.file: no such file or directory",
		},
		{
			name:        "test unexpexted flag",
			input:       []string{"--flag_unexpected"},
			expectedErr: `flag provided but not defined: -flag_unexpected`,
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			err := uniq.Run(testCase.input)
			if err == nil {
				t.Errorf("NOT RETURN ERR\n But EXPECTED: %v", testCase.expectedErr)
			} else if err.Error() != testCase.expectedErr {
				t.Errorf("RETURN ERR: %v\n But EXPECTED: %v", err.Error(), testCase.expectedErr)
			}
		})
	}
}
