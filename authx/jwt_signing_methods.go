package authx

type JwtSignMethod string

func (m JwtSignMethod) String() string {
	return string(m)
}

const (
	JwtSignMethodHS256 JwtSignMethod = "HS256"
	JwtSignMethodRS256 JwtSignMethod = "RS256"
)
