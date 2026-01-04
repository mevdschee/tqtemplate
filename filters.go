package tqtemplate

import (
	"encoding/json"
	"fmt"
	"math"
	"net/url"
	"reflect"
	"strings"
)

// getBuiltinFilters returns all builtin filters for the template engine
func getBuiltinFilters() map[string]any {
	return map[string]any{
		"abs":            filterAbs,
		"attr":           filterAttr,
		"capitalize":     filterCapitalize,
		"default":        filterDefault,
		"filesizeformat": filterFileSizeFormat,
		"first":          filterFirst,
		"sprintf":        filterSprintf,
		"join":           filterJoin,
		"split":          filterSplit,
		"last":           filterLast,
		"length":         filterLength,
		"count":          filterLength, // alias for length
		"lower":          filterLower,
		"debug":          filterDebug,
		"d":              filterDebug, // alias for debug
		"replace":        filterReplace,
		"reverse":        filterReverse,
		"round":          filterRound,
		"sum":            filterSum,
		"title":          filterTitle,
		"trim":           filterTrim,
		"truncate":       filterTruncate,
		"upper":          filterUpper,
		"urlencode":      filterURLEncode,
		"raw":            filterRaw,
	}
}

// filterRaw marks a value that should not be escaped
func filterRaw(value any) RawValue {
	return RawValue{Value: toString(value)}
}

// filterAbs returns the absolute value of a number
func filterAbs(value any) any {
	num, ok := toNumber(value)
	if !ok {
		return value
	}
	return math.Abs(num)
}

// filterAttr gets an attribute of an object by name
func filterAttr(obj any, name any) any {
	attrName := toString(name)

	// Handle map
	if m, ok := obj.(map[string]any); ok {
		if val, exists := m[attrName]; exists {
			return val
		}
		return nil
	}

	// Handle struct using reflection
	v := reflect.ValueOf(obj)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() == reflect.Struct {
		field := v.FieldByName(attrName)
		if field.IsValid() {
			return field.Interface()
		}
	}

	return nil
}

// filterCapitalize capitalizes the first character of a string
func filterCapitalize(value any) string {
	s := toString(value)
	if s == "" {
		return s
	}
	return strings.ToUpper(s[:1]) + s[1:]
}

// filterDefault returns a default value if the input is nil or empty
func filterDefault(value any, args ...any) any {
	defaultValue := ""
	boolean := false

	if len(args) > 0 {
		defaultValue = toString(args[0])
	}
	if len(args) > 1 {
		boolean = toBool(args[1])
	}

	// If boolean mode, check truthiness
	if boolean {
		if !toBool(value) {
			return defaultValue
		}
		return value
	}

	// Default mode: only check for nil
	if value == nil {
		return defaultValue
	}

	return value
}

// filterFileSizeFormat formats a number as a human-readable file size
func filterFileSizeFormat(value any, args ...any) string {
	num, ok := toNumber(value)
	if !ok {
		return toString(value)
	}

	binary := false
	if len(args) > 0 {
		binary = toBool(args[0])
	}

	var base float64
	var units []string

	if binary {
		base = 1024
		units = []string{"Bytes", "KiB", "MiB", "GiB", "TiB", "PiB"}
	} else {
		base = 1000
		units = []string{"Bytes", "kB", "MB", "GB", "TB", "PB"}
	}

	if num < base {
		return fmt.Sprintf("%.0f %s", num, units[0])
	}

	exp := int(math.Log(num) / math.Log(base))
	if exp >= len(units) {
		exp = len(units) - 1
	}

	val := num / math.Pow(base, float64(exp))
	return fmt.Sprintf("%.1f %s", val, units[exp])
}

// filterFirst returns the first n items of a slice
func filterFirst(value any, args ...any) any {
	number := 1
	if len(args) > 0 {
		if num, ok := toNumber(args[0]); ok {
			number = int(num)
		}
	}

	// Convert to slice
	slice := toSlice(value)
	if slice == nil {
		return value
	}

	if number >= len(slice) {
		return slice
	}

	if number == 1 && len(slice) > 0 {
		return slice[0]
	}

	return slice[:number]
}

