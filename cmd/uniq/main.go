package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
)

const helpMessage = `uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]

With no options uniq remove duplicate lines and show first occurrence

Options`

var ErrHelp = errors.New("for help use argument -h or --help")

const (
	inputFileExist  int8 = 1
	outputFileExist int8 = 2
)

type ArgsUniqCommand struct {
	cFlag bool
	dFlag bool
	uFlag bool
	iFlag bool

	fFlag uint
	sFlag uint

	Reader io.Reader
	Writer io.Writer
}

type UniqCommand struct {
	args               *ArgsUniqCommand
	uniqPrev           bool
	duplicationCounter int
	prevLine           string
}

func (uq *UniqCommand) Writeln(p []byte) (int, error) {
	return uq.args.Writer.Write(append(p, []byte("\n")...))
}

func (uq *UniqCommand) defaultHandlerLine(curLine string) error {
	if uq.prevLine != curLine {
		_, err := uq.Writeln([]byte(uq.prevLine))
		if err != nil {
			return fmt.Errorf("in defaultHandleLine(): Error is: %w", err)
		}

		uq.prevLine = curLine
	}

	return nil
}

func (uq *UniqCommand) handleLastLine(handler func(string) error) error {
	return handler("")
}

func (uq *UniqCommand) Run() error {
	scanner := bufio.NewScanner(uq.args.Reader)
	if scanner.Scan() {
		uq.prevLine = scanner.Text()

		defer func(uq *UniqCommand) {
			err := uq.handleLastLine(uq.defaultHandlerLine)
			if err != nil {
				log.Printf("())Run() defer: Error is: %v", err)
			}
		}(uq)
	}

	for scanner.Scan() {
		err := uq.defaultHandlerLine(scanner.Text())
		if err != nil {
			return fmt.Errorf("in Run(): Error is: %w", err)
		}
	}

	return nil
}

func NewUniqCommand(args *ArgsUniqCommand) *UniqCommand {
	if args.Reader == nil {
		args.Reader = os.Stdin
	}

	if args.Writer == nil {
		args.Writer = os.Stdout
	}

	return &UniqCommand{args: args, duplicationCounter: 1} //nolint:exhaustivestruct,exhaustruct
}

func calcCountTrueFlags(flags ...bool) int {
	countFlags := 0

	for _, value := range flags {
		if value {
			countFlags++
		}
	}

	return countFlags
}

func parseFilesArgs(argsFiles []string) (io.Reader, io.Writer, error) {
	countArgs := len(argsFiles)
	if countArgs > 2 {
		return nil, nil, fmt.Errorf("in parseFilesArgs(): Too much arguments countArgs == %d. %w", countArgs, ErrHelp)
	}

	var err error

	var reader io.Reader
	if countArgs >= int(inputFileExist) {
		reader, err = os.Open(argsFiles[0])
		if err != nil {
			return nil, nil, fmt.Errorf("in parseFilesArgs(): Error is: %w", err)
		}
	}

	var writer io.Writer
	if countArgs == int(outputFileExist) {
		writer, err = os.OpenFile(argsFiles[1], os.O_WRONLY|os.O_CREATE, 0222)
		if err != nil {
			return nil, nil, fmt.Errorf("in parseFilesArgs(): Error is: %w", err)
		}
	}

	return reader, writer, nil
}

func parseArgs() (*UniqCommand, error) {
	flagSet := flag.NewFlagSet(helpMessage, flag.ExitOnError)
	args := ArgsUniqCommand{} //nolint:exhaustruct,exhaustivestruct

	flagSet.BoolVar(&args.cFlag, "c", false, "calculate count duplicated lines")
	flagSet.BoolVar(&args.dFlag, "d", false, "print only duplicated lines")
	flagSet.BoolVar(&args.uFlag, "u", false, "print only unique lines")
	flagSet.BoolVar(&args.iFlag, "i", false, "ignore differences in case when comparing")

	flagSet.UintVar(&args.fFlag, "f", 0, "not compare first `num` fields in every line")
	flagSet.UintVar(&args.sFlag, "s", 0, "not compare first `chars` rune in every line")

	if err := flagSet.Parse(os.Args[1:]); err != nil {
		return nil, fmt.Errorf("in parseArgs(): %w", err)
	}

	countFirstFlags := calcCountTrueFlags(args.cFlag, args.dFlag, args.uFlag)
	if countFirstFlags > 1 {
		return nil, fmt.Errorf("in parseArgs(): you can`t use c, d, u flags together. %w", ErrHelp)
	}

	argsFiles := flagSet.Args()

	var err error

	args.Reader, args.Writer, err = parseFilesArgs(argsFiles)
	if err != nil {
		return nil, fmt.Errorf("in parseArgs(): Error is: %w", err)
	}

	return NewUniqCommand(&args), nil
}

func main() {
	uniqCommand, err := parseArgs()
	if err != nil {
		log.Fatal(fmt.Errorf("error in main(): %w", err))
	}

	err = uniqCommand.Run()
	if err != nil {
		log.Fatal(fmt.Errorf("error in main(): %w", err))
	}
}
