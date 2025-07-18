//nolint:gosec,gocyclo,funlen
package convert

import "strconv"

// IntFromBool convert bool to int.
//
// If x equals true returns 1, otherwise 0
func IntFromBool(x bool) int {
	var result int
	if x {
		result = 1
	}

	return result
}

// IntFromFloat32 convert float32 to any integer type
func IntFromFloat32[T Integer](x float32) T {
	return T(x)
}

// IntFromFloat64 convert float64 to any integer type
func IntFromFloat64[T Integer](x float64) T {
	return T(x)
}

// Int8FromBool convert bool to int8.
//
// If x equals true returns 1, otherwise 0
func Int8FromBool(x bool) int8 {
	return int8(IntFromBool(x))
}

// Int16FromBool convert bool to int16.
//
// If x equals true returns 1, otherwise 0
func Int16FromBool(x bool) int16 {
	return int16(IntFromBool(x))
}

// Int32FromBool convert bool to int32.
//
// If x equals true returns 1, otherwise 0
func Int32FromBool(x bool) int32 {
	return int32(IntFromBool(x))
}

// Int64FromBool convert bool to int64.
//
// If x equals true returns 1, otherwise 0
func Int64FromBool(x bool) int64 {
	return int64(IntFromBool(x))
}

// UintFromBool convert bool to uint.
//
// If x equals true returns 1, otherwise 0
func UintFromBool(x bool) uint {
	return uint(IntFromBool(x))
}

// Uint8FromBool convert bool to uint8.
//
// If x equals true returns 1, otherwise 0
func Uint8FromBool(x bool) uint8 {
	return uint8(IntFromBool(x))
}

// Uint16FromBool convert bool to uint16.
//
// If x equals true returns 1, otherwise 0
func Uint16FromBool(x bool) uint16 {
	return uint16(IntFromBool(x))
}

// Uint32FromBool convert bool to uint32.
//
// If x equals true returns 1, otherwise 0
func Uint32FromBool(x bool) uint32 {
	return uint32(IntFromBool(x))
}

// Uint64FromBool convert bool to uint64.
//
// If x equals true returns 1, otherwise 0
func Uint64FromBool(x bool) uint64 {
	return uint64(IntFromBool(x))
}

// IntFromString convert string to int
func IntFromString(x string) int {
	if x == "" {
		return 0
	}

	result, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return 0
	}

	return int(result)
}

// Int8FromString convert string to int8
func Int8FromString(x string) int8 {
	return int8(IntFromString(x))
}

// Int16FromString convert string to int16
func Int16FromString(x string) int16 {
	return int16(IntFromString(x))
}

// Int32FromString convert string to int32
func Int32FromString(x string) int32 {
	return int32(IntFromString(x))
}

// Int64FromString convert string to int64
func Int64FromString(x string) int64 {
	if x == "" {
		return 0
	}

	result, err := strconv.ParseInt(x, 10, 64)
	if err != nil {
		return 0
	}

	return result
}

// UintFromString convert string to uint
func UintFromString(x string) uint {
	return uint(IntFromString(x))
}

// Uint8FromString convert string to uint8
func Uint8FromString(x string) uint8 {
	return uint8(IntFromString(x))
}

// Uint16FromString convert string to uint16
func Uint16FromString(x string) uint16 {
	return uint16(IntFromString(x))
}

// Uint32FromString convert string to uint32
func Uint32FromString(x string) uint32 {
	return uint32(IntFromString(x))
}

// Uint64FromString convert string to uint64
func Uint64FromString(x string) uint64 {
	return uint64(Int64FromString(x))
}

