package uniq_test

import "bytes"

type nopCloseBuffer struct {
	bytes.Buffer
}

func (nopCloseBuffer) Close() error {
	return nil
}
