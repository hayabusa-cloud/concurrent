// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"errors"
)

var (
	// ErrTemporaryUnavailable is the error used for Enqueue operations
	// on a fulled queue or Dequeue operations on an empty queue
	ErrTemporaryUnavailable = errors.New("temporary unavailable")
)

// Pause executes CPU pause instructions to reduce energy consumption in spin-wait loops.
//
// Defaults to 20 cycles if not specified. Uses optimized assembly on amd64/arm64.
//
// Usage:
//
//	Pause()     // 20 cycles (default)
//	Pause(1)    // 1 cycle
//	Pause(40)   // 40 cycles
func Pause(cycles ...int) {
	n := 20
	if len(cycles) > 0 && cycles[0] > 0 {
		n = cycles[0]
	}
	pause(n)
}

// Producer is the interface that wraps the Enqueue method
type Producer[T any] interface {
	// Enqueue pushes item to FIFO queue.
	// if the queue is fulled, ErrTemporaryUnavailable will be returned
	Enqueue(elem *T) error
}

// Consumer is the interface that wraps the Dequeue method
type Consumer[T any] interface {
	// Dequeue pops items from FIFO queue.
	// if the queue is empty, ErrTemporaryUnavailable will be returned
	Dequeue() (elem *T, err error)
}

// Closer is the interface that wraps the Close method
type Closer interface {
	// Close closes the queue.
	// Enqueue and Dequeue operations on a closed queue is undefined
	Close() error
}
