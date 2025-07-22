// Package sql provide tools for manipulating connections and helper tools (sqlx extension).
// Features:
// - Client for single database and for sharding (common interface - DB).
// - More simple connecting.
// - Arguments tool. Helps to set arguments for multiple insert.
// - Connection builder. Building connection string or connecting.
// - Any driver support.
// - Migrations.
// - Transactor implementation. Implementation based on manipulating transaction from context.
package sql

type SortBy struct {
	Field string
	Asc   bool
}

// Page returns offset & limit by pagination
func Page(pageSize int64, page int) (offset int, limit int) {
	if page == 0 {
		page = 1
	}

	offset = (page - 1) * int(pageSize)
	limit = int(pageSize)
	return offset, limit
}
