package uniq

import (
	"bytes"
)

type nopCloseBuffer struct {
	bytes.Buffer
}

func (nopCloseBuffer) Close() error {
	return nil
}
