package sorts

import (
	"fmt"
)

const (
	Asc  = "ASC"
	Desc = "DESC"
)

type Params struct {
	Field string
	Asc   bool
}

func (p Params) Empty() bool {
	return p.Field == ""
}

func (p Params) Query(alias ...string) string {
	var direction string
	if p.Asc {
		direction = Asc
	} else {
		direction = Desc
	}
	return wrapFieldWithAlias(p.Field, alias...) + " " + direction
}

func (p Params) Direction() string {
	if p.Asc {
		return Asc
	}
	return Desc
}

func wrapFieldWithAlias(field string, alias ...string) string {
	var setAlias string
	if len(alias) > 0 {
		setAlias = alias[0]
	}

	return fmt.Sprintf("%s%s", setAlias, field)
}
