package echox

import "github.com/boostgo/core/sql"

type SortByParams struct {
	Field string `json:"field" query:"sort-by-field" form:"sort-by-field"`
	Asc   bool   `json:"asc" query:"sort-by-asc" form:"sort-by-asc"`
}

func (s SortByParams) SortBy() sql.SortBy {
	return sql.SortBy{
		Field: s.Field,
		Asc:   s.Asc,
	}
}
