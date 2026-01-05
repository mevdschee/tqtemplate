package tqtemplate

import (
	"strings"
	"testing"
)

// String filter tests

func TestFilterLower(t *testing.T) {
	result, _ := template.Render("{{ text|lower }}", map[string]any{"text": "HELLO WORLD"})
	expected := "hello world"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterUpper(t *testing.T) {
	result, _ := template.Render("{{ text|upper }}", map[string]any{"text": "hello world"})
	expected := "HELLO WORLD"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterCapitalize(t *testing.T) {
	result, _ := template.Render("{{ text|capitalize }}", map[string]any{"text": "hello world"})
	expected := "Hello world"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTitle(t *testing.T) {
	result, _ := template.Render("{{ text|title }}", map[string]any{"text": "hello world"})
	expected := "Hello World"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTrim(t *testing.T) {
	result, _ := template.Render("{{ text|trim }}", map[string]any{"text": "  hello  "})
	expected := "hello"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTruncate(t *testing.T) {
	result, _ := template.Render("{{ text|truncate(10) }}", map[string]any{"text": "Hello World"})
	expected := "Hello..."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTruncateCustomEnd(t *testing.T) {
	result, _ := template.Render(`{{ text|truncate(10, ">>")|raw }}`, map[string]any{"text": "Hello World"})
	expected := "Hello>>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTruncateNoTruncation(t *testing.T) {
	result, _ := template.Render("{{ text|truncate(20) }}", map[string]any{"text": "Hello"})
	expected := "Hello"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTruncateLongWord(t *testing.T) {
	result, _ := template.Render("{{ text|truncate(10) }}", map[string]any{"text": "Supercalifragilistic"})
	expected := "Superca..."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTruncateMultipleWords(t *testing.T) {
	result, _ := template.Render("{{ text|truncate(20) }}", map[string]any{"text": "The quick brown fox jumps"})
	expected := "The quick brown..."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterTruncateWithTrailingSpace(t *testing.T) {
	result, _ := template.Render("{{ text|truncate(15) }}", map[string]any{"text": "Hello world and more"})
	expected := "Hello world..."
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterReplace(t *testing.T) {
	result, _ := template.Render(`{{ text|replace("Hello", "Goodbye") }}`, map[string]any{"text": "Hello World"})
	expected := "Goodbye World"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterReplaceWithCount(t *testing.T) {
	result, _ := template.Render(`{{ text|replace("a", "o", 2) }}`, map[string]any{"text": "banana"})
	expected := "bonona"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterSplit(t *testing.T) {
	result, _ := template.Render(`{{ text|split(",")|join("|") }}`, map[string]any{"text": "1,2,3"})
	expected := "1|2|3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterSplitChars(t *testing.T) {
	result, _ := template.Render(`{{ text|split|join("|") }}`, map[string]any{"text": "abc"})
	expected := "a|b|c"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterURLEncode(t *testing.T) {
	result, _ := template.Render("{{ text|urlencode }}", map[string]any{"text": "hello world"})
	expected := "hello+world"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterURLEncodeSpecialChars(t *testing.T) {
	result, _ := template.Render("{{ text|urlencode }}", map[string]any{"text": "hello&world=test"})
	expected := "hello%26world%3Dtest"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Numeric filter tests

func TestFilterAbs(t *testing.T) {
	result, _ := template.Render("{{ num|abs }}", map[string]any{"num": -42})
	expected := "42"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterAbsPositive(t *testing.T) {
	result, _ := template.Render("{{ num|abs }}", map[string]any{"num": 42})
	expected := "42"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterRound(t *testing.T) {
	result, _ := template.Render("{{ num|round }}", map[string]any{"num": 42.55})
	expected := "43"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterRoundWithPrecision(t *testing.T) {
	result, _ := template.Render(`{{ num|round(1, "floor") }}`, map[string]any{"num": 42.55})
	expected := "42.5"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterRoundCeil(t *testing.T) {
	result, _ := template.Render(`{{ num|round(0, "ceil") }}`, map[string]any{"num": 42.1})
	expected := "43"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterSprintf(t *testing.T) {
	result, _ := template.Render(`{{ num|sprintf("%.2f") }}`, map[string]any{"num": 3.14159})
	expected := "3.14"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterFileSizeFormat(t *testing.T) {
	result, _ := template.Render("{{ size|filesizeformat }}", map[string]any{"size": 13000})
	expected := "13.0 kB"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterFileSizeFormatBinary(t *testing.T) {
	result, _ := template.Render("{{ size|filesizeformat(true) }}", map[string]any{"size": 1024})
	expected := "1.0 KiB"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterFileSizeFormatLarge(t *testing.T) {
	result, _ := template.Render("{{ size|filesizeformat }}", map[string]any{"size": 1500000})
	expected := "1.5 MB"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Array/Collection filter tests

func TestFilterLength(t *testing.T) {
	result, _ := template.Render("{{ items|length }}", map[string]any{"items": []any{1, 2, 3}})
	expected := "3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterCount(t *testing.T) {
	result, _ := template.Render("{{ items|count }}", map[string]any{"items": []any{1, 2, 3, 4}})
	expected := "4"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterLengthString(t *testing.T) {
	result, _ := template.Render("{{ text|length }}", map[string]any{"text": "hello"})
	expected := "5"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterFirst(t *testing.T) {
	result, _ := template.Render("{{ items|first }}", map[string]any{"items": []any{1, 2, 3, 4}})
	expected := "1"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterFirstMultiple(t *testing.T) {
	result, _ := template.Render("{{ items|first(2)|join(\",\") }}", map[string]any{"items": []any{1, 2, 3, 4}})
	expected := "1,2"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterLast(t *testing.T) {
	result, _ := template.Render("{{ items|last }}", map[string]any{"items": []any{1, 2, 3, 4}})
	expected := "4"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterLastMultiple(t *testing.T) {
	result, _ := template.Render("{{ items|last(2)|join(\",\") }}", map[string]any{"items": []any{1, 2, 3, 4}})
	expected := "3,4"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterJoin(t *testing.T) {
	result, _ := template.Render(`{{ items|join("|") }}`, map[string]any{"items": []any{1, 2, 3}})
	expected := "1|2|3"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterJoinNoSeparator(t *testing.T) {
	result, _ := template.Render("{{ items|join }}", map[string]any{"items": []any{1, 2, 3}})
	expected := "123"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterJoinAttribute(t *testing.T) {
	users := []any{
		map[string]any{"name": "Alice"},
		map[string]any{"name": "Bob"},
	}
	result, _ := template.Render(`{{ users|join(", ", "name") }}`, map[string]any{"users": users})
	expected := "Alice, Bob"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterReverse(t *testing.T) {
	result, _ := template.Render("{{ items|reverse|join(\",\") }}", map[string]any{"items": []any{1, 2, 3}})
	expected := "3,2,1"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterReverseString(t *testing.T) {
	result, _ := template.Render("{{ text|reverse }}", map[string]any{"text": "hello"})
	expected := "olleh"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterSum(t *testing.T) {
	result, _ := template.Render("{{ items|sum }}", map[string]any{"items": []any{1, 2, 3}})
	expected := "6"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterSumAttribute(t *testing.T) {
	items := []any{
		map[string]any{"price": 10},
		map[string]any{"price": 20},
		map[string]any{"price": 30},
	}
	result, _ := template.Render(`{{ items|sum("price") }}`, map[string]any{"items": items})
	expected := "60"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Utility filter tests

func TestFilterDefault(t *testing.T) {
	// Use a nil value instead of missing to test default filter
	result, _ := template.Render(`{{ value|default("N/A") }}`, map[string]any{"value": nil})
	expected := "N/A"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterDefaultWithValue(t *testing.T) {
	result, _ := template.Render(`{{ value|default("N/A") }}`, map[string]any{"value": "exists"})
	expected := "exists"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterDefaultBoolean(t *testing.T) {
	result, _ := template.Render(`{{ value|default("empty", true) }}`, map[string]any{"value": ""})
	expected := "empty"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterDefaultBooleanZero(t *testing.T) {
	result, _ := template.Render(`{{ value|default("zero", true) }}`, map[string]any{"value": 0})
	expected := "zero"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterAttr(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name":  "Alice",
			"email": "alice@example.com",
		},
	}
	result, _ := template.Render(`{{ user|attr("email") }}`, data)
	expected := "alice@example.com"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterAttrMissing(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
		},
	}
	result, _ := template.Render(`{{ user|attr("missing") }}`, data)
	expected := ""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterDebug(t *testing.T) {
	data := map[string]any{
		"user": map[string]any{
			"name": "Alice",
			"age":  30,
		},
	}
	result, _ := template.Render("{{ user|debug|raw }}", data)
	// Should contain JSON formatted output
	if !strings.Contains(result, `"name"`) || !strings.Contains(result, `"Alice"`) {
		t.Errorf("Expected JSON output, got '%s'", result)
	}
}

func TestFilterDebugAlias(t *testing.T) {
	data := map[string]any{"value": 42}
	result, _ := template.Render("{{ value|d }}", data)
	expected := "42"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterRaw(t *testing.T) {
	result, _ := template.Render("{{ html|raw }}", map[string]any{"html": "<strong>Bold</strong>"})
	expected := "<strong>Bold</strong>"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterRawWithoutEscaping(t *testing.T) {
	result, _ := template.Render("{{ html }}", map[string]any{"html": "<strong>Bold</strong>"})
	expected := "&lt;strong&gt;Bold&lt;/strong&gt;"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Filter chaining tests

func TestFilterChaining(t *testing.T) {
	result, _ := template.Render(`{{ text|trim|upper|replace("WORLD", "FRIEND") }}`, map[string]any{"text": "  hello world  "})
	expected := "HELLO FRIEND"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterChainingArrays(t *testing.T) {
	result, _ := template.Render(`{{ items|first(3)|reverse|join(", ") }}`, map[string]any{"items": []any{1, 2, 3, 4, 5}})
	expected := "3, 2, 1"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterChainingComplex(t *testing.T) {
	users := []any{
		map[string]any{"name": "alice"},
		map[string]any{"name": "bob"},
		map[string]any{"name": "charlie"},
	}
	result, _ := template.Render(`{{ users|join(", ", "name")|upper }}`, map[string]any{"users": users})
	expected := "ALICE, BOB, CHARLIE"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Edge case tests

func TestFilterEmptyArray(t *testing.T) {
	result, _ := template.Render("{{ items|length }}", map[string]any{"items": []any{}})
	expected := "0"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterEmptyString(t *testing.T) {
	result, _ := template.Render("{{ text|upper }}", map[string]any{"text": ""})
	expected := ""
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterNilValue(t *testing.T) {
	result, _ := template.Render(`{{ value|default("nil") }}`, map[string]any{"value": nil})
	expected := "nil"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestFilterNumericString(t *testing.T) {
	result, _ := template.Render("{{ num|abs }}", map[string]any{"num": "-42"})
	expected := "42"
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}
