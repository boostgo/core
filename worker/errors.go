package worker

import "github.com/boostgo/core/errorx"

var (
	ErrLocked = errorx.New("worker.locker.locked")
)
