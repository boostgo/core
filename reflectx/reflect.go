// Package reflectx contains reflect package helper tools.
package reflectx

import (
	"reflect"
)

func IsPointer(x any) bool {
	if x == nil {
		return true
	}

	return reflect.TypeOf(x).Kind() == reflect.Ptr
}

func CheckExport(x any) error {
	if !IsPointer(x) {
		return ErrCheckExport
	}

	return nil
}
