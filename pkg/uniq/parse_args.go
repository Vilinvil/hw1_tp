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
		return nil, nil, fmt.Errorf(errTemplate, NewError("count args = %d. %s", countArgs, ErrTooMuchFilesArgs.Error()))
	}

	var err error

	var readCloser io.ReadCloser
	if countArgs >= int(inputFileExist) {
		readCloser, err = os.Open(argsFiles[0])
		if err != nil {
			return nil, nil, fmt.Errorf(errTemplate, err)
		}
	}

	var writeCloser io.WriteCloser
	if countArgs == int(outputFileExist) {
		writeCloser, err = os.OpenFile(argsFiles[1], os.O_WRONLY|os.O_CREATE, readeWriteEnable) //nolint:nosnakecase
		if err != nil {
			return nil, nil, fmt.Errorf(errTemplate, err)
		}
	}

	return readCloser, writeCloser, nil
}

func ParseArgs(args []string) (*Command, error) {
	flagSet := flag.NewFlagSet(helpMessage, flag.ContinueOnError)
	argsCommand := ArgsCommand{} //nolint:exhaustruct,exhaustivestruct

	flagSet.BoolVar(&argsCommand.CFlag, "c", false, "calculate count duplicated lines")
	flagSet.BoolVar(&argsCommand.DFlag, "d", false, "print only duplicated lines")
	flagSet.BoolVar(&argsCommand.UFlag, "u", false, "print only unique lines")
	flagSet.BoolVar(&argsCommand.IFlag, "i", false, "ignore differences in case when comparing")

	flagSet.UintVar(&argsCommand.FFlag, "f", 0, "not compare first `num` fields in every line")
	flagSet.UintVar(&argsCommand.SFlag, "s", 0, "not compare first `chars` rune in every line")

	if err := flagSet.Parse(args); err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}

	countFirstFlags := CalcCountTrueFlags(argsCommand.CFlag, argsCommand.DFlag, argsCommand.UFlag)
	if countFirstFlags > 1 {
		return nil, fmt.Errorf(errTemplate, ErrTogetherFlagsCDU)
	}

	argsFiles := flagSet.Args()

	var err error

	argsCommand.ReadCloser, argsCommand.WriteCloser, err = parseFilesArgs(argsFiles)
	if err != nil {
		return nil, fmt.Errorf(errTemplate, err)
	}

	return NewCommand(&argsCommand), nil
}
