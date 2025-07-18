package queuex

import "github.com/boostgo/core/errorx"

var (
	ErrQueueEmpty = errorx.New("queue.empty")
	ErrQueueFull  = errorx.New("queue.full")
)