// filterSprintf applies sprintf-style formatting
func filterSprintf(value any, args ...any) string {
	if len(args) == 0 {
		return toString(value)
	}

	format := toString(args[0])
	return fmt.Sprintf(format, value)
}

// filterJoin concatenates items in a slice with a separator
func filterJoin(value any, args ...any) string {
	separator := ""
	attribute := ""

	if len(args) > 0 {
		separator = toString(args[0])
	}
	if len(args) > 1 {
		attribute = toString(args[1])
	}

	slice := toSlice(value)
	if slice == nil {
		return toString(value)
	}

	var parts []string
	for _, item := range slice {
		if attribute != "" {
			// Join by attribute
			item = filterAttr(item, attribute)
		}
		parts = append(parts, toString(item))
	}

	return strings.Join(parts, separator)
}

// filterSplit splits a string into a slice
func filterSplit(value any, args ...any) any {
	s := toString(value)
	separator := ""

	if len(args) > 0 {
		separator = toString(args[0])
	}

	if separator == "" {
		// Split into characters
		chars := []any{}
		for _, char := range s {
			chars = append(chars, string(char))
		}
		return chars
	}

	parts := strings.Split(s, separator)
	result := []any{}
	for _, part := range parts {
		result = append(result, part)
	}
	return result
}

// filterLast returns the last n items of a slice
func filterLast(value any, args ...any) any {
	number := 1
	if len(args) > 0 {
		if num, ok := toNumber(args[0]); ok {
			number = int(num)
		}
	}

	slice := toSlice(value)
	if slice == nil {
		return value
	}

	if number >= len(slice) {
		return slice
	}

	if number == 1 && len(slice) > 0 {
		return slice[len(slice)-1]
	}

	return slice[len(slice)-number:]
}

// filterLength returns the number of items in a collection
func filterLength(value any) int {
	slice := toSlice(value)
	if slice != nil {
		return len(slice)
	}

	// Handle string
	if s, ok := value.(string); ok {
		return len(s)
	}

	// Handle map
	if m, ok := value.(map[string]any); ok {
		return len(m)
	}

	return 0
}

// filterLower converts a string to lowercase
func filterLower(value any) string {
	return strings.ToLower(toString(value))
}

// filterDebug returns a debug representation of a value
func filterDebug(value any) string {
	bytes, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Sprintf("%#v", value)
	}
	return string(bytes)
}

// filterReplace replaces occurrences of a substring
func filterReplace(value any, args ...any) string {
	s := toString(value)

	if len(args) < 2 {
		return s
	}

	old := toString(args[0])
	new := toString(args[1])
	count := -1

	if len(args) > 2 {
		if num, ok := toNumber(args[2]); ok {
			count = int(num)
		}
	}

	if count < 0 {
		return strings.ReplaceAll(s, old, new)
	}

	return strings.Replace(s, old, new, count)
}

// filterReverse reverses a slice or string
func filterReverse(value any) any {
	// Handle slice
	slice := toSlice(value)
	if slice != nil {
		reversed := make([]any, len(slice))
		for i, v := range slice {
			reversed[len(slice)-1-i] = v
		}
		return reversed
	}

	// Handle string
	if s, ok := value.(string); ok {
		runes := []rune(s)
		for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
			runes[i], runes[j] = runes[j], runes[i]
		}
		return string(runes)
	}

	return value
}

