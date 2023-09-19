package uniq

import "io"

type ArgsCommand struct {
	CFlag bool
	DFlag bool
	UFlag bool
	IFlag bool

	FFlag uint
	SFlag uint

	ReadCloser  io.ReadCloser
	WriteCloser io.WriteCloser
}
