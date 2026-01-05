package tqtemplate

import (
	"reflect"
	"strings"
)

// undefinedSentinel is a sentinel type to distinguish between nil and undefined
type undefinedSentinel struct{}

var undefinedValue = &undefinedSentinel{}

// getBuiltinTests returns all builtin tests for the template engine
func getBuiltinTests() map[string]any {
	return map[string]any{
		"defined":     testDefined,
		"undefined":   testUndefined,
		"divisibleby": testDivisibleBy,
		"even":        testEven,
		"odd":         testOdd,
		"iterable":    testIterable,
		"null":        testNull,
		"number":      testNumber,
		"string":      testString,
		"__istest__":  filterIsTest,
		"__isnot__":   filterIsNot,
	}
}

// filterIsTest applies a test function to a value
func filterIsTest(value any, args ...any) bool {
	if len(args) == 0 {
		return false
	}

	// First arg is the test name
	testName := toString(args[0])

	// Get the test function
	tests := getBuiltinTests()
	testFn, exists := tests[testName]
	if !exists {
		return false
	}

	// Call the test with remaining args
	remaining := args[1:]

	// Type switch on test function signature
	switch fn := testFn.(type) {
	case func(any) bool:
		return fn(value)
	case func(any, ...any) bool:
		return fn(value, remaining...)
	default:
		return false
	}
}

// filterIsNot is the negation of filterIsTest
func filterIsNot(value any, args ...any) bool {
	return !filterIsTest(value, args...)
}

// testDefined returns true if the value is not undefined (even if it's nil)
func testDefined(value any) bool {
	// Check if it's the undefined sentinel
	if _, isUndefined := value.(*undefinedSentinel); isUndefined {
		return false
	}
	// nil is still defined if the variable exists in the data
	return true
}

// testUndefined returns true if the value is undefined (not just nil)
func testUndefined(value any) bool {
	// Check if it's the undefined sentinel
	_, isUndefined := value.(*undefinedSentinel)
	return isUndefined
}

// testDivisibleBy checks if a number is divisible by another number
func testDivisibleBy(value any, args ...any) bool {
	if len(args) == 0 {
		return false
	}

	num, ok := toNumber(value)
	if !ok {
		return false
	}

	divisor, ok := toNumber(args[0])
	if !ok || divisor == 0 {
		return false
	}

	return int(num)%int(divisor) == 0
}

// testEven returns true if the value is an even number
func testEven(value any) bool {
	num, ok := toNumber(value)
	if !ok {
		return false
	}
	return int(num)%2 == 0
}

// testOdd returns true if the value is an odd number
func testOdd(value any) bool {
	num, ok := toNumber(value)
	if !ok {
		return false
	}
	return int(num)%2 != 0
}

// testIterable checks if it's possible to iterate over an object
func testIterable(value any) bool {
	if value == nil {
		return false
	}

	v := reflect.ValueOf(value)
	kind := v.Kind()

	return kind == reflect.Slice ||
		kind == reflect.Array ||
		kind == reflect.Map ||
		kind == reflect.String
}

// testNull returns true if the value is nil
func testNull(value any) bool {
	return value == nil
}

// testNumber returns true if the value is a number
func testNumber(value any) bool {
	_, ok := toNumber(value)
	return ok
}

// testString returns true if the value is a string
func testString(value any) bool {
	_, ok := value.(string)
	return ok
}

// processIsTests preprocesses an expression to convert "is test" syntax to filter calls
// Returns the modified expression and extracted test filter (if any)
func processIsTests(expr string) (string, string) {
	// Pattern: value is testname or value is testname(args) or value is not testname
	// This converts "variable is defined" to: variable with filter __istest__("defined")

	// Look for " is " pattern
	isIdx := strings.Index(expr, " is ")
	if isIdx == -1 {
		return expr, ""
	}

	left := strings.TrimSpace(expr[:isIdx])
	right := strings.TrimSpace(expr[isIdx+4:])

	// Check for negation: "is not testname"
	isNegated := false
	if strings.HasPrefix(right, "not ") {
		isNegated = true
		right = strings.TrimSpace(right[4:])
	}

	// Extract test name and arguments
	testFilter := ""
	if strings.Contains(right, "(") {
		// Has arguments like "divisibleby(3)"
		parenIdx := strings.Index(right, "(")
		testName := right[:parenIdx]
		argsWithParens := right[parenIdx:] // includes parentheses and content
		argsInner := strings.TrimSuffix(strings.TrimPrefix(argsWithParens, "("), ")")

		if isNegated {
			testFilter = "__isnot__(\"" + testName + "\", " + argsInner + ")"
		} else {
			testFilter = "__istest__(\"" + testName + "\", " + argsInner + ")"
		}
	} else {
		// Simple test like "defined"
		if isNegated {
			testFilter = "__isnot__(\"" + right + "\")"
		} else {
			testFilter = "__istest__(\"" + right + "\")"
		}
	}

	return left, testFilter
}
