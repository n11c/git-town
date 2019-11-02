package command

import (
	"io"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestByteStreamScanner_receiveAfterSubscribe(t *testing.T) {
	reader, writer := io.Pipe()
	scanner := NewByteStreamScanner(reader)
	helloChan := scanner.WaitForText("hello")
	go func() {
		time.Sleep(10 * time.Microsecond)
		writer.Write([]byte("one"))
		writer.Write([]byte("a hello to you!"))
		writer.Write([]byte("two"))
		writer.Close()
	}()
	received := <-helloChan
	assert.Equal(t, "hello", received)
}

func TestByteStreamScanner_receiveBeforeSubscribe(t *testing.T) {
	reader, writer := io.Pipe()
	scanner := NewByteStreamScanner(reader)
	go func() {
		writer.Write([]byte("query"))
		writer.Close()
	}()
	time.Sleep(10 * time.Microsecond)
	helloChan := scanner.WaitForText("query")
	received := <-helloChan
	assert.Equal(t, "query", received)
}

func TestSubscription_Matches(t *testing.T) {
	tests := map[string]bool{
		"a hello to you": true,
		"zonk":           false,
	}
	sub := NewSubscription("hello")
	for input, expected := range tests {
		assert.Equal(t, expected, sub.Matches(input))
	}
}

func TestSubscription_Notify(t *testing.T) {
	sub := NewSubscription("query")
	go func() {
		sub.Notify()
	}()
	received := <-sub.Channel()
	assert.Equal(t, "query", received)
}

func TestSubScriptionList(t *testing.T) {
	s1 := NewSubscription("one")
	s2 := NewSubscription("two")
	s3 := NewSubscription("three")
	list := subscriptionList{s1}
	list = list.Add(s2, s3)

	removedMiddle := list.Remove(1)
	assert.Len(t, removedMiddle, 2)
	assert.Equal(t, s1, removedMiddle[0])
	assert.Equal(t, s3, removedMiddle[1])

	removedEnd := removedMiddle.Remove(1)
	assert.Len(t, removedEnd, 1)
	assert.Equal(t, s1, removedEnd[0])

	removedBeginning := removedEnd.Remove(0)
	assert.Len(t, removedBeginning, 0)
}
