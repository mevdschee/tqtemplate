package tqtemplate

import (
	"fmt"
	"strconv"
)

// toBool converts a value to boolean
func toBool(value any) bool {
	switch v := value.(type) {
	case bool:
		return v
	case int:
		return v != 0
	case float64:
		return v != 0
	case string:
		return v != ""
	case nil:
		return false
	default:
		return true
	}
}

// toNumber converts a value to float64
func toNumber(value any) (float64, bool) {
	switch v := value.(type) {
	case int:
		return float64(v), true
	case float64:
		return v, true
	case string:
		if f, err := strconv.ParseFloat(v, 64); err == nil {
			return f, true
		}
		return 0, false
	default:
		return 0, false
	}
}

// toString converts a value to string
func toString(value any) string {
	switch v := value.(type) {
	case string:
		return v
	case int:
		return strconv.Itoa(v)
	case float64:
		// Format number without unnecessary trailing zeros
		str := strconv.FormatFloat(v, 'f', -1, 64)
		return str
	case bool:
		if v {
			return "1"
		}
		return ""
	case nil:
		return ""
	default:
		return fmt.Sprintf("%v", v)
	}
}

// compare compares two values and returns -1, 0, or 1
func compare(left, right any) int {
	// Try numeric comparison first
	leftNum, leftIsNum := toNumber(left)
	rightNum, rightIsNum := toNumber(right)
	if leftIsNum && rightIsNum {
		if leftNum < rightNum {
			return -1
		} else if leftNum > rightNum {
			return 1
		}
		return 0
	}

	// Fall back to string comparison
	leftStr := toString(left)
	rightStr := toString(right)
	if leftStr < rightStr {
		return -1
	} else if leftStr > rightStr {
		return 1
	}
	return 0
}

// callFunction calls a function with the given arguments
func callFunction(fn any, args []any) (any, error) {
	switch f := fn.(type) {
	// RawValue functions
	case func(any) RawValue:
		if len(args) > 0 {
			return f(args[0]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(string) RawValue:
		if len(args) > 0 {
			if str, ok := args[0].(string); ok {
				return f(str), nil
			}
		}
		return nil, fmt.Errorf("invalid arguments for function")

	// String functions
	case func(any) string:
		if len(args) > 0 {
			return f(args[0]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(string) string:
		if len(args) > 0 {
			return f(toString(args[0])), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(string, string) string:
		if len(args) >= 2 {
			return f(toString(args[0]), toString(args[1])), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(any, ...any) string:
		if len(args) > 0 {
			return f(args[0], args[1:]...), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")

	// Int functions
	case func(any) int:
		if len(args) > 0 {
			return f(args[0]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")

	// Float64 functions
	case func(any) float64:
		if len(args) > 0 {
			return f(args[0]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(any, ...any) float64:
		if len(args) > 0 {
			return f(args[0], args[1:]...), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")

	// Any functions (generic)
	case func(any) any:
		if len(args) > 0 {
			return f(args[0]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(any, any) any:
		if len(args) >= 2 {
			return f(args[0], args[1]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(any, ...any) any:
		if len(args) > 0 {
			return f(args[0], args[1:]...), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")

	// Boolean functions
	case func(any) bool:
		if len(args) > 0 {
			return f(args[0]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(any, ...any) bool:
		if len(args) > 0 {
			return f(args[0], args[1:]...), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(any, any) bool:
		if len(args) >= 2 {
			return f(args[0], args[1]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(any, any, any) bool:
		if len(args) >= 3 {
			return f(args[0], args[1], args[2]), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")
	case func(int, int) bool:
		if len(args) >= 2 {
			left, _ := toNumber(args[0])
			right, _ := toNumber(args[1])
			return f(int(left), int(right)), nil
		}
		return nil, fmt.Errorf("invalid arguments for function")

	default:
		return nil, fmt.Errorf("unsupported function type")
	}
}
