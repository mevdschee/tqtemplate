package tqtemplate

import (
	"strings"
	"testing"
	"time"
)

var template *Template

func init() {
	template = NewTemplate()
}

func TestRenderWithCustomFunction(t *testing.T) {
	functions := map[string]any{
		"capitalize": func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
	}
	result, _ := template.Render("hello {{ name|capitalize }}", map[string]any{"name": "world"}, functions)
	if result != "hello World" {
		t.Errorf("Expected 'hello World', got '%s'", result)
	}
}

func TestRenderWithHtmlEscaping(t *testing.T) {
	result, _ := template.Render("<br>hello {{ name }}", map[string]any{"name": "<br>world"}, nil)
	if result != "<br>hello &lt;br&gt;world" {
		t.Errorf("Expected '<br>hello &lt;br&gt;world', got '%s'", result)
	}
}

func TestRenderWithMissingFunction(t *testing.T) {
	functions := map[string]any{
		"capitalize": func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
	}
	result, _ := template.Render("hello {{ name|failure }}", map[string]any{"name": "world"}, functions)
	if result != "hello {{name|failure!!function `failure` not found}}" {
		t.Errorf("Expected error message, got '%s'", result)
	}
}

func TestRenderIfWithNestedPath(t *testing.T) {
	functions := map[string]any{
		"eq": func(a, b any) bool { return a == b },
	}
	result, _ := template.Render("hello {% if n.m|eq(3) %}m is 3{% endif %}", map[string]any{
		"n": map[string]any{"m": 3},
	}, functions)
	if result != "hello m is 3" {
		t.Errorf("Expected 'hello m is 3', got '%s'", result)
	}
}

func TestRenderIfElse(t *testing.T) {
	result, _ := template.Render("hello {% if n %}n{% else %}not n{% endif %}", map[string]any{"n": false}, nil)
	if result != "hello not n" {
		t.Errorf("Expected 'hello not n', got '%s'", result)
	}
}

func TestRenderWithFunctionLiteralArgument(t *testing.T) {
	functions := map[string]any{
		"dateFormat": func(dateStr string, format string) string {
			t, _ := time.Parse("January 2, 2006", dateStr)
			return t.Format("2006-01-02")
		},
	}
	result, _ := template.Render("hello {{ name|dateFormat(\"Y-m-d\") }}", map[string]any{"name": "May 13, 1980"}, functions)
	if result != "hello 1980-05-13" {
		t.Errorf("Expected 'hello 1980-05-13', got '%s'", result)
	}
}

func TestRenderWithFunctionDataArgument(t *testing.T) {
	functions := map[string]any{
		"dateFormat": func(dateStr string, format string) string {
			t, _ := time.Parse("January 2, 2006", dateStr)
			return t.Format("2006-01-02")
		},
	}
	data := map[string]any{
		"name":   "May 13, 1980",
		"format": "Y-m-d",
	}
	result, _ := template.Render("hello {{ name|dateFormat(format) }}", data, functions)
	if result != "hello 1980-05-13" {
		t.Errorf("Expected 'hello 1980-05-13', got '%s'", result)
	}
}

func TestRenderWithFunctionComplexLiteralArgument(t *testing.T) {
	functions := map[string]any{
		"dateFormat": func(dateStr string, format string) string {
			t, _ := time.Parse("January 2, 2006", dateStr)
			return t.Format("Jan 2, 2006")
		},
	}
	result, _ := template.Render("hello {{ name|dateFormat(\"M j, Y\") }}", map[string]any{"name": "May 13, 1980"}, functions)
	if result != "hello May 13, 1980" {
		t.Errorf("Expected 'hello May 13, 1980', got '%s'", result)
	}
}

func TestRenderWithFunctionArgumentWithWhitespace(t *testing.T) {
	functions := map[string]any{
		"dateFormat": func(dateStr string, format string) string {
			t, _ := time.Parse("January 2, 2006", dateStr)
			return t.Format("Jan 2, 2006")
		},
	}
	result, _ := template.Render("hello {{ name|dateFormat( \"M j, Y\") }}", map[string]any{"name": "May 13, 1980"}, functions)
	if result != "hello May 13, 1980" {
		t.Errorf("Expected 'hello May 13, 1980', got '%s'", result)
	}
}

