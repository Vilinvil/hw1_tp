package uniq

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"hw1_tp/internal/uniq"
)

func TestParseArgsSuccessful(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name         string
		input        string
		expectedArgs *uniq.ArgsCommand
	}

	testCases := [...]TestCase{
		{name: "test parse without flags",
			input:        "",
			expectedArgs: &uniq.ArgsCommand{ReadCloser: os.Stdin, WriteCloser: os.Stdout}, //nolint:exhaustivestruct
		},
		{name: "test parse with all flags based on CFlag",
			input: "-c -i -f 10 -s 10",
			expectedArgs: &uniq.ArgsCommand{
				CFlag: true, IFlag: true, FFlag: 10, SFlag: 10, ReadCloser: os.Stdin, WriteCloser: os.Stdout}, //nolint:exhaustivestruct
		},
		{name: "test parse with all flags based on DFlag",
			input: "-d -i -f 10 -s 10",
			expectedArgs: &uniq.ArgsCommand{
				DFlag: true, IFlag: true, FFlag: 10, SFlag: 10, ReadCloser: os.Stdin, WriteCloser: os.Stdout}, //nolint:exhaustivestruct
		},
		{name: "test parse with all flags based on UFlag",
			input: "-u -i -f 10 -s 10",
			expectedArgs: &uniq.ArgsCommand{
				UFlag: true, IFlag: true, FFlag: 10, SFlag: 10, ReadCloser: os.Stdin, WriteCloser: os.Stdout}, //nolint:exhaustivestruct
		}}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			command, err := uniq.ParseArgs(strings.Fields(testCase.input))
			if err != nil {
				t.Errorf("in uniq.ParseArgs(): Error is: %v", err)
			}

			if !reflect.DeepEqual(testCase.expectedArgs, command.Args) {
				t.Errorf("not deepEqual testCase.expectedCommand: %+v and received command: %+v",
					testCase.expectedArgs, command.Args)
			}
		})
	}

}
