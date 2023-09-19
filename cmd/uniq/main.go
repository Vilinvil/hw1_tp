package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

const helpMessage = `uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]

With no options uniq remove duplicate lines and show first occurrence

Options`

var ErrHelp = errors.New("for help use argument -h or --help")

const (
	inputFileExist  int8 = 1
	outputFileExist int8 = 2
)

func calcCountTrueFlags(flags ...bool) int {
	countFlags := 0

	for _, value := range flags {
		if value {
			countFlags++
		}
	}

	return countFlags
}

func TrimFirstFields(str string, count uint) string {
	idx := strings.Index(str, " ")
	for ; count > 0; count-- {
		if idx == -1 {
			return ""
		}

		str = str[idx+1:]
		idx = strings.Index(str, " ")
	}

	return str
}

type comparatorString func(first, second string) bool

type ArgsUniqCommand struct {
	cFlag bool
	dFlag bool
	uFlag bool
	iFlag bool

	fFlag uint
	sFlag uint

	readCloser  io.ReadCloser
	writeCloser io.WriteCloser
}

type UniqCommand struct {
	args               *ArgsUniqCommand
	uniqPrev           bool
	duplicationCounter int
	prevLine           string
	compStr            comparatorString
}

func NewUniqCommand(args *ArgsUniqCommand) *UniqCommand {
	if args.readCloser == nil {
		args.readCloser = os.Stdin
	}

	if args.writeCloser == nil {
		args.writeCloser = os.Stdout
	}

	return &UniqCommand{args: args, duplicationCounter: 1, uniqPrev: true} //nolint:exhaustivestruct,exhaustruct
}

func (uq *UniqCommand) Close() error {
	err := uq.args.readCloser.Close()
	if err != nil {
		return fmt.Errorf("in Close(): Error is: %w", err)
	}

	err = uq.args.writeCloser.Close()
	if err != nil {
		return fmt.Errorf("in Close(): Error is: %w", err)
	}

	return nil
}

func (uq *UniqCommand) Writeln(p []byte) (int, error) {
	return uq.args.writeCloser.Write(append(p, []byte("\n")...)) //nolint:wrapcheck
}

func (uq *UniqCommand) fieldsHandleLine(curLine string) string {
	if uq.args.fFlag != 0 {
		curLine = TrimFirstFields(curLine, uq.args.fFlag)
	}

	return curLine
}

func (uq *UniqCommand) skipCharsHandleLine(curLine string) string {
	if uq.args.sFlag != 0 {
		if uint(len(curLine)) >= uq.args.sFlag {
			curLine = curLine[uq.args.sFlag:]
		} else {
			curLine = ""
		}
	}

	return curLine
}

func (uq *UniqCommand) ignoreCaseHandleLine(curLine string) string {
	if uq.args.iFlag {
		curLine = strings.ToLower(curLine)
	}

	return curLine
}

func (uq *UniqCommand) initComparator() {
	uq.compStr = func(first, second string) bool {
		first = uq.ignoreCaseHandleLine(uq.skipCharsHandleLine(uq.fieldsHandleLine(first)))
		second = uq.ignoreCaseHandleLine(uq.skipCharsHandleLine(uq.fieldsHandleLine(second)))

		return first == second
	}
}

func (uq *UniqCommand) defaultHandleLine(curLine string) error {
	if !uq.compStr(uq.prevLine, curLine) {
		_, err := uq.Writeln([]byte(uq.prevLine))
		if err != nil {
			return fmt.Errorf("in defaultHandleLine(): Error is: %w", err)
		}

		uq.prevLine = curLine
	}

	return nil
}

func (uq *UniqCommand) calcHandleLine(curLine string) error {
	if uq.compStr(uq.prevLine, curLine) {
		uq.duplicationCounter++

		return nil
	}

	_, err := uq.Writeln([]byte(strconv.Itoa(uq.duplicationCounter) + " " + uq.prevLine))
	if err != nil {
		return fmt.Errorf("in calcHandlerLine(): Error is: %w", err)
	}

	uq.prevLine = curLine
	uq.duplicationCounter = 1

	return nil
}

func (uq *UniqCommand) uniqHandleLine(curLine string) error {
	if uq.compStr(uq.prevLine, curLine) {
		uq.uniqPrev = false

		return nil
	}

	if uq.uniqPrev {
		_, err := uq.Writeln([]byte(uq.prevLine))
		if err != nil {
			return fmt.Errorf("in calcHandlerLine(): Error is: %w", err)
		}
	}

	uq.prevLine = curLine
	uq.uniqPrev = true

	return nil
}

func (uq *UniqCommand) duplicateHandleLine(curLine string) error {
	if uq.compStr(uq.prevLine, curLine) {
		uq.uniqPrev = false

		return nil
	}

	if !uq.uniqPrev {
		_, err := uq.Writeln([]byte(uq.prevLine))
		if err != nil {
			return fmt.Errorf("in calcHandlerLine(): Error is: %w", err)
		}
	}

	uq.prevLine = curLine
	uq.uniqPrev = true

	return nil
}

func (uq *UniqCommand) handleLastLine(handler func(string) error) error {
	uq.compStr = func(string, string) bool { return false }

	return handler("")
}

func (uq *UniqCommand) getEndHandler() func(string) error {
	switch {
	case uq.args.cFlag:
		return uq.calcHandleLine
	case uq.args.uFlag:
		return uq.uniqHandleLine
	case uq.args.dFlag:
		return uq.duplicateHandleLine
	default:
		return uq.defaultHandleLine
	}
}

func (uq *UniqCommand) Run() error {
	uq.initComparator()
	endHandler := uq.getEndHandler()

	scanner := bufio.NewScanner(uq.args.readCloser)
	if scanner.Scan() {
		uq.prevLine = scanner.Text()

		defer func(uq *UniqCommand) {
			err := uq.handleLastLine(endHandler)
			if err != nil {
				log.Printf("in Run() defer: Error is: %v", err)
			}
		}(uq)
	}

	for scanner.Scan() {
		err := endHandler(scanner.Text())
		if err != nil {
			return fmt.Errorf("in Run(): Error is: %w", err)
		}
	}

	return nil
}

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

	args.readCloser, args.writeCloser, err = parseFilesArgs(argsFiles)
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

	defer func() {
		err = uniqCommand.Close()
		if err != nil {
			log.Println(fmt.Errorf("error in main(): %w", err))
		}
	}()

	err = uniqCommand.Run()
	if err != nil {
		log.Println(fmt.Errorf("error in main(): %w", err))
	}
}
