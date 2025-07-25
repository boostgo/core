package sql

import (
	"context"
	"database/sql"

	"github.com/boostgo/core/contextx"
	"github.com/boostgo/core/convert"
	"github.com/boostgo/core/log"
	"github.com/boostgo/core/storage"

	"github.com/jmoiron/sqlx"
)

type clientSingle struct {
	conn      *sqlx.DB
	enableLog bool
}

// Client creates DB implementation by single client
func Client(conn *sqlx.DB, enableLog ...bool) DB {
	var enable bool
	if len(enableLog) > 0 {
		enable = enableLog[0]
	}

	return &clientSingle{
		conn:      conn,
		enableLog: enable,
	}
}

func (c *clientSingle) Connection() *sqlx.DB {
	return c.conn
}

func (c *clientSingle) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	if err := contextx.Validate(ctx); err != nil {
		return nil, err
	}

	c.printLog(ctx, "ExecContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return c.conn.ExecContext(ctx, query, args...)
}

func (c *clientSingle) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	if err := contextx.Validate(ctx); err != nil {
		return nil, err
	}

	c.printLog(ctx, "QueryContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return c.conn.QueryContext(ctx, query, args...)
}

func (c *clientSingle) QueryxContext(ctx context.Context, query string, args ...interface{}) (*sqlx.Rows, error) {
	if err := contextx.Validate(ctx); err != nil {
		return nil, err
	}

	c.printLog(ctx, "QueryxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryxContext(ctx, query, args...)
	}

	return c.conn.QueryxContext(ctx, query, args...)
}

func (c *clientSingle) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	c.printLog(ctx, "QueryRowxContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryRowxContext(ctx, query, args...)
	}

	return c.conn.QueryRowxContext(ctx, query, args...)
}

func (c *clientSingle) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if err := contextx.Validate(ctx); err != nil {
		return nil, err
	}

	c.printLog(ctx, "PrepareContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareContext(ctx, query)
	}

	return c.conn.PrepareContext(ctx, query)
}

func (c *clientSingle) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	if err := contextx.Validate(ctx); err != nil {
		return nil, err
	}

	c.printLog(ctx, "NamedExecContext", query, arg)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.NamedExecContext(ctx, query, arg)
	}

	return c.conn.NamedExecContext(ctx, query, arg)
}

func (c *clientSingle) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if err := contextx.Validate(ctx); err != nil {
		return err
	}

	c.printLog(ctx, "SelectContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return c.conn.SelectContext(ctx, dest, query, args...)
}

func (c *clientSingle) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	if err := contextx.Validate(ctx); err != nil {
		return err
	}

	c.printLog(ctx, "GetContext", query, args...)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.GetContext(ctx, dest, query, args...)
	}

	return c.conn.GetContext(ctx, dest, query, args...)
}

func (c *clientSingle) PrepareNamedContext(ctx context.Context, query string) (*sqlx.NamedStmt, error) {
	if err := contextx.Validate(ctx); err != nil {
		return nil, err
	}

	c.printLog(ctx, "PrepareNamedContext", query)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.PrepareNamedContext(ctx, query)
	}

	return c.conn.PrepareNamedContext(ctx, query)
}

func (c *clientSingle) NamedQueryRowxContext(ctx context.Context, query string, arg any) *sqlx.Row {
	if err := contextx.Validate(ctx); err != nil {
		return nil
	}

	c.printLog(ctx, "NamedQueryRowxContext", query)

	// convert name vars to ?
	namedQuery, args, err := sqlx.Named(query, arg)
	if err != nil {
		log.
			Error().
			Ctx(ctx).
			Err(err).
			Str("query", query).
			Obj("arg", arg).
			Msg("NamedQueryRowxContext: Named query error")
		return nil
	}

	// convert ? to format $1, $2, $3...
	convertedQuery := sqlx.Rebind(sqlx.DOLLAR, namedQuery)

	tx, ok := GetTx(ctx)
	if ok {
		return tx.QueryRowxContext(ctx, convertedQuery, args...)
	}

	return c.conn.QueryRowxContext(ctx, convertedQuery, args...)
}

func (c *clientSingle) EachShard(_ func(conn DB) error) error {
	return ErrMethodNotSuppoertedInSingle
}

func (c *clientSingle) EachShardAsync(_ func(conn DB) error, _ ...int) error {
	return ErrMethodNotSuppoertedInSingle
}

func (c *clientSingle) printLog(ctx context.Context, queryType, query string, args ...any) {
	if !c.enableLog || storage.IsNoLog(ctx) {
		return
	}

	convertedArgs := make([]string, 0, len(args))
	for _, arg := range args {
		convertedArgs = append(convertedArgs, convert.String(arg))
	}

	log.
		Info().
		Ctx(ctx).
		Str("queryType", queryType).
		Str("query", query).
		Strs("args", convertedArgs).
		Send()
}