// filterRound rounds a number to a given precision
func filterRound(value any, args ...any) float64 {
	num, ok := toNumber(value)
	if !ok {
		return 0
	}

	precision := 0
	method := "common"

	if len(args) > 0 {
		if p, ok := toNumber(args[0]); ok {
			precision = int(p)
		}
	}
	if len(args) > 1 {
		method = toString(args[1])
	}

	multiplier := math.Pow(10, float64(precision))
	scaled := num * multiplier

	var rounded float64
	switch method {
	case "ceil":
		rounded = math.Ceil(scaled)
	case "floor":
		rounded = math.Floor(scaled)
	case "up", "common":
		// Round half up
		if scaled >= 0 {
			rounded = math.Floor(scaled + 0.5)
		} else {
			rounded = math.Ceil(scaled - 0.5)
		}
	case "down":
		// Round half down
		if scaled >= 0 {
			rounded = math.Ceil(scaled - 0.5)
		} else {
			rounded = math.Floor(scaled + 0.5)
		}
	case "even", "banker":
		// Round half to even
		rounded = math.RoundToEven(scaled)
	case "odd":
		// Round half to odd
		floor := math.Floor(scaled)
		ceil := math.Ceil(scaled)
		diff := scaled - floor
		if diff == 0.5 {
			if int(floor)%2 == 0 {
				rounded = ceil
			} else {
				rounded = floor
			}
		} else if diff > 0.5 {
			rounded = ceil
		} else {
			rounded = floor
		}
	case "awayzero":
		// Round half away from zero
		if scaled >= 0 {
			rounded = math.Floor(scaled + 0.5)
		} else {
			rounded = math.Ceil(scaled - 0.5)
		}
	case "tozero":
		// Round half towards zero
		if scaled >= 0 {
			rounded = math.Ceil(scaled - 0.5)
		} else {
			rounded = math.Floor(scaled + 0.5)
		}
	default:
		// Default to common (up)
		if scaled >= 0 {
			rounded = math.Floor(scaled + 0.5)
		} else {
			rounded = math.Ceil(scaled - 0.5)
		}
	}

	return rounded / multiplier
}

// filterSum returns the sum of numbers in a slice
func filterSum(value any, args ...any) float64 {
	attribute := ""
	if len(args) > 0 {
		attribute = toString(args[0])
	}

	slice := toSlice(value)
	if slice == nil {
		return 0
	}

	sum := 0.0
	for _, item := range slice {
		if attribute != "" {
			item = filterAttr(item, attribute)
		}

		if num, ok := toNumber(item); ok {
			sum += num
		}
	}

	return sum
}

// filterTitle converts a string to title case
func filterTitle(value any) string {
	s := toString(value)
	return strings.Title(s)
}

// filterTrim removes leading and trailing whitespace
func filterTrim(value any) string {
	return strings.TrimSpace(toString(value))
}

// filterTruncate truncates a string to a maximum length
func filterTruncate(value any, args ...any) string {
	s := toString(value)

	length := 255
	end := "..."

	if len(args) > 0 {
		if l, ok := toNumber(args[0]); ok {
			length = int(l)
		}
	}
	if len(args) > 1 {
		end = toString(args[1])
	}

	if len(s) <= length {
		return s
	}

	truncated := s[:length-len(end)]
	return truncated + end
}

// filterUpper converts a string to uppercase
func filterUpper(value any) string {
	return strings.ToUpper(toString(value))
}

// filterURLEncode encodes a string for use in URLs
func filterURLEncode(value any) string {
	// Handle map
	if m, ok := value.(map[string]any); ok {
		values := url.Values{}
		for k, v := range m {
			values.Add(k, toString(v))
		}
		return values.Encode()
	}

	// Handle string
	return url.QueryEscape(toString(value))
}

// toSlice converts a value to a slice of any
func toSlice(value any) []any {
	if value == nil {
		return nil
	}

	// Already a slice of any
	if slice, ok := value.([]any); ok {
		return slice
	}

	// Handle other slice types using reflection
	v := reflect.ValueOf(value)
	if v.Kind() == reflect.Slice || v.Kind() == reflect.Array {
		result := make([]any, v.Len())
		for i := 0; i < v.Len(); i++ {
			result[i] = v.Index(i).Interface()
		}
		return result
	}

	return nil
}
