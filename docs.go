// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"errors"
	"slices"
)

var (
	// ErrTemporaryUnavailable is the error used for Enqueue operations
	// on a fulled queue or Dequeue operations on an empty queue
	ErrTemporaryUnavailable = errors.New("temporary unavailable")
)

// QueueOptions is a struct that contains options for creating a queue.
type QueueOptions struct {
	SingleProducer bool // TODO: implement
	SingleConsumer bool // TODO: implement
	LowContention  bool // TODO: implement
	DistinctValues bool
}

var defaultQueueOptions = QueueOptions{
	DistinctValues: true,
}

// NewQueue creates a new queue with the given capacity and options.
func NewQueue[T any](capacity int, opts ...func(opts *QueueOptions)) (Consumer[T], Producer[T]) {
	opt := defaultQueueOptions
	for o := range slices.Values(opts) {
		o(&opt)
	}
	if !opt.DistinctValues {
		panic("not implement")
	}
	c, p := NewMPMCQueue[T](capacity)
	return c, p
}

// Producer is the interface that wraps the Enqueue method
type Producer[T any] interface {
	// Enqueue pushes item to FIFO queue.
	// if the queue is fulled, ErrTemporaryUnavailable will be returned
	Enqueue(elem *T) error
}

// Consumer is the interface that wraps the Dequeue method
type Consumer[T any] interface {
	// Dequeue pops items from the FIFO queue.
	// if the queue is empty, ErrTemporaryUnavailable will be returned
	Dequeue() (elem *T, err error)
}

// ProducerIndirect is a non-generic producer interface that enqueues uintptr values.
// Use this for indirect references such as indices, handles, or other integer-like identifiers.
type ProducerIndirect interface {
	Enqueue(elem uintptr) error
}

// ConsumerIndirect is a non-generic consumer interface that dequeues uintptr values.
// Use this for indirect references such as indices, handles, or other integer-like identifiers.
type ConsumerIndirect interface {
	Dequeue() (elem uintptr, err error)
}

// Closer is the interface that wraps the Close method
type Closer interface {
	// Close closes the queue.
	// Enqueue and Dequeue operations on a closed queue are undefined
	Close() error
}
