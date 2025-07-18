//nolint:gocyclo,funlen
package convert

// BoolFromString convert string to bool
func BoolFromString(x string) bool {
	return x == "true" || x == "TRUE"
}

// BoolFromInt convert any integer type to bool
func BoolFromInt[T Integer](x T) bool {
	return x == 1
}

// Bool convert any value to bool.
//
// If x is nil returns false.
//
// If x is string, value should be "true" or "TRUE"
func Bool(x any) bool {
	if x == nil {
		return false
	}

	switch v := x.(type) {
	case bool:
		return v
	case *bool:
		return *v
	case string:
		return BoolFromString(v)
	case *string:
		if v == nil {
			return false
		}

		return BoolFromString(*v)
	case int:
		return BoolFromInt(v)
	case *int:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case int8:
		return BoolFromInt(v)
	case *int8:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case int16:
		return BoolFromInt(v)
	case *int16:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case int32:
		return BoolFromInt(v)
	case *int32:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case int64:
		return BoolFromInt(v)
	case *int64:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case uint:
		return BoolFromInt(v)
	case *uint:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case uint8:
		return BoolFromInt(v)
	case *uint8:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case uint16:
		return BoolFromInt(v)
	case *uint16:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case uint32:
		return BoolFromInt(v)
	case *uint32:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	case uint64:
		return BoolFromInt(v)
	case *uint64:
		if v == nil {
			return false
		}

		return BoolFromInt(*v)
	default:
		return false
	}
}
