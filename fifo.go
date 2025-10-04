// ©Hayabusa Cloud Co., Ltd. 2023. All rights reserved.
// Use of this source code is governed by a MIT-style
// license that can be found in the LICENSE file.

package concurrent

import (
	"unsafe"
)

// SPSCQueue represents simple producer single consumer FIFO queue
type SPSCQueue[T any] struct{}

// NewSPSCQueue creates a new simple producer single consumer
// FIFO queue with the given capacity
func NewSPSCQueue[T any](capacity int) (Consumer[T], Producer[T]) {
	panic("not implement")
}

// MPSCQueue represents multiple producers single consumer FIFO queue
type MPSCQueue[T any] struct{}

// NewMPSCQueue creates a new multiple producers single consumer
// FIFO queue with the given capacity
func NewMPSCQueue[T any](capacity int) (Consumer[T], Producer[T]) {
	panic("not implement")
}

// SPMCQueue represents single producer multiple consumers FIFO queue
type SPMCQueue[T any] struct{}

// NewSPMCQueue creates a new single producer multiple consumers
// FIFO queue with the given capacity
func NewSPMCQueue[T any](capacity int) (Consumer[T], Producer[T]) {
	panic("not implement")
}

// MPMCQueue represents multiple producers multiple consumers FIFO queue
type MPMCQueue[T any] struct {
	*rmfLF
}

// NewMPMCQueue creates a new multiple producers multiple consumers
// FIFO queue with the given capacity
func NewMPMCQueue[T any](capacity int) (Consumer[T], Producer[T]) {
	if capacity < 2 {
		panic("bad capacity")
	}
	capacity--
	order := 0
	for capacity > 0 {
		order++
		capacity >>= 1
	}
	q := MPMCQueue[T]{rmfLF: newRmfLF(order)}

	return &q, &q
}

// Enqueue pushes the given item to a FIFO queue
func (q *MPMCQueue[T]) Enqueue(elem *T) error {
	ok := q.offer(uintptr(unsafe.Pointer(elem)))
	if !ok {
		return ErrTemporaryUnavailable
	}

	return nil
}

// Dequeue pops items from FIFO queue
func (q *MPMCQueue[T]) Dequeue() (elem *T, err error) {
	ptr, ok := q.poll()
	if !ok {
		return elem, ErrTemporaryUnavailable
	}
	elem = (*T)(unsafe.Pointer(ptr))

	return
}

// EnqueueWait pushes the given item to a fifo queue.
// the operation will block until success or error occurred
func EnqueueWait[T any](p Producer[T], elem *T) error {
	sw := SpinWait{}
	for {
		err := p.Enqueue(elem)
		if err == ErrTemporaryUnavailable {
			sw.Once()
			continue
		}

		return err
	}
}

// DequeueWait pops items from fifo queue.
// the operation will block until
func DequeueWait[T any](c Consumer[T]) (elem *T, err error) {
	sw := SpinWait{}
	for {
		elem, err = c.Dequeue()
		if err == ErrTemporaryUnavailable {
			sw.Once()
			continue
		}

		return
	}
}
