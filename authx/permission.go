package authx

import "slices"

type Permission string

func (p Permission) String() string {
	return string(p)
}

func (p Permission) Is(permissions ...Permission) bool {
	return slices.Contains(permissions, p)
}
