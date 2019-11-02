package command

import (
	"fmt"
	"io"
)

// ByteStreamAccumulator accumulates the textual content received from an io.ReadCloser.
// It notifies clients about updates from the stream over a channel.
type ByteStreamAccumulator struct {
	accumulated string      // the text accumulated so far
	observed    io.Reader   // the pipe to listen to
	updates     chan string // channel that provides updated accumulated values
}

// NewByteStreamAccumulator provides instances ByteStreamAccumulator that observe the given io.Reader.
func NewByteStreamAccumulator(stream io.Reader) *ByteStreamAccumulator {
	result := ByteStreamAccumulator{observed: stream, updates: make(chan string)}
	go result.read()
	return &result
}

// String returns the content accumulated so far.
func (acc *ByteStreamAccumulator) String() string {
	return acc.accumulated
}

func (acc *ByteStreamAccumulator) read() {
	p := make([]byte, 100)
	read, err := acc.observed.Read(p)
	if read > 0 {
		acc.accumulated = acc.accumulated + string(p)[:read]
		acc.updates <- acc.accumulated
	}
	switch err {
	case nil: // no error --> keep reading
		acc.read()
	case io.EOF: // EOF --> done here
	default: // other error
		fmt.Printf("unexpected error: %v\n", err)
	}
}

// Updates provides the channel that sends content updates.
func (acc *ByteStreamAccumulator) Updates() chan string {
	return acc.updates
}
