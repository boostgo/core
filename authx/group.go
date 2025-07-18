package authx

import (
	"slices"
	"strings"
)

type Group string

func (g Group) String() string {
	return string(g)
}

func (g Group) NoSuffix() string {
	return strings.Split(g.String(), "-")[0]
}

func (g Group) Is(groups ...Group) bool {
	return slices.Contains(groups, g)
}

func (g Group) IsString(groups ...string) bool {
	return slices.Contains(groups, g.String())
}
