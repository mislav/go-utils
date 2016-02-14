package utils

import (
	"io"
)

func ClosingTeeReader(origin io.ReadCloser, destination io.WriteCloser) *closingTeeReader {
	return &closingTeeReader{
		origin:      origin,
		destination: destination,
		Reader:      io.TeeReader(origin, destination),
	}
}

type closingTeeReader struct {
	origin      io.ReadCloser
	destination io.WriteCloser
	io.Reader
}

func (c *closingTeeReader) Close() error {
	err1 := c.destination.Close()
	err2 := c.origin.Close()
	if err1 == nil {
		return err2
	} else {
		return err1
	}
}
