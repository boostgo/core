package sql

import (
	"github.com/boostgo/core/errorx"
	"github.com/boostgo/core/sql/duplicate"
)

var (
	ErrOpenConnect = errorx.New("sql.open_connect")
	ErrPing        = errorx.New("sql.ping")

	ErrConnectionStringEmpty     = errorx.New("sql.connection_string_empty")
	ErrConnectionStringDuplicate = errorx.New("sql.connection_string_duplicate")
	ErrConnectionIsNotShard      = errorx.New("sql.connection_is_not_shard")

	ErrMethodNotSuppoertedInSingle = errorx.New("sql.method_not_suppoerted_in_single")

	ErrMigrateOpenConn          = errorx.New("migrate.open_conn")
	ErrMigrateGetDriver         = errorx.New("migrate.get_driver")
	ErrMigrateLock              = errorx.New("migrate.lock")
	ErrMigrateReadMigrationsDir = errorx.New("migrate.read_migrations_dir")
	ErrMigrateUp                = errorx.New("migrate.up")

	ErrTransactorBegin    = errorx.New("transactor.begin")
	ErrTransactorCommit   = errorx.New("transactor.commit")
	ErrTransactorRollback = errorx.New("transactor.rollback")

	ErrStatementPrepare      = errorx.New("sql.statement.prepare")
	ErrStatementExecuteQuery = errorx.New("sql.statement.execute_query")

	ErrDuplicate           = errorx.New("sql.duplicate")
	ErrForeignKeyViolation = errorx.New("sql.foreign_key_violation")
	ErrNotNull             = errorx.New("sql.not_null")
)

type openConnectContext struct {
	Driver           string `json:"driver"`
	ConnectionString string `json:"connection_string"`
}

func NewOpenConnectError(err error, driver, connectionString string) error {
	return ErrOpenConnect.
		SetError(err).
		SetData(openConnectContext{
			Driver:           driver,
			ConnectionString: connectionString,
		})
}

type prepareStatementContext struct {
	Operation string `json:"operation"`
	Error     error  `json:"error"`
}

func NewPrepareStatementError(err error, operation string) error {
	return ErrStatementPrepare.
		SetError(err).
		SetData(prepareStatementContext{
			Operation: operation,
			Error:     err,
		})
}

type executeQueryContext struct {
	Operation string `json:"operation"`
	Error     error  `json:"error"`
}

func NewExecuteQueryError(err error, operation string) error {
	return ErrStatementExecuteQuery.
		SetError(err).
		SetData(executeQueryContext{
			Operation: operation,
			Error:     err,
		})
}

type duplicateContext struct {
	Field      string `json:"field"`
	Value      string `json:"value"`
	Constraint string `json:"constraint"`
}

func NewDuplicateError(duplicateErr *duplicate.Error) error {
	return ErrDuplicate.SetData(duplicateContext{
		Field:      duplicateErr.Field,
		Value:      duplicateErr.Value,
		Constraint: duplicateErr.Constraint,
	})
}

type foreignKeyViolationContext struct {
	Details string `json:"details"`
}

func NewForeignKeyViolationError(details string) error {
	return ErrForeignKeyViolation.SetData(foreignKeyViolationContext{
		Details: details,
	})
}

type notNullContext struct {
	Column string `json:"column"`
}

func NewNotNullError(column string) error {
	return ErrNotNull.SetData(notNullContext{
		Column: column,
	})
}
