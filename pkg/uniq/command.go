package uniq

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type comparatorString func(first, second string) bool

type Command struct {
	Args               *ArgsCommand
	uniqPrev           bool
	duplicationCounter int
	prevLine           string
	compStr            comparatorString
}

func NewCommand(args *ArgsCommand) *Command {
	command := Command{duplicationCounter: 1, uniqPrev: true} //nolint:exhaustruct,exhaustivestruct

	if args.ReadCloser == nil {
		args.ReadCloser = os.Stdin
	}

	if args.WriteCloser == nil {
		args.WriteCloser = os.Stdout
	}

	command.initComparator()
	command.Args = args

	return &command //nolint:exhaustivestruct,exhaustruct
}

func (c *Command) Close() error {
	err := c.Args.ReadCloser.Close()
	if err != nil {
		return fmt.Errorf("in Close(): Error is: %w", err)
	}

	err = c.Args.WriteCloser.Close()
	if err != nil {
		return fmt.Errorf("in Close(): Error is: %w", err)
	}

	return nil
}

func (c *Command) writeln(p []byte) (int, error) {
	return c.Args.WriteCloser.Write(append(p, []byte("\n")...)) //nolint:wrapcheck
}

func (c *Command) fieldsHandleLine(curLine string) string {
	if c.Args.FFlag != 0 {
		curLine = TrimFirstFields(curLine, c.Args.FFlag)
	}

	return curLine
}

func (c *Command) skipCharsHandleLine(curLine string) string {
	if c.Args.SFlag != 0 {
		if uint(len(curLine)) >= c.Args.SFlag {
			curLine = curLine[c.Args.SFlag:]
		} else {
			curLine = ""
		}
	}

	return curLine
}

func (c *Command) ignoreCaseHandleLine(curLine string) string {
	if c.Args.IFlag {
		curLine = strings.ToLower(curLine)
	}

	return curLine
}

func (c *Command) initComparator() {
	c.compStr = func(first, second string) bool {
		first = c.ignoreCaseHandleLine(c.skipCharsHandleLine(c.fieldsHandleLine(first)))
		second = c.ignoreCaseHandleLine(c.skipCharsHandleLine(c.fieldsHandleLine(second)))

		return first == second
	}
}

func (c *Command) defaultHandleLine(curLine string) error {
	if !c.compStr(c.prevLine, curLine) {
		_, err := c.writeln([]byte(c.prevLine))
		if err != nil {
			return fmt.Errorf("in defaultHandleLine(): Error is: %w", err)
		}

		c.prevLine = curLine
	}

	return nil
}

func (c *Command) calcHandleLine(curLine string) error {
	if c.compStr(c.prevLine, curLine) {
		c.duplicationCounter++

		return nil
	}

	_, err := c.writeln([]byte(strconv.Itoa(c.duplicationCounter) + " " + c.prevLine))
	if err != nil {
		return fmt.Errorf("in calcHandlerLine(): Error is: %w", err)
	}

	c.prevLine = curLine
	c.duplicationCounter = 1

	return nil
}

func (c *Command) uniqHandleLine(curLine string) error {
	if c.compStr(c.prevLine, curLine) {
		c.uniqPrev = false

		return nil
	}

	if c.uniqPrev {
		_, err := c.writeln([]byte(c.prevLine))
		if err != nil {
			return fmt.Errorf("in calcHandlerLine(): Error is: %w", err)
		}
	}

	c.prevLine = curLine
	c.uniqPrev = true

	return nil
}

func (c *Command) duplicateHandleLine(curLine string) error {
	if c.compStr(c.prevLine, curLine) {
		c.uniqPrev = false

		return nil
	}

	if !c.uniqPrev {
		_, err := c.writeln([]byte(c.prevLine))
		if err != nil {
			return fmt.Errorf("in calcHandlerLine(): Error is: %w", err)
		}
	}

	c.prevLine = curLine
	c.uniqPrev = true

	return nil
}

func (c *Command) handleLastLine(handler func(string) error) error {
	c.compStr = func(string, string) bool { return false }

	return handler("")
}

func (c *Command) getEndHandler() func(string) error {
	switch {
	case c.Args.CFlag:
		return c.calcHandleLine
	case c.Args.UFlag:
		return c.uniqHandleLine
	case c.Args.DFlag:
		return c.duplicateHandleLine
	default:
		return c.defaultHandleLine
	}
}

func (c *Command) Run() error {
	endHandler := c.getEndHandler()

	scanner := bufio.NewScanner(c.Args.ReadCloser)
	if scanner.Scan() {
		c.prevLine = scanner.Text()

		defer func(c *Command) {
			err := c.handleLastLine(endHandler)
			if err != nil {
				log.Printf("in Run() defer: Error is: %v", err)
			}
		}(c)
	}

	for scanner.Scan() {
		err := endHandler(scanner.Text())
		if err != nil {
			return fmt.Errorf("in Run(): Error is: %w", err)
		}
	}

	return nil
}
