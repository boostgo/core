package sorts

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

func (p Params) Query() string {
	var direction string
	if p.Asc {
		direction = Asc
	} else {
		direction = Desc
	}
	return p.Field + " " + direction
}

func (p Params) Direction() string {
	if p.Asc {
		return Asc
	}
	return Desc
}
