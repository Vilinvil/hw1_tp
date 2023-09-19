package uniq

import (
	"flag"
	"fmt"
	"io"
	"os"
)

func parseFilesArgs(argsFiles []string) (io.ReadCloser, io.WriteCloser, error) {
	countArgs := len(argsFiles)
	if countArgs > int(outputFileExist) {
		return nil, nil, fmt.Errorf("in parseFilesArgs(): Too much arguments countArgs == %d. %w", countArgs, ErrHelp)
	}

	var err error

	var readCloser io.ReadCloser
	if countArgs >= int(inputFileExist) {
		readCloser, err = os.Open(argsFiles[0])
		if err != nil {
			return nil, nil, fmt.Errorf("in parseFilesArgs(): Error is: %w", err)
		}
	}

	var writeCloser io.WriteCloser
	if countArgs == int(outputFileExist) {
		writeCloser, err = os.OpenFile(argsFiles[1], os.O_WRONLY|os.O_CREATE, 0666)
		if err != nil {
			return nil, nil, fmt.Errorf("in parseFilesArgs(): Error is: %w", err)
		}
	}

	return readCloser, writeCloser, nil
}

func ParseArgs() (*Command, error) {
	flagSet := flag.NewFlagSet(helpMessage, flag.ExitOnError)
	args := ArgsCommand{} //nolint:exhaustruct,exhaustivestruct

	flagSet.BoolVar(&args.CFlag, "c", false, "calculate count duplicated lines")
	flagSet.BoolVar(&args.DFlag, "d", false, "print only duplicated lines")
	flagSet.BoolVar(&args.UFlag, "u", false, "print only unique lines")
	flagSet.BoolVar(&args.IFlag, "i", false, "ignore differences in case when comparing")

	flagSet.UintVar(&args.FFlag, "f", 0, "not compare first `num` fields in every line")
	flagSet.UintVar(&args.SFlag, "s", 0, "not compare first `chars` rune in every line")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("in ParseArgs(): %w", err)
	}

	countFirstFlags := calcCountTrueFlags(args.CFlag, args.DFlag, args.UFlag)
	if countFirstFlags > 1 {
		return nil, fmt.Errorf("in ParseArgs(): you can`t use c, d, u flags together. %w", ErrHelp)
	}

	argsFiles := flagSet.Args()

	var err error

	args.ReadCloser, args.WriteCloser, err = parseFilesArgs(argsFiles)
	if err != nil {
		return nil, fmt.Errorf("in ParseArgs(): Error is: %w", err)
	}

	return NewCommand(&args), nil
}
