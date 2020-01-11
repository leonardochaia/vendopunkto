package invoice

import (
	vendopunkto "github.com/leonardochaia/vendopunkto/internal"
	"sync"
)

// topicImpl implements InvoiceTopic
// credits to https://github.com/tv42/topic/blob/master/topic.go
type topicImpl struct {
	// Producer sends messages on this channel. Close the channel
	// to shutdown the topic.
	Broadcast chan<- vendopunkto.Invoice

	lock        sync.Mutex
	connections map[string][]chan<- vendopunkto.Invoice
}

// NewTopic creates a new topic. Messages can be broadcast on this topic,
// and registered consumers are guaranteed to either receive them, or
// see a channel close.
func NewTopic() vendopunkto.InvoiceTopic {
	t := &topicImpl{}
	broadcast := make(chan vendopunkto.Invoice, 100)
	t.Broadcast = broadcast
	t.connections = make(map[string][]chan<- vendopunkto.Invoice)
	go t.run(broadcast)
	return t
}

func (t *topicImpl) run(broadcast <-chan vendopunkto.Invoice) {
	for invoice := range broadcast {
		func() {
			t.lock.Lock()
			defer t.lock.Unlock()
			subs, ok := t.connections[invoice.ID]
			if ok {
				for _, ch := range subs {
					ch <- invoice
				}
			}
		}()
	}

	t.lock.Lock()
	defer t.lock.Unlock()
	for id, subs := range t.connections {
		delete(t.connections, id)
		for _, ch := range subs {
			close(ch)
		}
	}
}

// Register starts receiving messages on the given channel. If a
// channel close is seen, either the topic has been shut down, or the
// consumer was too slow, and should re-register.
func (t *topicImpl) Register(invoiceID string) <-chan vendopunkto.Invoice {
	t.lock.Lock()
	defer t.lock.Unlock()

	ch := make(chan vendopunkto.Invoice)

	subs, ok := t.connections[invoiceID]
	if !ok {
		subs = []chan<- vendopunkto.Invoice{}
	}

	t.connections[invoiceID] = append(subs, ch)

	return ch
}

// Unregister stops receiving messages on this channel.
func (t *topicImpl) Unregister(invoiceID string) {
	t.lock.Lock()
	defer t.lock.Unlock()

	// double-close is not safe, so make sure we didn't already
	// drop this consumer as too slow
	subs, ok := t.connections[invoiceID]
	if ok {
		delete(t.connections, invoiceID)
		for _, ch := range subs {
			close(ch)
		}
	}
}

func (t *topicImpl) Send(msg vendopunkto.Invoice) {
	t.Broadcast <- msg
}

func (t *topicImpl) Close() {
	close(t.Broadcast)
}
