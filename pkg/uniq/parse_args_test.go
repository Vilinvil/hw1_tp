package uniq_test

import (
	"os"
	"reflect"
	"strings"
	"testing"

	"hw1_tp/pkg/uniq"
)

func TestParseArgsSuccessful(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name         string
		input        string
		expectedArgs *uniq.ArgsCommand
	}

	testCases := [...]TestCase{
		{
			name:         "test parse without flags",
			input:        "",
			expectedArgs: &uniq.ArgsCommand{ReadCloser: os.Stdin, WriteCloser: os.Stdout}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name:  "test parse with all flags based on CFlag",
			input: "-c -i -f 10 -s 10",
			expectedArgs: &uniq.ArgsCommand{ //nolint:exhaustruct,exhaustivestruct
				CFlag: true, IFlag: true, FFlag: 10, SFlag: 10,
				ReadCloser: os.Stdin, WriteCloser: os.Stdout,
			},
		},
		{
			name:  "test parse with all flags based on DFlag",
			input: "-d -i -f 10 -s 10",
			expectedArgs: &uniq.ArgsCommand{ //nolint:exhaustruct,exhaustivestruct
				DFlag: true, IFlag: true, FFlag: 10, SFlag: 10,
				ReadCloser: os.Stdin, WriteCloser: os.Stdout,
			},
		},
		{
			name:  "test parse with all flags based on UFlag",
			input: "-u -i -f 10 -s 10",
			expectedArgs: &uniq.ArgsCommand{ //nolint:exhaustruct,exhaustivestruct
				UFlag: true, IFlag: true, FFlag: 10, SFlag: 10,
				ReadCloser: os.Stdin, WriteCloser: os.Stdout,
			},
		},
	}

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

func TestParseArgsErrors(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name        string
		input       string
		expectedErr string
	}

	testCases := [...]TestCase{
		{
			name:        "test parse with unnecessary args of files ",
			input:       "input.file output.file unnecessary.file",
			expectedErr: "count args = 3. too much files arguments. For help use argument -h or --help",
		},
		{
			name:        "test parse not existing input file",
			input:       "not_exist.file",
			expectedErr: "open not_exist.file: no such file or directory",
		},
		{
			name:        "test parse with uncorrected flags",
			input:       "--uncorrected",
			expectedErr: "flag provided but not defined: -uncorrected",
		},
		{
			name:        "test parse with CFlag and DFlag together ",
			input:       "-c -d",
			expectedErr: "you can`t use c, d, u flags together. For help use argument -h or --help",
		},
		{
			name:        "test parse with UFlag and DFlag together ",
			input:       "-u -d",
			expectedErr: "you can`t use c, d, u flags together. For help use argument -h or --help",
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			_, err := uniq.ParseArgs(strings.Fields(testCase.input))
			if err == nil {
				t.Errorf("uniq.ParseArgs() not RETURN ERR\n But EXPECTED: %v", testCase.expectedErr)
			} else if err.Error() != testCase.expectedErr {
				t.Errorf("uniq.ParseArgs() RETURN ERR: %v\n But EXPECTED: %v", err.Error(), testCase.expectedErr)
			}
		})
	}
}
