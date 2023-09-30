package uniq_test

import (
	"bytes"
	"io"
	"strings"
	"testing"

	"hw1_tp/pkg/uniq"
)

func generateLongString(str string) string {
	for i := 0; i < 10; i++ {
		str += "\n" + str
	}

	return str
}

type nopCloseBuffer struct {
	bytes.Buffer
}

func (nopCloseBuffer) Close() error {
	return nil
}

//nolint:funlen
func TestUniqCommandSuccessful(t *testing.T) {
	t.Parallel()

	type TestCase struct {
		name     string
		input    string
		expected string
		args     uniq.ArgsCommand
	}

	testCases := [...]TestCase{
		{
			name: "Test default command without arguments",
			input: `I love music.
I love music.
I love music.
I love music sooo much.

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
I love music of Kartik.`,
			expected: `I love music.
I love music sooo much.

I love music of Kartik.
Thanks.
I love music of Kartik.
`,
			args: uniq.ArgsCommand{}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name:     "Test default command with empty input",
			input:    ``,
			expected: ``,
			args:     uniq.ArgsCommand{}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name:  "Test default command with a lot of line in input",
			input: generateLongString("test line"),
			expected: `test line
`,
			args: uniq.ArgsCommand{}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name: "Test -c command",
			input: `I love music.
I love music.
I love music.

I love music of Kartik.
I love music of Kartik.
I don't now Katrik.
I don't now Katrik.
I don't now Katrik.
Thanks.
I love music of Kartik.
I love music of Kartik.`,
			expected: `3 I love music.
1 
2 I love music of Kartik.
3 I don't now Katrik.
1 Thanks.
2 I love music of Kartik.
`,
			args: uniq.ArgsCommand{CFlag: true}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name: "Test -d command",
			input: `I love music.
I love go gym

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
I love music.
I love go gym
I love go gym
I love go gym`,
			expected: `I love music of Kartik.
I love music of Kartik.
I love go gym
`,
			args: uniq.ArgsCommand{DFlag: true}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name: "Test -u command",
			input: `I love music.
I love go gym

I love music of Kartik.
I love music of Kartik.
Thanks.
I love music of Kartik.
I love music of Kartik.
I love music.
I love go gym
I love go gym
I love go gym`,
			expected: `I love music.
I love go gym

Thanks.
I love music.
`,
			args: uniq.ArgsCommand{UFlag: true}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name: "Test -i command",
			input: `I lOVe music.
I love go gym

I love music of Kartik.
I love music OF Kartik.
Thanks.
I Love music of Kartik.
I LOVE MUSIC of Kartik.
I love music.
I love go GYM
I LOVE go gym
I love GO gym`,
			expected: `I lOVe music.
I love go gym

I love music of Kartik.
Thanks.
I Love music of Kartik.
I love music.
I love go GYM
`,
			args: uniq.ArgsCommand{IFlag: true}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name: "Test -f command",
			input: `I lOVe music.
I like music.

I love music of Kartik.
I love films of Kartik.
Thanks.
I.
I love 
We love go gym
I LOVE go gym
People love go gym`,
			expected: `I lOVe music.

I love music of Kartik.
I love films of Kartik.
Thanks.
We love go gym
`,
			args: uniq.ArgsCommand{FFlag: 2}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name: "Test -s command",
			input: `I lOVe music.
I like music.

We love music of Kartik.
Me love music of Kartik.
We love go gym
I love go gym
People love go gym`,
			expected: `I lOVe music.
I like music.

We love music of Kartik.
We love go gym
I love go gym
People love go gym
`,
			args: uniq.ArgsCommand{SFlag: 2}, //nolint:exhaustruct,exhaustivestruct
		},
		{
			name: "Test command with all args together based on -c",
			input: `I lOVe music.
I like MUSIC.

We love music of Kartik.
Me love tusic of Kartik.
We love go gym
I love Go GYM
People love go gym`,
			expected: `2 I lOVe music.
1 
2 We love music of Kartik.
3 We love go gym
`,
			args: uniq.ArgsCommand{CFlag: true, IFlag: true, SFlag: 1, FFlag: 2}, //nolint:exhaustruct,exhaustivestruct
		},
	}

	for _, testCase := range testCases {
		testCase := testCase

		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()
			testCase.args.ReadCloser = io.NopCloser(strings.NewReader(testCase.input))
			writerBuf := &nopCloseBuffer{*bytes.NewBufferString("")}
			testCase.args.WriteCloser = writerBuf
			uniqCommand := uniq.NewCommand(&testCase.args)
			defer func(uc *uniq.Command) {
				err := uc.Close()
				if err != nil {
					t.Errorf("in defer: Error is: %v", err)
				}
			}(uniqCommand)

			err := uniqCommand.Run()
			if err != nil {
				t.Fatalf("in uniqCommand.Run(): Error is: %v", err)
			}

			received := writerBuf.String()
			if testCase.expected != received {
				t.Errorf("Not equal EXPECTED: %s \nAnd RECEIVED: %s", testCase.expected, received)
			}
		})
	}
}
