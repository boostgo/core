package queuex

// baseQueue implements basic queue operations
type baseQueue[T any] struct {
	items   []T
	maxSize int
	mode    Mode
}

// New creates a new queue with given configuration
func New[T any](cfg Config) Queue[T] {
	q := &baseQueue[T]{
		items:   make([]T, 0),
		maxSize: cfg.MaxSize,
		mode:    cfg.Mode,
	}

	if cfg.ThreadSafe {
		return newSafeQueue(q)
	}

	return q
}

// Push adds items to the queue
func (q *baseQueue[T]) Push(items ...T) error {
	if len(items) == 0 {
		return nil
	}

	// Calculate how many items we can add
	if q.maxSize > 0 {
		availableSpace := q.maxSize - len(q.items)
		toEvict := len(items) - availableSpace

		if toEvict > 0 {
			// Need to evict some items
			if toEvict >= len(q.items) {
				// All current items will be evicted
				q.items = q.items[:0]
			} else {
				// Evict oldest items
				q.items = q.items[toEvict:]
			}
		}

		// If we still have too many items, take only the newest ones
		if len(items) > q.maxSize {
			items = items[len(items)-q.maxSize:]
		}
	}

	q.items = append(q.items, items...)
	return nil
}

// Pop removes and returns items from the queue
func (q *baseQueue[T]) Pop(count ...int) ([]T, error) {
	n := 1
	if len(count) > 0 && count[0] > 0 {
		n = count[0]
	}

	if len(q.items) == 0 {
		return nil, ErrQueueEmpty
	}

	// Adjust n if it's more than available items
	if n > len(q.items) {
		n = len(q.items)
	}

	var items []T
	switch q.mode {
	case ModeFIFO:
		items = make([]T, n)
		copy(items, q.items[:n])
		q.items = q.items[n:]
	case ModeLIFO:
		start := len(q.items) - n
		items = make([]T, n)
		copy(items, q.items[start:])
		q.items = q.items[:start]
	}

	return items, nil
}

// PushFront adds items to the front of the queue
func (q *baseQueue[T]) PushFront(items ...T) error {
	if len(items) == 0 {
		return nil
	}

	// Calculate total size after adding
	totalSize := len(q.items) + len(items)

	if q.maxSize > 0 && totalSize > q.maxSize {
		// Need to evict from the back
		toKeep := q.maxSize - len(items)
		if toKeep > 0 {
			q.items = q.items[:toKeep]
		} else {
			q.items = q.items[:0]
			// If items are more than maxSize, keep only the last maxSize items
			if len(items) > q.maxSize {
				items = items[:q.maxSize]
			}
		}
	}

	q.items = append(items, q.items...)
	return nil
}

// PopBack removes and returns items from the back of the queue
func (q *baseQueue[T]) PopBack(count ...int) ([]T, error) {
	n := 1
	if len(count) > 0 && count[0] > 0 {
		n = count[0]
	}

	if len(q.items) == 0 {
		return nil, ErrQueueEmpty
	}

	// Adjust n if it's more than available items
	if n > len(q.items) {
		n = len(q.items)
	}

	start := len(q.items) - n
	items := make([]T, n)
	copy(items, q.items[start:])

	// For FIFO, PopBack should return items in reverse order (newest first)
	if q.mode == ModeFIFO {
		// Reverse the slice
		for i := 0; i < len(items)/2; i++ {
			items[i], items[len(items)-1-i] = items[len(items)-1-i], items[i]
		}
	}

	q.items = q.items[:start]

	return items, nil
}

// Size returns the number of items in the queue
func (q *baseQueue[T]) Size() int {
	return len(q.items)
}

// IsEmpty returns true if the queue is empty
func (q *baseQueue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// IsFull returns true if the queue has reached max size
func (q *baseQueue[T]) IsFull() bool {
	return q.maxSize > 0 && len(q.items) >= q.maxSize
}

// Clear removes all items from the queue
func (q *baseQueue[T]) Clear() {
	q.items = q.items[:0]
}

// evictOldest removes the oldest item based on queue mode
func (q *baseQueue[T]) evictOldest() {
	if len(q.items) == 0 {
		return
	}

	switch q.mode {
	case ModeFIFO:
		q.items = q.items[1:]
	case ModeLIFO:
		q.items = q.items[:len(q.items)-1]
	}
}

// evictNewest removes the newest item based on queue mode
func (q *baseQueue[T]) evictNewest() {
	if len(q.items) == 0 {
		return
	}

	switch q.mode {
	case ModeFIFO:
		q.items = q.items[:len(q.items)-1]
	case ModeLIFO:
		q.items = q.items[1:]
	}
}

// Additional utility methods for baseQueue
func (q *baseQueue[T]) Peek() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, ErrQueueEmpty
	}

	switch q.mode {
	case ModeFIFO:
		return q.items[0], nil
	case ModeLIFO:
		return q.items[len(q.items)-1], nil
	}

	return zero, nil
}

func (q *baseQueue[T]) PeekBack() (T, error) {
	var zero T
	if len(q.items) == 0 {
		return zero, ErrQueueEmpty
	}

	return q.items[len(q.items)-1], nil
}

func (q *baseQueue[T]) ToSlice() []T {
	result := make([]T, len(q.items))
	copy(result, q.items)
	return result
}
