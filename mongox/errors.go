package mongox

import (
	"github.com/boostgo/core/errorx"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCreateIndexes = errorx.New("mongo.create_indexes")

	ErrReadPrefInvalidMode = errorx.New("mongo.read_pref.invalid_mode")
	ErrReadPrefCreate      = errorx.New("mongo.read_pref.create")

	ErrConcernWriteUnsupported = errorx.New("mongo.concern.write.unsupported")
	ErrConcernReadUnsupported  = errorx.New("mongo.concern.read.unsupported")

	ErrMigrate = errorx.New("mongo.migrate")

	ErrTxWrongObject      = errorx.New("mongo.tx.wrong_object")
	ErrTxStartSession     = errorx.New("mongo.tx.start_session")
	ErrTxStartTransaction = errorx.New("mongo.tx.start_transaction")
	ErrTxNoTransaction    = errorx.New("mongo.tx.no_transaction")
	ErrTxCommit           = errorx.New("mongo.tx.commit")
	ErrTxRollback         = errorx.New("mongo.tx.rollback")
)

type createIndexContext struct {
	Collection string             `json:"collection"`
	Indexes    []mongo.IndexModel `json:"indexes"`
	Error      error              `json:"error"`
}

func newCreateIndexError(err error, collection string, indexes []mongo.IndexModel) error {
	return ErrCreateIndexes.
		SetError(err).
		SetData(createIndexContext{
			Collection: collection,
			Indexes:    indexes,
			Error:      err,
		})
}

type migrateContext struct {
	Num   int   `json:"num"`
	Error error `json:"error"`
}

func newMigrateError(err error, num int) error {
	return ErrMigrate.
		SetError(err).
		SetData(migrateContext{
			Num:   num,
			Error: err,
		})
}

type unsupportedConcernContext struct {
	Provide     string `json:"provide"`
	ConcernType string `json:"concern_type"`
}

func newUnsupportedConcernError(err *errorx.Error, provide, concernType string) error {
	return err.
		SetData(unsupportedConcernContext{
			Provide:     provide,
			ConcernType: concernType,
		})
}
