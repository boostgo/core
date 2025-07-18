package queuex

import (
	"sync"
	"testing"
)

func TestBasicQueue(t *testing.T) {
	t.Run("FIFO operations", func(t *testing.T) {
		q := New[int](Config{Mode: ModeFIFO})

		// Test Push and Pop
		q.Push(1, 2, 3)

		if q.Size() != 3 {
			t.Errorf("expected size 3, got %d", q.Size())
		}

		val, err := q.Pop()
		if err != nil || val[0] != 1 {
			t.Errorf("expected 1, got %d, err: %v", val, err)
		}

		val, err = q.Pop()
		if err != nil || val[0] != 2 {
			t.Errorf("expected 2, got %d, err: %v", val, err)
		}
	})

	t.Run("LIFO operations", func(t *testing.T) {
		q := New[int](Config{Mode: ModeLIFO})

		q.Push(1)
		q.Push(2)
		q.Push(3)

		val, err := q.Pop()
		if err != nil || val[0] != 3 {
			t.Errorf("expected 3, got %d, err: %v", val, err)
		}

		val, err = q.Pop()
		if err != nil || val[0] != 2 {
			t.Errorf("expected 2, got %d, err: %v", val, err)
		}
	})

	t.Run("empty queue", func(t *testing.T) {
		q := New[string](Config{})

		if !q.IsEmpty() {
			t.Error("new queue should be empty")
		}

		_, err := q.Pop()
		if err != ErrQueueEmpty {
			t.Errorf("expected ErrQueueEmpty, got %v", err)
		}

		_, err = q.PopBack()
		if err != ErrQueueEmpty {
			t.Errorf("expected ErrQueueEmpty, got %v", err)
		}
	})

	t.Run("max size with auto-eviction", func(t *testing.T) {
		q := New[int](Config{MaxSize: 3, Mode: ModeFIFO})

		q.Push(1, 2, 3, 4)

		if q.Size() != 3 {
			t.Errorf("expected size 3, got %d", q.Size())
		}

		val, _ := q.Pop()
		if val[0] != 2 {
			t.Errorf("expected 2 (1 was evicted), got %d", val)
		}
	})

	t.Run("PushFront and PopBack", func(t *testing.T) {
		q := New[int](Config{Mode: ModeFIFO})

		q.Push(2)
		q.PushFront(1)
		q.Push(3)

		// Queue should be [1, 2, 3]
		val, err := q.PopBack()
		if err != nil || val[0] != 3 {
			t.Errorf("expected 3, got %d, err: %v", val, err)
		}

		val, err = q.Pop()
		if err != nil || val[0] != 1 {
			t.Errorf("expected 1, got %d, err: %v", val, err)
		}
	})
}

func TestThreadSafety(t *testing.T) {
	q := New[int](Config{ThreadSafe: true, MaxSize: 100})

	var wg sync.WaitGroup
	numGoroutines := 10
	itemsPerGoroutine := 100

	// Concurrent pushes
	for i := 0; i < numGoroutines; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				q.Push(id*1000 + j)
			}
		}(i)
	}

	// Concurrent pops
	for i := 0; i < numGoroutines/2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < itemsPerGoroutine; j++ {
				q.Pop()
			}
		}()
	}

	wg.Wait()

	// Check final state
	if q.Size() > 100 {
		t.Errorf("size should not exceed max size, got %d", q.Size())
	}
}

func TestClear(t *testing.T) {
	q := New[string](Config{})

	q.Push("a")
	q.Push("b")
	q.Push("c")

	q.Clear()

	if !q.IsEmpty() {
		t.Error("queue should be empty after Clear()")
	}

	if q.Size() != 0 {
		t.Errorf("size should be 0, got %d", q.Size())
	}
}

func TestBatchOperations(t *testing.T) {
	t.Run("batch push", func(t *testing.T) {
		q := New[int](Config{Mode: ModeFIFO})

		q.Push(1, 2, 3, 4, 5)

		if q.Size() != 5 {
			t.Errorf("expected size 5, got %d", q.Size())
		}

		items, _ := q.Pop(3)
		if len(items) != 3 || items[0] != 1 || items[2] != 3 {
			t.Errorf("expected [1,2,3], got %v", items)
		}
	})

	t.Run("batch with max size", func(t *testing.T) {
		q := New[int](Config{MaxSize: 3, Mode: ModeFIFO})

		q.Push(1, 2, 3, 4, 5) // Should keep only [3,4,5]

		items, _ := q.Pop(10) // Request more than available
		if len(items) != 3 || items[0] != 3 {
			t.Errorf("expected [3,4,5], got %v", items)
		}
	})
}
