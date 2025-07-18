package queuex

import "sync"

// safeQueue wraps a queue with mutex for thread safety
type safeQueue[T any] struct {
	queue Queue[T]
	mu    sync.RWMutex
}

// newSafeQueue creates a thread-safe wrapper around any queue
func newSafeQueue[T any](q Queue[T]) Queue[T] {
	return &safeQueue[T]{
		queue: q,
	}
}

func (s *safeQueue[T]) Push(items ...T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.Push(items...)
}

func (s *safeQueue[T]) PushFront(items ...T) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.PushFront(items...)
}

func (s *safeQueue[T]) Pop(count ...int) ([]T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.Pop(count...)
}

func (s *safeQueue[T]) PopBack(count ...int) ([]T, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.queue.PopBack(count...)
}

func (s *safeQueue[T]) Size() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.queue.Size()
}

func (s *safeQueue[T]) IsEmpty() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.queue.IsEmpty()
}

func (s *safeQueue[T]) IsFull() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.queue.IsFull()
}

func (s *safeQueue[T]) Clear() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.queue.Clear()
}

func (s *safeQueue[T]) Peek() (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.queue.Peek()
}

func (s *safeQueue[T]) PeekBack() (T, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.queue.PeekBack()
}

func (s *safeQueue[T]) ToSlice() []T {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.queue.ToSlice()
}
