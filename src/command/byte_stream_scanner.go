package command

import (
	"io"
	"strings"
	"sync"
)

// ByteStreamScanner allows observing an io.ReadCloser stream for textual content.
type ByteStreamScanner struct {
	acc           *ByteStreamAccumulator
	subscriptions subscriptionList
	subMux        sync.Mutex // synchronizes access to subscriptions
}

// NewByteStreamScanner provides instances of ByteStreamScanner.
func NewByteStreamScanner(stream io.ReadCloser) *ByteStreamScanner {
	scanner := ByteStreamScanner{acc: NewByteStreamAccumulator(stream)}
	go scanner.processStreamUpdates()
	return &scanner
}

// checkSubscriptions notifies subscribers about possible matches in the given accumulated text.
func (scanner *ByteStreamScanner) checkSubscriptions(accumulated string) {
	scanner.subMux.Lock()
	defer scanner.subMux.Unlock()
	for i := 0; i < len(scanner.subscriptions); i++ {
		subscription := scanner.subscriptions[i]
		if subscription.Matches(accumulated) {
			subscription.Notify()
			scanner.subscriptions = scanner.subscriptions.Remove(i)
			i--
		}
	}
}

// processStreamUpdates handles updates received from the stream.
func (scanner *ByteStreamScanner) processStreamUpdates() {
	for {
		streamText := <-scanner.acc.Updates()
		scanner.checkSubscriptions(streamText)
	}
}

// ReceivedText provides the text received from the stream so far.
func (scanner *ByteStreamScanner) ReceivedText() string {
	return scanner.acc.String()
}

// WaitForText returns a channel that emits the given text
// when it occurs in the observed byte stream.
// Each subscription in notified only once, then the channel is closed.
func (scanner *ByteStreamScanner) WaitForText(text string) chan string {
	result := make(chan string)
	scanner.subMux.Lock()
	scanner.subscriptions = scanner.subscriptions.Add(Subscription{query: text, channel: result})
	scanner.subMux.Unlock()
	scanner.checkSubscriptions(scanner.acc.accumulated)
	return result
}

// Subscription represents a request of a client to be notified when the given text occurs in the stream.
type Subscription struct {
	query   string      // the text to look for
	channel chan string // the channel to notify the subscriber
}

// NewSubscription provides new subscription instances.
func NewSubscription(query string) Subscription {
	return Subscription{query: query, channel: make(chan string)}
}

// Channel provides the channel that notifies clients when the query has been found.
func (sub *Subscription) Channel() chan string {
	return sub.channel
}

// Matches indicates whether this subscription matches the given text.
func (sub *Subscription) Matches(text string) bool {
	return strings.Contains(text, sub.query)
}

// Notify informs the client for this subscription that the text has been found
func (sub *Subscription) Notify() {
	go func() {
		sub.channel <- sub.query
		close(sub.channel)
	}()
}

// subscriptionList is a list of subscriptions
type subscriptionList []Subscription

// Add appends the given Subscription to this list.
func (list subscriptionList) Add(sub ...Subscription) subscriptionList {
	return append(list, sub...)
}

// Remove provides a subscriptionList with the element at the given index removed.
func (list subscriptionList) Remove(index int) subscriptionList {
	list[index] = list[len(list)-1]
	return list[:len(list)-1]
}
