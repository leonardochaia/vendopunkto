package invoice

import (
	"sync"
)

// Topic is a pub-sub mechanism where consumers can Register to
// receive messages sent to using Send.
// credits to https://github.com/tv42/topic/blob/master/topic.go
type Topic interface {
	Register(invoiceID string) <-chan Invoice
	Unregister(invoiceID string)
	Send(msg Invoice)
	Close()
}

type topicImpl struct {
	// Producer sends messages on this channel. Close the channel
	// to shutdown the topic.
	Broadcast chan<- Invoice

	lock        sync.Mutex
	connections map[string]chan<- Invoice
}

// NewTopic creates a new topic. Messages can be broadcast on this topic,
// and registered consumers are guaranteed to either receive them, or
// see a channel close.
func NewTopic() Topic {
	t := &topicImpl{}
	broadcast := make(chan Invoice, 100)
	t.Broadcast = broadcast
	t.connections = make(map[string]chan<- Invoice)
	go t.run(broadcast)
	return t
}

func (t *topicImpl) run(broadcast <-chan Invoice) {
	for invoice := range broadcast {
		func() {
			t.lock.Lock()
			defer t.lock.Unlock()
			ch, ok := t.connections[invoice.ID]
			if ok {
				ch <- invoice
			}
		}()
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	for id, ch := range t.connections {
		delete(t.connections, id)
		close(ch)
	}
}

// Register starts receiving messages on the given channel. If a
// channel close is seen, either the topic has been shut down, or the
// consumer was too slow, and should re-register.
func (t *topicImpl) Register(invoiceID string) <-chan Invoice {
	t.lock.Lock()
	defer t.lock.Unlock()
	ch := make(chan Invoice)
	t.connections[invoiceID] = ch
	return ch
}

// Unregister stops receiving messages on this channel.
func (t *topicImpl) Unregister(invoiceID string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// double-close is not safe, so make sure we didn't already
	// drop this consumer as too slow
	ch, ok := t.connections[invoiceID]
	if ok {
		delete(t.connections, invoiceID)
		close(ch)
	}
}

func (t *topicImpl) Send(msg Invoice) {
	t.Broadcast <- msg
}

func (t *topicImpl) Close() {
	close(t.Broadcast)
}
