package mongox

import (
	"context"

	"github.com/boostgo/core/storage"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readconcern"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"go.mongodb.org/mongo-driver/mongo/writeconcern"
)

type TxOption func(opts *options.TransactionOptions)

const (
	sessionContextKey = "mongodb_session"
	txnContextKey     = "mongodb_transaction"
)

type transactor struct {
	client Client
	opts   *options.TransactionOptions
}

func NewTransactor(client Client, opts ...TxOption) storage.Transactor {
	txOptions := options.
		Transaction().
		SetReadConcern(readconcern.Majority()).
		SetWriteConcern(writeconcern.New(writeconcern.WMajority())).
		SetReadPreference(readpref.Primary())

	for _, opt := range opts {
		opt(txOptions)
	}

	return &transactor{
		client: client,
		opts:   txOptions,
	}
}

func (t *transactor) Key() string {
	return ""
}

func (t *transactor) IsTx(ctx context.Context) bool {
	inTxn, ok := ctx.Value(txnContextKey).(bool)
	return ok && inTxn
}

func (t *transactor) Begin(ctx context.Context) (storage.Transaction, error) {
	if t.IsTx(ctx) {
		session, ok := t.session(ctx)
		if !ok {
			return nil, ErrTxWrongObject
		}

		return newTransaction(ctx, session), nil
	}

	session, err := t.client.StartSession()
	if err != nil {
		return nil, ErrTxStartSession.SetError(err)
	}

	if err = session.StartTransaction(t.opts); err != nil {
		return nil, ErrTxStartTransaction.SetError(err)
	}

	ctx = mongo.NewSessionContext(ctx, session)

	return newTransaction(ctx, session), nil
}

func (t *transactor) BeginCtx(ctx context.Context) (context.Context, error) {
	if t.IsTx(ctx) {
		return ctx, nil
	}

	session, err := t.client.StartSession()
	if err != nil {
		return nil, ErrTxStartSession.SetError(err)
	}

	if err = session.StartTransaction(t.opts); err != nil {
		return nil, ErrTxStartTransaction.SetError(err)
	}

	ctx = mongo.NewSessionContext(ctx, session)
	newCtx := context.WithValue(ctx, sessionContextKey, session)
	newCtx = context.WithValue(newCtx, txnContextKey, true)

	return newCtx, nil
}

func (t *transactor) CommitCtx(ctx context.Context) error {
	session, ok := t.session(ctx)
	if !ok {
		return nil
	}

	if !t.IsTx(ctx) {
		return ErrTxNoTransaction
	}

	defer session.EndSession(ctx)

	if err := session.CommitTransaction(ctx); err != nil {
		return ErrTxCommit.SetError(err)
	}

	return nil
}

func (t *transactor) RollbackCtx(ctx context.Context) error {
	session, ok := t.session(ctx)
	if !ok {
		return nil
	}

	if !t.IsTx(ctx) {
		return ErrTxNoTransaction
	}

	defer session.EndSession(ctx)

	if err := session.AbortTransaction(ctx); err != nil {
		return ErrTxRollback.SetError(err)
	}

	return nil
}

func (t *transactor) TryCommit(ctx context.Context, err *error) {
	if err != nil {
		_ = t.RollbackCtx(ctx)
	}

	_ = t.CommitCtx(ctx)
}

func (t *transactor) session(ctx context.Context) (mongo.Session, bool) {
	session, ok := ctx.Value(sessionContextKey).(mongo.Session)
	return session, ok
}
