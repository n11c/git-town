package helpers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRemoveChannel(t *testing.T) {
	c1 := make(chan string)
	c2 := make(chan string)
	c3 := make(chan string)
	channels := []chan string{c1, c2, c3}

	removedMiddle := removeStringChannel(channels, 1)
	assert.Len(t, removedMiddle, 2)
	assert.Equal(t, removedMiddle[0], c1, "0 == c1")
	assert.Equal(t, removedMiddle[1], c3, "1 == c3")

	removedEnd := removeStringChannel(removedMiddle, 1)
	assert.Len(t, removedEnd, 1)
	assert.Equal(t, removedEnd[0], c1, "0 == c1")

	removedOnly := removeStringChannel(removedEnd, 0)
	assert.Len(t, removedOnly, 0)
}