func TestRenderWithEscapedSpecialCharacters(t *testing.T) {
	functions := map[string]any{
		"dateFormat": func(dateStr string, format string) string {
			// Convert PHP date format to Go time format
			// " M ()}}\",|:.j, Y" becomes " Jan ()}}\",|:.2, 2006"
			goFormat := strings.ReplaceAll(format, "M", "Jan")
			goFormat = strings.ReplaceAll(goFormat, "j", "2")
			goFormat = strings.ReplaceAll(goFormat, "Y", "2006")

			t, _ := time.Parse("January 2, 2006", dateStr)
			return t.Format(goFormat)
		},
	}

	// Template with escaped quotes and special characters in the filter argument
	tmpl := "hello \"{{ name|dateFormat(\" M ()}}\\\",|:.j, Y\") }}\""
	expected := "hello \" May ()}}&#34;,|:.13, 1980\""

	result, _ := template.Render(tmpl, map[string]any{"name": "May 13, 1980"}, functions)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestRenderForLoopWithValues(t *testing.T) {
	result, _ := template.Render("test{% for i in counts %} {{ i }}{% endfor %}", map[string]any{
		"counts": []any{1, 2, 3},
	}, nil)
	if result != "test 1 2 3" {
		t.Errorf("Expected 'test 1 2 3', got '%s'", result)
	}
}

func TestRenderForLoopWithKeysAndValues(t *testing.T) {
	result, _ := template.Render("test{% for k, v in counts %} {{ k }}={{ v }}{% endfor %}", map[string]any{
		"counts": map[string]any{"a": 1, "b": 2, "c": 3},
	}, nil)
	// Note: map iteration order is not guaranteed in Go, so we need to check all possibilities
	validResults := []string{
		"test a=1 b=2 c=3",
		"test a=1 c=3 b=2",
		"test b=2 a=1 c=3",
		"test b=2 c=3 a=1",
		"test c=3 a=1 b=2",
		"test c=3 b=2 a=1",
	}
	found := false
	for _, valid := range validResults {
		if result == valid {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("Expected one of the valid results, got '%s'", result)
	}
}

func TestRenderNestedForLoops(t *testing.T) {
	result, _ := template.Render("test{% for x in steps %}{% for y in steps %} ({{ x }},{{ y }}){% endfor %}{% endfor %}", map[string]any{
		"steps": []any{-1, 1},
	}, nil)
	if result != "test (-1,-1) (-1,1) (1,-1) (1,1)" {
		t.Errorf("Expected 'test (-1,-1) (-1,1) (1,-1) (1,1)', got '%s'", result)
	}
}

func TestRenderForLoopWithIfElseIf(t *testing.T) {
	functions := map[string]any{
		"eq": func(a, b any) bool { return a == b },
	}
	result, _ := template.Render("hello{% for i in counts %} {% if i|eq(1) %}one{% elseif i|eq(2) %}two{% else %}three{% endif %}{% endfor %}", map[string]any{
		"counts": []any{1, 2, 3},
	}, functions)
	if result != "hello one two three" {
		t.Errorf("Expected 'hello one two three', got '%s'", result)
	}
}

func TestEscape(t *testing.T) {
	result, _ := template.Render("{{ a }}", map[string]any{"a": "<script>alert(\"xss\")</script>"}, nil)
	// Both &quot; and &#34; are valid HTML entities for quotes
	if result != "&lt;script&gt;alert(&#34;xss&#34;)&lt;/script&gt;" {
		t.Errorf("Expected escaped output, got '%s'", result)
	}
}

func TestRawEscape(t *testing.T) {
	result, _ := template.Render("{{ a|raw }}", map[string]any{"a": "<script>alert(\"xss\")</script>"}, nil)
	if result != "<script>alert(\"xss\")</script>" {
		t.Errorf("Expected raw output, got '%s'", result)
	}
}

// Expression tests - Basic comparison operators
func TestExpressionEquals(t *testing.T) {
	result, _ := template.Render("{% if a == 5 %}equal{% endif %}", map[string]any{"a": 5}, nil)
	if result != "equal" {
		t.Errorf("Expected 'equal', got '%s'", result)
	}

	result, _ = template.Render("{% if a == 5 %}equal{% endif %}", map[string]any{"a": 3}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionNotEquals(t *testing.T) {
	result, _ := template.Render("{% if a != 5 %}not equal{% endif %}", map[string]any{"a": 3}, nil)
	if result != "not equal" {
		t.Errorf("Expected 'not equal', got '%s'", result)
	}

	result, _ = template.Render("{% if a != 5 %}not equal{% endif %}", map[string]any{"a": 5}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLessThan(t *testing.T) {
	result, _ := template.Render("{% if a < 10 %}less{% endif %}", map[string]any{"a": 5}, nil)
	if result != "less" {
		t.Errorf("Expected 'less', got '%s'", result)
	}

	result, _ = template.Render("{% if a < 10 %}less{% endif %}", map[string]any{"a": 15}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionGreaterThan(t *testing.T) {
	result, _ := template.Render("{% if a > 10 %}greater{% endif %}", map[string]any{"a": 15}, nil)
	if result != "greater" {
		t.Errorf("Expected 'greater', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 10 %}greater{% endif %}", map[string]any{"a": 5}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLessThanOrEqual(t *testing.T) {
	result, _ := template.Render("{% if a <= 10 %}yes{% endif %}", map[string]any{"a": 10}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if a <= 10 %}yes{% endif %}", map[string]any{"a": 5}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if a <= 10 %}yes{% endif %}", map[string]any{"a": 15}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionGreaterThanOrEqual(t *testing.T) {
	result, _ := template.Render("{% if a >= 10 %}yes{% endif %}", map[string]any{"a": 10}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if a >= 10 %}yes{% endif %}", map[string]any{"a": 15}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if a >= 10 %}yes{% endif %}", map[string]any{"a": 5}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

// Expression tests - Logical operators
func TestExpressionLogicalAnd(t *testing.T) {
	result, _ := template.Render("{% if a > 5 && b < 20 %}both true{% endif %}", map[string]any{"a": 10, "b": 15}, nil)
	if result != "both true" {
		t.Errorf("Expected 'both true', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 && b < 20 %}both true{% endif %}", map[string]any{"a": 3, "b": 15}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 && b < 20 %}both true{% endif %}", map[string]any{"a": 10, "b": 25}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLogicalOr(t *testing.T) {
	result, _ := template.Render("{% if a > 5 || b < 20 %}at least one{% endif %}", map[string]any{"a": 10, "b": 25}, nil)
	if result != "at least one" {
		t.Errorf("Expected 'at least one', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 || b < 20 %}at least one{% endif %}", map[string]any{"a": 3, "b": 15}, nil)
	if result != "at least one" {
		t.Errorf("Expected 'at least one', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 || b < 20 %}at least one{% endif %}", map[string]any{"a": 3, "b": 25}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLogicalNot(t *testing.T) {
	result, _ := template.Render("{% if not a %}not true{% endif %}", map[string]any{"a": false}, nil)
	if result != "not true" {
		t.Errorf("Expected 'not true', got '%s'", result)
	}

	result, _ = template.Render("{% if not a %}not true{% endif %}", map[string]any{"a": true}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLogicalAndWordBased(t *testing.T) {
	result, _ := template.Render("{% if a > 5 and b < 20 %}both true{% endif %}", map[string]any{"a": 10, "b": 15}, nil)
	if result != "both true" {
		t.Errorf("Expected 'both true', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 and b < 20 %}both true{% endif %}", map[string]any{"a": 3, "b": 15}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 and b < 20 %}both true{% endif %}", map[string]any{"a": 10, "b": 25}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLogicalOrWordBased(t *testing.T) {
	result, _ := template.Render("{% if a > 5 or b < 20 %}at least one{% endif %}", map[string]any{"a": 10, "b": 25}, nil)
	if result != "at least one" {
		t.Errorf("Expected 'at least one', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 or b < 20 %}at least one{% endif %}", map[string]any{"a": 3, "b": 15}, nil)
	if result != "at least one" {
		t.Errorf("Expected 'at least one', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 or b < 20 %}at least one{% endif %}", map[string]any{"a": 3, "b": 25}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLogicalMixedWordAndSymbol(t *testing.T) {
	result, _ := template.Render("{% if a > 5 and b < 20 or c == 10 %}yes{% endif %}", map[string]any{"a": 10, "b": 15, "c": 0}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if a > 5 || b < 20 and c == 10 %}yes{% endif %}", map[string]any{"a": 10, "b": 25, "c": 5}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestExpressionLogicalNotWithAnd(t *testing.T) {
	result, _ := template.Render("{% if not a and b %}yes{% endif %}", map[string]any{"a": false, "b": true}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if not a and b %}yes{% endif %}", map[string]any{"a": true, "b": true}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}

	result, _ = template.Render("{% if not a and b %}yes{% endif %}", map[string]any{"a": false, "b": false}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionLogicalNotWithOr(t *testing.T) {
	result, _ := template.Render("{% if not a or b %}yes{% endif %}", map[string]any{"a": false, "b": false}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if not a or b %}yes{% endif %}", map[string]any{"a": true, "b": true}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if not a or b %}yes{% endif %}", map[string]any{"a": true, "b": false}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

// Expression tests - Arithmetic operators
func TestExpressionAddition(t *testing.T) {
	result, _ := template.Render("{{ a + b }}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "15" {
		t.Errorf("Expected '15', got '%s'", result)
	}
}

func TestExpressionSubtraction(t *testing.T) {
	result, _ := template.Render("{{ a - b }}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "5" {
		t.Errorf("Expected '5', got '%s'", result)
	}
}

func TestExpressionMultiplication(t *testing.T) {
	result, _ := template.Render("{{ a * b }}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "50" {
		t.Errorf("Expected '50', got '%s'", result)
	}
}

func TestExpressionDivision(t *testing.T) {
	result, _ := template.Render("{{ a / b }}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "2" {
		t.Errorf("Expected '2', got '%s'", result)
	}
}

func TestExpressionDivisionByZero(t *testing.T) {
	result, _ := template.Render("{{ a / 0 }}", map[string]any{"a": 10}, nil)
	if result != "{{a / 0!!division by zero}}" {
		t.Errorf("Expected error message, got '%s'", result)
	}
}

func TestExpressionModulo(t *testing.T) {
	result, _ := template.Render("{{ a % b }}", map[string]any{"a": 10, "b": 3}, nil)
	if result != "1" {
		t.Errorf("Expected '1', got '%s'", result)
	}
}

func TestExpressionModuloByZero(t *testing.T) {
	result, _ := template.Render("{{ a % 0 }}", map[string]any{"a": 10}, nil)
	if result != "{{a % 0!!modulo by zero}}" {
		t.Errorf("Expected error message, got '%s'", result)
	}
}

func TestExpressionNotEnoughOperandsUnary(t *testing.T) {
	result, _ := template.Render("{{ not }}", map[string]any{}, nil)
	if result != "{{not!!not enough operands for &#39;not&#39;}}" {
		t.Errorf("Expected error message, got '%s'", result)
	}
}

func TestExpressionNotEnoughOperandsBinary(t *testing.T) {
	result, _ := template.Render("{{ 5 + }}", map[string]any{}, nil)
	if result != "{{5 +!!not enough operands for &#39;+&#39;}}" {
		t.Errorf("Expected error message, got '%s'", result)
	}
}

func TestExpressionMalformedExpression(t *testing.T) {
	result, _ := template.Render("{{ 5 5 }}", map[string]any{}, nil)
	if result != "{{5 5!!malformed expression}}" {
		t.Errorf("Expected error message, got '%s'", result)
	}
}

// Expression tests - Operator precedence
func TestExpressionPrecedenceArithmetic(t *testing.T) {
	result, _ := template.Render("{{ a + b * c }}", map[string]any{"a": 2, "b": 3, "c": 4}, nil)
	if result != "14" {
		t.Errorf("Expected '14', got '%s'", result)
	}
}

func TestExpressionPrecedenceComparison(t *testing.T) {
	result, _ := template.Render("{% if a + b > c %}yes{% endif %}", map[string]any{"a": 5, "b": 10, "c": 12}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestExpressionPrecedenceLogical(t *testing.T) {
	// && has higher precedence than ||
	result, _ := template.Render("{% if a == 1 || b == 2 && c == 3 %}yes{% endif %}", map[string]any{"a": 5, "b": 2, "c": 3}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if a == 1 || b == 2 && c == 3 %}yes{% endif %}", map[string]any{"a": 5, "b": 2, "c": 5}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

// Expression tests - Nested expressions with parentheses
func TestExpressionParenthesesArithmetic(t *testing.T) {
	result, _ := template.Render("{{ (a + b) * c }}", map[string]any{"a": 2, "b": 4, "c": 3}, nil)
	if result != "18" {
		t.Errorf("Expected '18', got '%s'", result)
	}
}

func TestExpressionParenthesesLogical(t *testing.T) {
	result, _ := template.Render("{% if (a == 1 || b == 2) && c == 3 %}yes{% endif %}", map[string]any{"a": 5, "b": 2, "c": 3}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if (a == 1 || b == 2) && c == 3 %}yes{% endif %}", map[string]any{"a": 5, "b": 5, "c": 3}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionNestedParentheses(t *testing.T) {
	result, _ := template.Render("{% if ((a + b) * c) > 20 %}yes{% endif %}", map[string]any{"a": 5, "b": 5, "c": 3}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

// Expression tests - Complex combined conditions
func TestExpressionComplexCondition1(t *testing.T) {
	result, _ := template.Render("{% if a > 5 && b < 20 || c == 10 %}match{% endif %}", map[string]any{"a": 10, "b": 15, "c": 0}, nil)
	if result != "match" {
		t.Errorf("Expected 'match', got '%s'", result)
	}
}

func TestExpressionComplexCondition2(t *testing.T) {
	result, _ := template.Render("{% if (a > 5 || b > 5) && (c < 20 || d < 20) %}match{% endif %}", map[string]any{"a": 3, "b": 10, "c": 25, "d": 15}, nil)
	if result != "match" {
		t.Errorf("Expected 'match', got '%s'", result)
	}
}

func TestExpressionInElseIf(t *testing.T) {
	result, _ := template.Render("{% if a > 10 %}first{% elseif a > 5 %}second{% else %}third{% endif %}", map[string]any{"a": 7}, nil)
	if result != "second" {
		t.Errorf("Expected 'second', got '%s'", result)
	}
}

func TestExpressionWithStringLiterals(t *testing.T) {
	result, _ := template.Render("{% if name == \"John\" %}match{% endif %}", map[string]any{"name": "John"}, nil)
	if result != "match" {
		t.Errorf("Expected 'match', got '%s'", result)
	}

	result, _ = template.Render("{% if name == \"John\" %}match{% endif %}", map[string]any{"name": "Jane"}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionWithNumericLiterals(t *testing.T) {
	result, _ := template.Render("{% if a + 5 > 10 %}yes{% endif %}", map[string]any{"a": 8}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

// Expression tests - String concatenation
func TestStringConcatenationSimple(t *testing.T) {
	result, _ := template.Render("{{ first + second }}", map[string]any{"first": "hello", "second": "world"}, nil)
	if result != "helloworld" {
		t.Errorf("Expected 'helloworld', got '%s'", result)
	}
}

func TestStringConcatenationWithSpace(t *testing.T) {
	result, _ := template.Render("{{ first + \" \" + second }}", map[string]any{"first": "hello", "second": "world"}, nil)
	if result != "hello world" {
		t.Errorf("Expected 'hello world', got '%s'", result)
	}
}

func TestStringConcatenationMultiple(t *testing.T) {
	result, _ := template.Render("{{ first + \" \" + middle + \" \" + last }}", map[string]any{"first": "John", "middle": "Doe", "last": "Smith"}, nil)
	if result != "John Doe Smith" {
		t.Errorf("Expected 'John Doe Smith', got '%s'", result)
	}
}

func TestStringConcatenationWithNumber(t *testing.T) {
	result, _ := template.Render("{{ \"Value: \" + num }}", map[string]any{"num": 42}, nil)
	if result != "Value: 42" {
		t.Errorf("Expected 'Value: 42', got '%s'", result)
	}
}

func TestNumericAdditionStillWorks(t *testing.T) {
	result, _ := template.Render("{{ a + b }}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "15" {
		t.Errorf("Expected '15', got '%s'", result)
	}

	result, _ = template.Render("{{ a + b }}", map[string]any{"a": 5, "b": 2.5}, nil)
	if result != "7.5" {
		t.Errorf("Expected '7.5', got '%s'", result)
	}
}

func TestExpressionWithNestedPaths(t *testing.T) {
	result, _ := template.Render("{% if user.age >= 18 %}yes{% endif %}", map[string]any{"user": map[string]any{"age": 21}}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{% if user.age >= 18 %}yes{% endif %}", map[string]any{"user": map[string]any{"age": 16}}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionInVariableOutput(t *testing.T) {
	result, _ := template.Render("{{ a + b }}", map[string]any{"a": 3, "b": 5}, nil)
	if result != "8" {
		t.Errorf("Expected '8', got '%s'", result)
	}
}

func TestExpressionWithMultipleConditions(t *testing.T) {
	result, _ := template.Render("{% if score >= 60 && score <= 100 && not failed %}passed{% endif %}", map[string]any{"score": 75, "failed": false}, nil)
	if result != "passed" {
		t.Errorf("Expected 'passed', got '%s'", result)
	}
}

// Expression tests - Without spaces
func TestExpressionWithoutSpacesEquals(t *testing.T) {
	result, _ := template.Render("{%if a==5%}equal{%endif%}", map[string]any{"a": 5}, nil)
	if result != "equal" {
		t.Errorf("Expected 'equal', got '%s'", result)
	}

	result, _ = template.Render("{%if a==5%}equal{%endif%}", map[string]any{"a": 3}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestExpressionWithoutSpacesComparison(t *testing.T) {
	result, _ := template.Render("{%if a<10%}yes{%endif%}", map[string]any{"a": 5}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{%if a>10%}yes{%endif%}", map[string]any{"a": 15}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{%if a<=10%}yes{%endif%}", map[string]any{"a": 10}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{%if a>=10%}yes{%endif%}", map[string]any{"a": 10}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestExpressionWithoutSpacesArithmetic(t *testing.T) {
	result, _ := template.Render("{{a+b}}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "15" {
		t.Errorf("Expected '15', got '%s'", result)
	}

	result, _ = template.Render("{{a-b}}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "5" {
		t.Errorf("Expected '5', got '%s'", result)
	}

	result, _ = template.Render("{{a*b}}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "50" {
		t.Errorf("Expected '50', got '%s'", result)
	}

	result, _ = template.Render("{{a/b}}", map[string]any{"a": 10, "b": 5}, nil)
	if result != "2" {
		t.Errorf("Expected '2', got '%s'", result)
	}
}

func TestExpressionWithoutSpacesLogical(t *testing.T) {
	result, _ := template.Render("{%if a>5&&b<20%}yes{%endif%}", map[string]any{"a": 10, "b": 15}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{%if a>5||b>20%}yes{%endif%}", map[string]any{"a": 10, "b": 15}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestExpressionWithoutSpacesPrecedence(t *testing.T) {
	result, _ := template.Render("{{a+b*c}}", map[string]any{"a": 2, "b": 3, "c": 4}, nil)
	if result != "14" {
		t.Errorf("Expected '14', got '%s'", result)
	}

	result, _ = template.Render("{{(a+b)*c}}", map[string]any{"a": 2, "b": 4, "c": 3}, nil)
	if result != "18" {
		t.Errorf("Expected '18', got '%s'", result)
	}
}

func TestExpressionWithoutSpacesComplex(t *testing.T) {
	result, _ := template.Render("{%if a==1||b==2&&c==3%}yes{%endif%}", map[string]any{"a": 5, "b": 2, "c": 3}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{%if (a+b)>10&&c<20%}yes{%endif%}", map[string]any{"a": 7, "b": 5, "c": 15}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestExpressionMixedSpacing(t *testing.T) {
	result, _ := template.Render("{%if a+b>10&&c<20%}yes{%endif%}", map[string]any{"a": 7, "b": 5, "c": 15}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	result, _ = template.Render("{{a +b* c}}", map[string]any{"a": 5, "b": 3, "c": 4}, nil)
	if result != "17" {
		t.Errorf("Expected '17', got '%s'", result)
	}
}

// Multiline template tests inspired by Jinja
func TestMultilineForLoopSimple(t *testing.T) {
	tmpl := "<ul>\n{% for item in items %}\n    <li>{{ item }}</li>\n{% endfor %}\n</ul>"
	expected := "<ul>\n    <li>apple</li>\n    <li>banana</li>\n    <li>cherry</li>\n</ul>"
	result, _ := template.Render(tmpl, map[string]any{"items": []any{"apple", "banana", "cherry"}}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineForLoopWithIndentation(t *testing.T) {
	tmpl := "<div>\n    <ul>\n    {% for user in users %}\n        <li>{{ user }}</li>\n    {% endfor %}\n    </ul>\n</div>"
	expected := "<div>\n    <ul>\n        <li>Alice</li>\n        <li>Bob</li>\n    </ul>\n</div>"
	result, _ := template.Render(tmpl, map[string]any{"users": []any{"Alice", "Bob"}}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineIfWithWhitespace(t *testing.T) {
	tmpl := "<div>\n    {% if active %}\n        <span>Active</span>\n    {% endif %}\n</div>"
	expected := "<div>\n        <span>Active</span>\n</div>"
	result, _ := template.Render(tmpl, map[string]any{"active": true}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineIfElseWithWhitespace(t *testing.T) {
	tmpl := "<div>\n    {% if active %}\n        <span>Active</span>\n    {% else %}\n        <span>Inactive</span>\n    {% endif %}\n</div>"
	expected := "<div>\n        <span>Inactive</span>\n</div>"
	result, _ := template.Render(tmpl, map[string]any{"active": false}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineNestedForLoops(t *testing.T) {
	tmpl := "<table>\n{% for row in rows %}\n    <tr>\n    {% for cell in row %}\n        <td>{{ cell }}</td>\n    {% endfor %}\n    </tr>\n{% endfor %}\n</table>"
	expected := "<table>\n    <tr>\n        <td>1</td>\n        <td>2</td>\n    </tr>\n    <tr>\n        <td>3</td>\n        <td>4</td>\n    </tr>\n</table>"
	result, _ := template.Render(tmpl, map[string]any{"rows": []any{[]any{1, 2}, []any{3, 4}}}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineComplexHtmlStructure(t *testing.T) {
	tmpl := "<!DOCTYPE html>\n<html>\n<head>\n    <title>{{ title }}</title>\n</head>\n<body>\n    <ul id=\"navigation\">\n    {% for item in navigation %}\n        <li><a href=\"{{ item.href }}\">{{ item.caption }}</a></li>\n    {% endfor %}\n    </ul>\n    <h1>{{ heading }}</h1>\n</body>\n</html>"

	data := map[string]any{
		"title":   "My Page",
		"heading": "Welcome",
		"navigation": []any{
			map[string]any{"href": "/home", "caption": "Home"},
			map[string]any{"href": "/about", "caption": "About"},
		},
	}

	expected := "<!DOCTYPE html>\n<html>\n<head>\n    <title>My Page</title>\n</head>\n<body>\n    <ul id=\"navigation\">\n        <li><a href=\"/home\">Home</a></li>\n        <li><a href=\"/about\">About</a></li>\n    </ul>\n    <h1>Welcome</h1>\n</body>\n</html>"

	result, _ := template.Render(tmpl, data, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestWhitespacePreservationWithLeadingSpaces(t *testing.T) {
	tmpl := "    Leading spaces\n{{ text }}\n    Trailing spaces    "
	expected := "    Leading spaces\nHello\n    Trailing spaces    "
	result, _ := template.Render(tmpl, map[string]any{"text": "Hello"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestWhitespacePreservationWithTabs(t *testing.T) {
	tmpl := "\t\tTabbed content\n{{ text }}\n\t\tMore tabs"
	expected := "\t\tTabbed content\nWorld\n\t\tMore tabs"
	result, _ := template.Render(tmpl, map[string]any{"text": "World"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestWhitespacePreservationEmptyLines(t *testing.T) {
	tmpl := "Line 1\n\n{{ text }}\n\nLine 4"
	expected := "Line 1\n\nTest\n\nLine 4"
	result, _ := template.Render(tmpl, map[string]any{"text": "Test"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineForLoopWithEmptyList(t *testing.T) {
	tmpl := "<ul>\n{% for item in items %}\n    <li>{{ item }}</li>\n{% endfor %}\n</ul>"
	expected := "<ul>\n</ul>"
	result, _ := template.Render(tmpl, map[string]any{"items": []any{}}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineIfWithFalseCondition(t *testing.T) {
	tmpl := "<div>\n    Content before\n    {% if show %}\n        This should not appear\n    {% endif %}\n    Content after\n</div>"
	expected := "<div>\n    Content before\n    Content after\n</div>"
	result, _ := template.Render(tmpl, map[string]any{"show": false}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineTextPreservation(t *testing.T) {
	tmpl := "First line\nSecond line\nThird line with {{ var }}\nFourth line"
	expected := "First line\nSecond line\nThird line with value\nFourth line"
	result, _ := template.Render(tmpl, map[string]any{"var": "value"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineWithMixedContentTypes(t *testing.T) {
	tmpl := "<p>\n    Text content\n    {{ text }}\n    {% if show %}\n        <strong>{{ emphasis }}</strong>\n    {% endif %}\n    More text\n</p>"
	expected := "<p>\n    Text content\n    Hello\n        <strong>Important</strong>\n    More text\n</p>"
	result, _ := template.Render(tmpl, map[string]any{"text": "Hello", "show": true, "emphasis": "Important"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineHtmlListWithData(t *testing.T) {
	tmpl := "<h1>Members</h1>\n<ul>\n{% for user in users %}\n  <li>{{ user.username }}</li>\n{% endfor %}\n</ul>"
	expected := "<h1>Members</h1>\n<ul>\n  <li>alice</li>\n  <li>bob</li>\n  <li>charlie</li>\n</ul>"
	data := map[string]any{
		"users": []any{
			map[string]any{"username": "alice"},
			map[string]any{"username": "bob"},
			map[string]any{"username": "charlie"},
		},
	}
	result, _ := template.Render(tmpl, data, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineNestedIfStatements(t *testing.T) {
	tmpl := "<div>\n{% if outer %}\n    <div class=\"outer\">\n    {% if inner %}\n        <div class=\"inner\">Content</div>\n    {% endif %}\n    </div>\n{% endif %}\n</div>"
	expected := "<div>\n    <div class=\"outer\">\n        <div class=\"inner\">Content</div>\n    </div>\n</div>"
	result, _ := template.Render(tmpl, map[string]any{"outer": true, "inner": true}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineWhitespaceOnlyBetweenTags(t *testing.T) {
	tmpl := "<div>   \n   {{ text }}   \n   </div>"
	expected := "<div>   \n   Value   \n   </div>"
	result, _ := template.Render(tmpl, map[string]any{"text": "Value"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineCommentLikeStructure(t *testing.T) {
	// Test Jinja-style {# #} comment syntax - comments should be completely removed
	tmpl := "<div>\n    {# This is a comment #}\n    {{ content }}\n    {# Another comment #}\n</div>"
	expected := "<div>\n    Data\n</div>"
	result, _ := template.Render(tmpl, map[string]any{"content": "Data"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineForLoopWithComplexData(t *testing.T) {
	tmpl := "<dl>\n{% for item in items %}\n  <dt>{{ item.key }}</dt>\n  <dd>{{ item.value }}</dd>\n{% endfor %}\n</dl>"
	expected := "<dl>\n  <dt>Name</dt>\n  <dd>John</dd>\n  <dt>Age</dt>\n  <dd>30</dd>\n</dl>"
	data := map[string]any{
		"items": []any{
			map[string]any{"key": "Name", "value": "John"},
			map[string]any{"key": "Age", "value": "30"},
		},
	}
	result, _ := template.Render(tmpl, data, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineTemplateWithNoWhitespace(t *testing.T) {
	tmpl := "<ul>{% for i in items %}<li>{{ i }}</li>{% endfor %}</ul>"
	expected := "<ul><li>A</li><li>B</li></ul>"
	result, _ := template.Render(tmpl, map[string]any{"items": []any{"A", "B"}}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestMultilineIndentationVariations(t *testing.T) {
	tmpl := "<div>\n  Two spaces\n    Four spaces\n\tOne tab\n{{ text }}\n</div>"
	expected := "<div>\n  Two spaces\n    Four spaces\n\tOne tab\nValue\n</div>"
	result, _ := template.Render(tmpl, map[string]any{"text": "Value"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

// Comment syntax tests - {# ... #}
func TestCommentSimple(t *testing.T) {
	result, _ := template.Render("hello {# comment #} world", map[string]any{}, nil)
	if result != "hello  world" {
		t.Errorf("Expected 'hello  world', got '%s'", result)
	}
}

func TestCommentWithVariables(t *testing.T) {
	result, _ := template.Render("hello {# this is ignored #} {{ text }} world", map[string]any{"text": "test"}, nil)
	if result != "hello  test world" {
		t.Errorf("Expected 'hello  test world', got '%s'", result)
	}
}

func TestCommentMultiline(t *testing.T) {
	tmpl := "Line 1\n{# This is\na multiline\ncomment #}\nLine 2"
	expected := "Line 1\nLine 2"
	result, _ := template.Render(tmpl, map[string]any{}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestCommentWithControlStructures(t *testing.T) {
	result, _ := template.Render("{# comment #}{% if true %}result{% endif %}{# another #}", map[string]any{"true": true}, nil)
	if result != "result" {
		t.Errorf("Expected 'result', got '%s'", result)
	}
}

func TestCommentMultiple(t *testing.T) {
	result, _ := template.Render("a{# one #}b{# two #}c{# three #}", map[string]any{}, nil)
	if result != "abc" {
		t.Errorf("Expected 'abc', got '%s'", result)
	}
}

func TestCommentWithSpecialChars(t *testing.T) {
	result, _ := template.Render("{# {{ }} {% %} #}text", map[string]any{}, nil)
	if result != "text" {
		t.Errorf("Expected 'text', got '%s'", result)
	}
}

func TestCommentInTemplate(t *testing.T) {
	tmpl := "{# Header comment #}\n<div>\n    {# Content comment #}\n    {{ content }}\n</div>\n{# Footer comment #}"
	expected := "<div>\n    Data\n</div>\n"
	result, _ := template.Render(tmpl, map[string]any{"content": "Data"}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestCommentBeforeAndAfterVariable(t *testing.T) {
	result, _ := template.Render("{# before #}{{ text }}{# after #}", map[string]any{"text": "Value"}, nil)
	if result != "Value" {
		t.Errorf("Expected 'Value', got '%s'", result)
	}
}

func TestCommentInForLoop(t *testing.T) {
	tmpl := "{% for i in items %}{# loop comment #}{{ i }}{% endfor %}"
	expected := "123"
	result, _ := template.Render(tmpl, map[string]any{"items": []any{1, 2, 3}}, nil)
	if result != expected {
		t.Errorf("Expected '%s', got '%s'", expected, result)
	}
}

func TestCommentEmpty(t *testing.T) {
	result, _ := template.Render("{##}text", map[string]any{}, nil)
	if result != "text" {
		t.Errorf("Expected 'text', got '%s'", result)
	}
}

// Newlines in expressions tests
func TestExpressionWithNewlineInVariable(t *testing.T) {
	tmpl := "{{ a\n+ b }}"
	result, _ := template.Render(tmpl, map[string]any{"a": 10, "b": 5}, nil)
	if result != "15" {
		t.Errorf("Expected '15', got '%s'", result)
	}
}

func TestExpressionWithMultipleNewlinesInVariable(t *testing.T) {
	tmpl := "{{ a\n+\nb\n*\nc }}"
	result, _ := template.Render(tmpl, map[string]any{"a": 2, "b": 3, "c": 4}, nil)
	if result != "14" {
		t.Errorf("Expected '14', got '%s'", result)
	}
}

func TestExpressionWithNewlineInIfCondition(t *testing.T) {
	tmpl := "{% if a\n>\n5 %}yes{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"a": 10}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestExpressionWithNewlineInComplexCondition(t *testing.T) {
	tmpl := "{% if a\n>\n5\n&&\nb\n<\n20 %}match{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"a": 10, "b": 15}, nil)
	if result != "match" {
		t.Errorf("Expected 'match', got '%s'", result)
	}
}

func TestExpressionWithNewlineInParentheses(t *testing.T) {
	tmpl := "{{ (\na\n+\nb\n)\n*\nc }}"
	result, _ := template.Render(tmpl, map[string]any{"a": 2, "b": 4, "c": 3}, nil)
	if result != "18" {
		t.Errorf("Expected '18', got '%s'", result)
	}
}

func TestExpressionWithNewlineInComparison(t *testing.T) {
	tmpl := "{% if a\n==\n5 %}equal{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"a": 5}, nil)
	if result != "equal" {
		t.Errorf("Expected 'equal', got '%s'", result)
	}
}

func TestExpressionWithNewlineInLogicalOperators(t *testing.T) {
	tmpl := "{% if a\nand\nb\nor\nc %}yes{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"a": false, "b": false, "c": true}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestExpressionWithNewlineInForLoop(t *testing.T) {
	tmpl := "{% for i\nin\nitems %}{{ i }}{% endfor %}"
	result, _ := template.Render(tmpl, map[string]any{"items": []any{1, 2, 3}}, nil)
	if result != "123" {
		t.Errorf("Expected '123', got '%s'", result)
	}
}

func TestExpressionWithNewlineInStringConcatenationAndInString(t *testing.T) {
	tmpl := "{{ first\n+\n\"\n\"\n+\nsecond }}"
	result, _ := template.Render(tmpl, map[string]any{"first": "hello", "second": "world"}, nil)
	if result != "hello\nworld" {
		t.Errorf("Expected 'hello\\nworld', got '%s'", result)
	}
}

func TestExpressionWithNewlineBeforeFilter(t *testing.T) {
	functions := map[string]any{
		"capitalize": func(s string) string { return strings.ToUpper(s[:1]) + s[1:] },
	}
	tmpl := "{{ name\n|capitalize }}"
	result, _ := template.Render(tmpl, map[string]any{"name": "world"}, functions)
	if result != "World" {
		t.Errorf("Expected 'World', got '%s'", result)
	}
}

func TestExpressionWithNewlineInFilterArguments(t *testing.T) {
	functions := map[string]any{
		"dateFormat": func(dateStr string, format string) string {
			t, _ := time.Parse("January 2, 2006", dateStr)
			return t.Format("2006-01-02")
		},
	}
	tmpl := "{{ name\n|dateFormat(\n\"Y-m-d\"\n) }}"
	result, _ := template.Render(tmpl, map[string]any{"name": "May 13, 1980"}, functions)
	if result != "1980-05-13" {
		t.Errorf("Expected '1980-05-13', got '%s'", result)
	}
}

func TestExpressionWithCarriageReturnNewline(t *testing.T) {
	tmpl := "{{ a\r\n+\r\nb }}"
	result, _ := template.Render(tmpl, map[string]any{"a": 10, "b": 5}, nil)
	if result != "15" {
		t.Errorf("Expected '15', got '%s'", result)
	}
}

func TestExpressionWithMixedWhitespaceAndNewlines(t *testing.T) {
	tmpl := "{{ a  \n  +  \n  b  \n  *  \n  c }}"
	result, _ := template.Render(tmpl, map[string]any{"a": 2, "b": 3, "c": 4}, nil)
	if result != "14" {
		t.Errorf("Expected '14', got '%s'", result)
	}
}

func TestExpressionWithNewlineInElseIfCondition(t *testing.T) {
	tmpl := "{% if a\n>\n10 %}first{% elseif a\n>\n5 %}second{% else %}third{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"a": 7}, nil)
	if result != "second" {
		t.Errorf("Expected 'second', got '%s'", result)
	}
}

// Tests for builtin test functions

func TestIsDefined(t *testing.T) {
	tmpl := "{% if variable is defined %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"variable": "value"}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test undefined variable
	tmpl = "{% if missing is defined %}yes{% else %}no{% endif %}"
	result, _ = template.Render(tmpl, map[string]any{}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsUndefined(t *testing.T) {
	tmpl := "{% if missing is undefined %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test defined variable
	tmpl = "{% if variable is undefined %}yes{% else %}no{% endif %}"
	result, _ = template.Render(tmpl, map[string]any{"variable": "value"}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsEven(t *testing.T) {
	tmpl := "{% if num is even %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"num": 4}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test odd number
	result, _ = template.Render(tmpl, map[string]any{"num": 5}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsOdd(t *testing.T) {
	tmpl := "{% if num is odd %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"num": 5}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test even number
	result, _ = template.Render(tmpl, map[string]any{"num": 4}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsDivisibleBy(t *testing.T) {
	tmpl := "{% if num is divisibleby(3) %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"num": 9}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test not divisible
	result, _ = template.Render(tmpl, map[string]any{"num": 10}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}

	// Test divisible by 2
	tmpl = "{% if num is divisibleby(2) %}yes{% else %}no{% endif %}"
	result, _ = template.Render(tmpl, map[string]any{"num": 8}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestIsIterable(t *testing.T) {
	tmpl := "{% if items is iterable %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"items": []any{1, 2, 3}}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test non-iterable
	result, _ = template.Render(tmpl, map[string]any{"items": 123}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}

	// Test string is iterable
	result, _ = template.Render(tmpl, map[string]any{"items": "hello"}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestIsNull(t *testing.T) {
	tmpl := "{% if value is null %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"value": nil}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test non-null
	result, _ = template.Render(tmpl, map[string]any{"value": "something"}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsNumber(t *testing.T) {
	tmpl := "{% if value is number %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"value": 42}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test float
	result, _ = template.Render(tmpl, map[string]any{"value": 3.14}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test string number
	result, _ = template.Render(tmpl, map[string]any{"value": "123"}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test non-number string
	result, _ = template.Render(tmpl, map[string]any{"value": "hello"}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsString(t *testing.T) {
	tmpl := "{% if value is string %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"value": "hello"}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test number
	result, _ = template.Render(tmpl, map[string]any{"value": 123}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsNotTest(t *testing.T) {
	tmpl := "{% if value is not null %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"value": "something"}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}

	// Test with null value
	result, _ = template.Render(tmpl, map[string]any{"value": nil}, nil)
	if result != "no" {
		t.Errorf("Expected 'no', got '%s'", result)
	}
}

func TestIsTestInVariable(t *testing.T) {
	tmpl := "{{ num is even }}"
	result, _ := template.Render(tmpl, map[string]any{"num": 4}, nil)
	if result != "1" {
		t.Errorf("Expected '1', got '%s'", result)
	}

	// Test false case
	result, _ = template.Render(tmpl, map[string]any{"num": 5}, nil)
	if result != "" {
		t.Errorf("Expected '', got '%s'", result)
	}
}

func TestIsTestWithComplexExpression(t *testing.T) {
	tmpl := "{% if (num + 1) is even %}yes{% else %}no{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"num": 3}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}

func TestMultipleIsTests(t *testing.T) {
	tmpl := "{% if a is defined %}{% if a is even %}yes{% else %}no{% endif %}{% endif %}"
	result, _ := template.Render(tmpl, map[string]any{"a": 4}, nil)
	if result != "yes" {
		t.Errorf("Expected 'yes', got '%s'", result)
	}
}
