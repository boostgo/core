//nolint:govet,gocritic
package convert

import (
	"encoding/json"
	"fmt"
	"reflect"
	"unsafe"
)

// BytesFromString converts string to bytes slice with no allocation
func BytesFromString(x string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{
		Data: (*(*reflect.StringHeader)(unsafe.Pointer(&x))).Data,
		Len:  len(x),
		Cap:  len(x),
	}))
}

// Bytes convert any value to bytes.
//
// If x is string calls BytesFromString function.
//
// If x is numeric convert it by json marshaller.
//
// If x is fmt.Stringer implementation calls .String() method and then to bytes.
func Bytes(x any) []byte {
	if x == nil {
		return nil
	}

	switch v := x.(type) {
	case []byte:
		return v
	case string:
		return BytesFromString(v)
	case *string:
		if v == nil {
			return nil
		}

		return BytesFromString(*v)
	case fmt.Stringer:
		return BytesFromString(v.String())
	default:
		blob, err := json.Marshal(x)
		if err != nil {
			return nil
		}

		return blob
	}
}
