package queuex

// Queue represents a generic queue interface
type Queue[T any] interface {
	Push(items ...T) error
	PushFront(items ...T) error
	Pop(count ...int) ([]T, error)
	PopBack(count ...int) ([]T, error)
	Size() int
	IsEmpty() bool
	IsFull() bool
	Clear()
	Peek() (T, error)
	PeekBack() (T, error)
	ToSlice() []T
}

// Config holds queue configuration
type Config struct {
	MaxSize    int // 0 means unlimited
	ThreadSafe bool
	Mode       Mode
}
