package mongox

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MigrationBuilder struct {
	client Client
	ops    []func(ctx context.Context) error
}

func NewMigrationBuilder(client Client) *MigrationBuilder {
	return &MigrationBuilder{
		client: client,
		ops:    make([]func(ctx context.Context) error, 0),
	}
}

func (mb *MigrationBuilder) CreateCollection(name string, validator bson.M) *MigrationBuilder {
	mb.ops = append(mb.ops, func(ctx context.Context) error {
		opts := options.CreateCollection()
		if validator != nil {
			opts.SetValidator(validator)
		}

		return mb.client.Database().CreateCollection(ctx, name, opts)
	})

	return mb
}

func (mb *MigrationBuilder) CreateIndex(collection string, keys bson.D, opts *options.IndexOptions) *MigrationBuilder {
	mb.ops = append(mb.ops, func(ctx context.Context) error {
		coll := mb.client.Collection(collection)
		_, err := coll.Indexes().CreateOne(ctx, mongo.IndexModel{
			Keys:    keys,
			Options: opts,
		})
		return err
	})

	return mb
}

func (mb *MigrationBuilder) InsertDocuments(collection string, documents []interface{}) *MigrationBuilder {
	mb.ops = append(mb.ops, func(ctx context.Context) error {
		coll := mb.client.Collection(collection)
		_, err := coll.InsertMany(ctx, documents)
		return err
	})

	return mb
}

func (mb *MigrationBuilder) UpdateDocuments(collection string, filter, update bson.M) *MigrationBuilder {
	mb.ops = append(mb.ops, func(ctx context.Context) error {
		coll := mb.client.Collection(collection)
		_, err := coll.UpdateMany(ctx, filter, update)
		return err
	})

	return mb
}

func (mb *MigrationBuilder) Execute(ctx context.Context) error {
	for i, op := range mb.ops {
		if err := op(ctx); err != nil {
			return newMigrateError(err, i+1)
		}
	}

	return nil
}
