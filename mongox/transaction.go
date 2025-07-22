package mongox

import (
	"context"

	"github.com/boostgo/core/storage"
	"go.mongodb.org/mongo-driver/mongo"
)

type transaction struct {
	ctx     context.Context
	session mongo.Session
}

func newTransaction(ctx context.Context, session mongo.Session) storage.Transaction {
	return &transaction{
		ctx:     ctx,
		session: session,
	}
}

func (tx *transaction) Context() context.Context {
	if tx.ctx == nil {
		return context.Background()
	}

	return tx.ctx
}

func (tx *transaction) Commit(ctx context.Context) error {
	defer tx.session.EndSession(tx.ctx)

	if err := tx.session.CommitTransaction(ctx); err != nil {
		return ErrTxCommit.SetError(err)
	}

	return nil
}

func (tx *transaction) Rollback(ctx context.Context) error {
	defer tx.session.EndSession(ctx)

	if err := tx.session.AbortTransaction(ctx); err != nil {
		return ErrTxRollback.SetError(err)
	}

	return nil
}