// Int convert any value to int.
//
// If x is nil returns 0
func Int(x any) int {
	switch v := x.(type) {
	case bool:
		return IntFromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return IntFromBool(*v)
	case string:
		return IntFromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return IntFromString(*v)
	case int:
		return v
	case *int:
		if v == nil {
			return 0
		}

		return *v
	case int8:
		return int(v)
	case *int8:
		if v == nil {
			return 0
		}

		return int(*v)
	case int16:
		return int(v)
	case *int16:
		if v == nil {
			return 0
		}

		return int(*v)
	case int32:
		return int(v)
	case *int32:
		if v == nil {
			return 0
		}

		return int(*v)
	case int64:
		return int(v)
	case *int64:
		if v == nil {
			return 0
		}

		return int(*v)
	case uint:
		return int(v)
	case *uint:
		if v == nil {
			return 0
		}

		return int(*v)
	case uint8:
		return int(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return int(*v)
	case uint16:
		return int(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return int(*v)
	case uint32:
		return int(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return int(*v)
	case uint64:
		return int(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return int(*v)
	case float32:
		return int(v)
	case *float32:
		if v == nil {
			return 0
		}

		return int(*v)
	case float64:
		return int(v)
	case *float64:
		if v == nil {
			return 0
		}

		return int(*v)
	default:
		return 0
	}
}

// Int8 convert any value to int8.
//
// If x is nil returns 0
func Int8(x any) int8 {
	switch v := x.(type) {
	case bool:
		return Int8FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Int8FromBool(*v)
	case string:
		return Int8FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Int8FromString(*v)
	case int8:
		return v
	case *int8:
		if v == nil {
			return 0
		}

		return *v
	case int:
		return int8(v)
	case *int:
		if v == nil {
			return 0
		}

		return int8(*v)
	case int16:
		return int8(v)
	case *int16:
		if v == nil {
			return 0
		}

		return int8(*v)
	case int32:
		return int8(v)
	case *int32:
		if v == nil {
			return 0
		}

		return int8(*v)
	case int64:
		return int8(v)
	case *int64:
		if v == nil {
			return 0
		}

		return int8(*v)
	case uint:
		return int8(v)
	case *uint:
		if v == nil {
			return 0
		}

		return int8(*v)
	case uint8:
		return int8(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return int8(*v)
	case uint16:
		return int8(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return int8(*v)
	case uint32:
		return int8(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return int8(*v)
	case uint64:
		return int8(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return int8(*v)
	case float32:
		return int8(v)
	case *float32:
		if v == nil {
			return 0
		}

		return int8(*v)
	case float64:
		return int8(v)
	case *float64:
		if v == nil {
			return 0
		}

		return int8(*v)
	default:
		return 0
	}
}

// Int16 convert any value to int16.
//
// If x is nil returns 0
func Int16(x any) int16 {
	switch v := x.(type) {
	case bool:
		return Int16FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Int16FromBool(*v)
	case string:
		return Int16FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Int16FromString(*v)
	case int16:
		return v
	case *int16:
		if v == nil {
			return 0
		}

		return *v
	case int8:
		return int16(v)
	case *int8:
		if v == nil {
			return 0
		}

		return int16(*v)
	case int:
		return int16(v)
	case *int:
		if v == nil {
			return 0
		}

		return int16(*v)
	case int32:
		return int16(v)
	case *int32:
		if v == nil {
			return 0
		}

		return int16(*v)
	case int64:
		return int16(v)
	case *int64:
		if v == nil {
			return 0
		}

		return int16(*v)
	case uint:
		return int16(v)
	case *uint:
		if v == nil {
			return 0
		}

		return int16(*v)
	case uint8:
		return int16(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return int16(*v)
	case uint16:
		return int16(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return int16(*v)
	case uint32:
		return int16(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return int16(*v)
	case uint64:
		return int16(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return int16(*v)
	case float32:
		return int16(v)
	case *float32:
		if v == nil {
			return 0
		}

		return int16(*v)
	case float64:
		return int16(v)
	case *float64:
		if v == nil {
			return 0
		}

		return int16(*v)
	default:
		return 0
	}
}

// Int32 convert any value to int32.
//
// If x is nil returns 0
func Int32(x any) int32 {
	switch v := x.(type) {
	case bool:
		return Int32FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Int32FromBool(*v)
	case string:
		return Int32FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Int32FromString(*v)
	case int8:
		return int32(v)
	case *int8:
		if v == nil {
			return 0
		}

		return int32(*v)
	case int16:
		return int32(v)
	case *int16:
		if v == nil {
			return 0
		}

		return int32(*v)
	case int32:
		return v
	case *int32:
		if v == nil {
			return 0
		}

		return *v
	case int64:
		return int32(v)
	case *int64:
		if v == nil {
			return 0
		}

		return int32(*v)
	case uint:
		return int32(v)
	case *uint:
		if v == nil {
			return 0
		}

		return int32(*v)
	case uint8:
		return int32(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return int32(*v)
	case uint16:
		return int32(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return int32(*v)
	case uint32:
		return int32(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return int32(*v)
	case uint64:
		return int32(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return int32(*v)
	case float32:
		return int32(v)
	case *float32:
		if v == nil {
			return 0
		}

		return int32(*v)
	case float64:
		return int32(v)
	case *float64:
		if v == nil {
			return 0
		}

		return int32(*v)
	default:
		return 0
	}
}

// Int64 convert any value to int64.
//
// If x is nil returns 0
func Int64(x any) int64 {
	switch v := x.(type) {
	case bool:
		return Int64FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Int64FromBool(*v)
	case string:
		return Int64FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Int64FromString(*v)
	case int8:
		return int64(v)
	case *int8:
		if v == nil {
			return 0
		}

		return int64(*v)
	case int16:
		return int64(v)
	case *int16:
		if v == nil {
			return 0
		}

		return int64(*v)
	case int32:
		return int64(v)
	case *int32:
		if v == nil {
			return 0
		}

		return int64(*v)
	case int64:
		return v
	case *int64:
		if v == nil {
			return 0
		}

		return *v
	case uint:
		return int64(v)
	case *uint:
		if v == nil {
			return 0
		}

		return int64(*v)
	case uint8:
		return int64(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return int64(*v)
	case uint16:
		return int64(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return int64(*v)
	case uint32:
		return int64(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return int64(*v)
	case uint64:
		return int64(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return int64(*v)
	case float32:
		return int64(v)
	case *float32:
		if v == nil {
			return 0
		}

		return int64(*v)
	case float64:
		return int64(v)
	case *float64:
		if v == nil {
			return 0
		}

		return int64(*v)
	default:
		return 0
	}
}

// Uint convert any value to uint.
//
// If x is nil returns 0
func Uint(x any) uint {
	switch v := x.(type) {
	case bool:
		return UintFromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return UintFromBool(*v)
	case string:
		return UintFromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return UintFromString(*v)
	case int8:
		return uint(v)
	case *int8:
		if v == nil {
			return 0
		}

		return uint(*v)
	case int16:
		return uint(v)
	case *int16:
		if v == nil {
			return 0
		}

		return uint(*v)
	case int32:
		return uint(v)
	case *int32:
		if v == nil {
			return 0
		}

		return uint(*v)
	case int64:
		return uint(v)
	case *int64:
		if v == nil {
			return 0
		}

		return uint(*v)
	case uint:
		return v
	case *uint:
		if v == nil {
			return 0
		}

		return *v
	case uint8:
		return uint(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return uint(*v)
	case uint16:
		return uint(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return uint(*v)
	case uint32:
		return uint(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return uint(*v)
	case uint64:
		return uint(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return uint(*v)
	case float32:
		return uint(v)
	case *float32:
		if v == nil {
			return 0
		}

		return uint(*v)
	case float64:
		return uint(v)
	case *float64:
		if v == nil {
			return 0
		}

		return uint(*v)
	default:
		return 0
	}
}

// Uint8 convert any value to uint8.
//
// If x is nil returns 0
func Uint8(x any) uint8 {
	switch v := x.(type) {
	case bool:
		return Uint8FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Uint8FromBool(*v)
	case string:
		return Uint8FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Uint8FromString(*v)
	case int8:
		return uint8(v)
	case *int8:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case int16:
		return uint8(v)
	case *int16:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case int32:
		return uint8(v)
	case *int32:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case int64:
		return uint8(v)
	case *int64:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case uint:
		return uint8(v)
	case *uint:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case uint8:
		return v
	case *uint8:
		if v == nil {
			return 0
		}

		return *v
	case uint16:
		return uint8(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case uint32:
		return uint8(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case uint64:
		return uint8(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case float32:
		return uint8(v)
	case *float32:
		if v == nil {
			return 0
		}

		return uint8(*v)
	case float64:
		return uint8(v)
	case *float64:
		if v == nil {
			return 0
		}

		return uint8(*v)
	default:
		return 0
	}
}

// Uint16 convert any value to uint16.
//
// If x is nil returns 0
func Uint16(x any) uint16 {
	switch v := x.(type) {
	case bool:
		return Uint16FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Uint16FromBool(*v)
	case string:
		return Uint16FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Uint16FromString(*v)
	case int8:
		return uint16(v)
	case *int8:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case int16:
		return uint16(v)
	case *int16:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case int32:
		return uint16(v)
	case *int32:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case int64:
		return uint16(v)
	case *int64:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case uint:
		return uint16(v)
	case *uint:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case uint8:
		return uint16(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case uint16:
		return v
	case *uint16:
		if v == nil {
			return 0
		}

		return *v
	case uint32:
		return uint16(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case uint64:
		return uint16(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case float32:
		return uint16(v)
	case *float32:
		if v == nil {
			return 0
		}

		return uint16(*v)
	case float64:
		return uint16(v)
	case *float64:
		if v == nil {
			return 0
		}

		return uint16(*v)
	default:
		return 0
	}
}

// Uint32 convert any value to uint32.
//
// If x is nil returns 0
func Uint32(x any) uint32 {
	switch v := x.(type) {
	case bool:
		return Uint32FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Uint32FromBool(*v)
	case string:
		return Uint32FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Uint32FromString(*v)
	case int8:
		return uint32(v)
	case *int8:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case int16:
		return uint32(v)
	case *int16:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case int32:
		return uint32(v)
	case *int32:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case int64:
		return uint32(v)
	case *int64:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case uint:
		return uint32(v)
	case *uint:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case uint8:
		return uint32(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case uint16:
		return uint32(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case uint32:
		return v
	case *uint32:
		if v == nil {
			return 0
		}

		return *v
	case uint64:
		return uint32(v)
	case *uint64:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case float32:
		return uint32(v)
	case *float32:
		if v == nil {
			return 0
		}

		return uint32(*v)
	case float64:
		return uint32(v)
	case *float64:
		if v == nil {
			return 0
		}

		return uint32(*v)
	default:
		return 0
	}
}

// Uint64 convert any value to uint64.
//
// If x is nil returns 0
func Uint64(x any) uint64 {
	switch v := x.(type) {
	case bool:
		return Uint64FromBool(v)
	case *bool:
		if v == nil {
			return 0
		}

		return Uint64FromBool(*v)
	case string:
		return Uint64FromString(v)
	case *string:
		if v == nil {
			return 0
		}

		return Uint64FromString(*v)
	case int8:
		return uint64(v)
	case *int8:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case int16:
		return uint64(v)
	case *int16:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case int32:
		return uint64(v)
	case *int32:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case int64:
		return uint64(v)
	case *int64:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case uint:
		return uint64(v)
	case *uint:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case uint8:
		return uint64(v)
	case *uint8:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case uint16:
		return uint64(v)
	case *uint16:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case uint32:
		return uint64(v)
	case *uint32:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case uint64:
		return v
	case *uint64:
		if v == nil {
			return 0
		}

		return *v
	case float32:
		return uint64(v)
	case *float32:
		if v == nil {
			return 0
		}

		return uint64(*v)
	case float64:
		return uint64(v)
	case *float64:
		if v == nil {
			return 0
		}

		return uint64(*v)
	default:
		return 0
	}
}
