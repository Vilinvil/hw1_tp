package uniq

import (
	"errors"
	"io/fs"
)

const helpMessage = `uniq [-c | -d | -u] [-i] [-f num] [-s chars] [input_file [output_file]]

With no options uniq remove duplicate lines and show first occurrence

Options`

var ErrHelp = errors.New("for help use argument -h or --help")

const (
	inputFileExist   int8        = 1
	outputFileExist  int8        = 2
	readeWriteEnable fs.FileMode = 0o666
)
