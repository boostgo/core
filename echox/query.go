package echox

import (
	"github.com/boostgo/core/sorts"
)

type SortByParams struct {
	Field string `json:"field" query:"sort-by-field" form:"sort-by-field"`
	Asc   bool   `json:"asc" query:"sort-by-asc" form:"sort-by-asc"`
}

func (s SortByParams) SortBy() sorts.Params {
	return sorts.Params{
		Field: s.Field,
		Asc:   s.Asc,
	}
}
