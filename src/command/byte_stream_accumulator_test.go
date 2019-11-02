package command

import (
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteStreamAccumulator(t *testing.T) {
	reader, writer := io.Pipe()
	acc := NewByteStreamAccumulator(reader)

	// write "one" to the stream and verify that it receives it
	go func() {
		writer.Write([]byte("one"))
	}()
	read := <-acc.Updates()
	assert.Equal(t, "one", read, "should provide the accumulated text in the notification")
	assert.Equal(t, "one", acc.String(), "should provide the accumulated text via String()")

	// write "two" to the stream and verify that it receives it
	go func() {
		writer.Write([]byte("two"))
	}()
	read = <-acc.Updates()
	assert.Equal(t, "onetwo", read, "should provide the accumulated text in the notification")
	assert.Equal(t, "onetwo", acc.String(), "should provide the accumulated text via String()")

	writer.Close()
}
