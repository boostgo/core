package sql

import (
	"context"
	"time"

	"github.com/boostgo/core/log"
	"github.com/boostgo/core/retry"

	"github.com/jmoiron/sqlx"
)

// StatementClose safely closes a statement with logging
func StatementClose(ctx context.Context, stmt *sqlx.NamedStmt) {
	if stmt == nil {
		return
	}

	if err := retry.Try(ctx, func(ctx context.Context) error {
		return stmt.Close()
	}, retry.Options{
		Policy: retry.NewFixedDelay(time.Millisecond, 3),
	}); err != nil {
		log.
			Error().
			Err(err).
			Msg("Close statement")
	}
}
