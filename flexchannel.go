// FlexChannel is a Go package that provides a non-blocking flexible channel implementation.
// It behaves like a channel with infinite size buffer,
// and it doesn't cause a panic even if you attempt to send after it has been closed.
package flexchannel

import (
	"errors"
	"sync"
)

var flexMessagePool = sync.Pool{
	New: func() interface{} {
		return &FlexChannelMessage{}
	},
}

// FlexChannelMessage represents a message in the FlexChannel.
type FlexChannelMessage struct {
	i    int
	data interface{}
}

// FlexChannel is a flexible channel with message pooling.
type FlexChannel struct {
	mu      sync.Mutex
	ch      chan *FlexChannelMessage
	started bool
	closed  bool
	seq     int
}

// NewFlexChannel creates a new FlexChannel.
func NewFlexChannel() *FlexChannel {
	fc := FlexChannel{
		ch: make(chan *FlexChannelMessage),
	}
	return &fc
}

// Start begins processing messages from the channel using the provided handler function.
// It returns an error if the channel has already been started.
func (fc *FlexChannel) Start(h func(interface{})) error {
	if h == nil {
		panic("handler is nil")
	}
	fc.mu.Lock()
	defer fc.mu.Unlock()
	if fc.started {
		return errors.New("already started")
	}
	fc.started = true
	go func() {
		cur := 0
		buf := make(map[int]interface{})
		var msg *FlexChannelMessage
		var ok bool
		for {
			msg, ok = <-fc.ch
			if !ok {
				fc.closed = true
				break
			}
			buf[msg.i] = msg.data
			flexMessagePool.Put(msg)
			for ; buf[cur] != nil; cur++ {
				h(buf[cur])
				delete(buf, cur)
			}
		}
	}()
	return nil
}

// Send sends data to the channel.
// It does nothing if the channel is closed and does not cause a panic.
func (fc *FlexChannel) Send(data interface{}) {
	if fc.closed {
		return
	}
	fc.mu.Lock()
	defer fc.mu.Unlock()
	seq := fc.seq
	fc.seq++
	go func() {
		defer (func() {
			recover()
		})()
		msg := flexMessagePool.Get().(*FlexChannelMessage)
		msg.i = seq
		msg.data = data
		fc.ch <- msg
	}()
}

// Close closes the channel.
func (fc *FlexChannel) Close() {
	fc.closed = true
	close(fc.ch)
}
