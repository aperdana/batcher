// Package batcher provides mechanism to perform batch operation.
package batcher

import (
	"time"
)

const (
	// DefaultBufferSize is Batcher's default internal buffer size.
	DefaultBufferSize = 10
)

// Batcher is a time-based batcher.
// When SetMaxBatchSize is used when initializing new Batcher,
// it will also add size as constraint.
type Batcher struct {
	size    int
	bufSize int
	timeout time.Duration
	buffer  chan interface{}
	do      func([]interface{})

	listening bool // flag to ensure Batcher is listening for message
}

// Option used to control Batcher
type Option func(*Batcher)

// SetBufferSize sets Batcher's internal buffer size.
// As buffer size increases, Batcher's ability to
// concurrently handles Batch() also increases.
func SetBufferSize(s int) Option {
	return func(b *Batcher) {
		if s > 0 {
			b.bufSize = s
		}
	}
}

// SetMaxBatchSize sets maximum batch size for every round
// of batching. When s > 0, Batcher will not only
// work as a time-limited batcher, but also size-limited.
func SetMaxBatchSize(s int) Option {
	return func(b *Batcher) {
		if s > 0 {
			b.size = s
		}
	}
}

// New returns Batcher that runs do on batched items every timeout.
func New(do func([]interface{}), timeout time.Duration, opts ...Option) *Batcher {
	b := &Batcher{
		timeout: timeout,
		do:      do,
		bufSize: DefaultBufferSize,
	}
	for _, opt := range opts {
		opt(b)
	}
	b.buffer = make(chan interface{}, b.bufSize)

	return b
}

// Batch sends m to be batched.
// Listen() must be called first before performing any Batch(),
// otherwise it will panic.
//
// It is safe to use Batch concurrently.
func (b *Batcher) Batch(m interface{}) {
	if !b.listening {
		panic("batcher is not yet listening for input")
	}
	b.buffer <- m
}

// Listen starts a message listener in background.
func (b *Batcher) Listen() {
	if b.listening {
		return
	}
	b.listening = true
	go b.listen()
}

func (b *Batcher) listen() {
	for {
		batched := b.receive()
		b.do(batched)
	}
}

func (b *Batcher) receive() []interface{} {
	batched := make([]interface{}, 0)

	timer := time.NewTimer(b.timeout)
	defer timer.Stop()

	var i int
	// Try to get as much items from buffer before timeouts.
	// When b.size is set, the number of items fetched from
	// buffer will be b.size at most.
	for b.size == 0 || i < b.size {
		select {
		case <-timer.C:
			return batched
		case m := <-b.buffer:
			batched = append(batched, m)
			i++
		}
	}
	return batched
}
