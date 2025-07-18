package translate

import "github.com/boostgo/core/errorx"

var (
	ErrReadFile       = errorx.New("translate.read_file")
	ErrLocaleNotFound = errorx.New("translate.locale_not_found")
	ErrKeyNotFound    = errorx.New("translate.key_not_found")
)
