package scraper

import "sync"

type CircularQueue[T any] struct {
	sync.RWMutex
	capacity int
	elements []T
	head     int
	tail     int
}

func NewCircularQueue[T any](capacity int) *CircularQueue[T] {
	if capacity < 1 {
		// default capacity
		capacity = 100
	} else if capacity == 1 {
		capacity += 1
	}

	return &CircularQueue[T]{
		capacity: capacity,
		elements: make([]T, capacity),
		head:     0,
		tail:     0,
	}
}

func (queue *CircularQueue[T]) IsEmpty() bool {
	queue.RLock()
	defer queue.RUnlock()
	return queue.head == queue.tail
}

func (queue *CircularQueue[T]) IsFull() bool {
	queue.RLock()
	defer queue.RUnlock()
	return queue.head == (queue.tail+1)%queue.capacity
}

func (queue *CircularQueue[T]) Push(element T) {
	queue.RLock()
	defer queue.RUnlock()
	if queue.IsFull() {
		return
	}
	queue.elements[queue.tail] = element
	queue.tail = (queue.tail + 1) % queue.capacity
}

func (queue *CircularQueue[T]) Shift() (element *T) {
	queue.RLock()
	defer queue.RUnlock()
	element = &queue.elements[queue.head]
	queue.head = (queue.head + 1) % queue.capacity
	return
}
