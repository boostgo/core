package mongox

import (
	"github.com/boostgo/core/errorx"
	"go.mongodb.org/mongo-driver/mongo"
)

var (
	ErrCreateIndexes = errorx.New("mongo.create_indexes")
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
