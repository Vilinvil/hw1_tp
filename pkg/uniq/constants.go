package uniq

import (
	"io/fs"
)

const (
	inputFileExist  int8 = 1
	outputFileExist int8 = 2

	readeWriteEnable fs.FileMode = 0o666
)
