package storage

import (
	"github.com/boostgo/core/errorx"
)

var (
	// ErrConnNotSelected returns if "shard client" does not choose connection to use
	ErrConnNotSelected = errorx.New("sql.connection_not_selected")
)
