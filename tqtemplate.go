package tqtemplate

import (
	"fmt"
	"html"
	"strings"
	"unicode/utf8"
)

// RawValue marks a value that should not be escaped
type RawValue struct {
	Value string
}

// TreeNode represents a node in the template syntax tree
type TreeNode struct {
	Type       string
	Expression string
	Children   []*TreeNode
	Value      any
}

// TemplateLoader is a function that loads template content by name
type TemplateLoader func(name string) (string, error)

// Template is the main template engine
type Template struct {
	loader  TemplateLoader
	filters map[string]any
	tests   map[string]any
}

// NewTemplate creates a new template engine
func NewTemplate() *Template {
	return NewTemplateWithLoaderAndFiltersAndTests(nil, nil, nil)
}

// NewTemplateWithLoader creates a new template engine with a custom template loader
func NewTemplateWithLoader(loader TemplateLoader) *Template {
	return NewTemplateWithLoaderAndFiltersAndTests(loader, nil, nil)

}

// NewTemplateWithLoader creates a new template engine with a custom template loader
func NewTemplateWithLoaderAndFilters(loader TemplateLoader, customFilters map[string]any) *Template {
	return NewTemplateWithLoaderAndFiltersAndTests(loader, customFilters, nil)
}

func NewTemplateWithLoaderAndFiltersAndTests(loader TemplateLoader, customFilters map[string]any, customTests map[string]any) *Template {
	return &Template{
		loader:  loader,
		filters: customFilters,
		tests:   customTests,
	}
}

// RenderFile renders a template file with the provided data
func (t *Template) RenderFile(templateFile string, data map[string]any) (string, error) {
	if t.loader == nil {
		return "", fmt.Errorf("no template loader defined")
	}
	templateContent, err := t.loader(templateFile)
	if err != nil {
		return "", err
	}
	return t.Render(templateContent, data)
}

// Render renders a template string with the provided data
func (t *Template) Render(template string, data map[string]any) (string, error) {
	tokens := t.tokenize(template)
	tree := t.createSyntaxTree(tokens)

	// Initialize filters map if needed
	filters := make(map[string]any)

	// Register all builtin filters
	builtins := getBuiltinFilters()
	for name, fn := range builtins {
		filters[name] = fn
	}

	// Collect all tests (builtin + custom)
	allTests := make(map[string]any)

	// Register all builtin tests
	tests := getBuiltinTests()
	for name, fn := range tests {
		allTests[name] = fn
		filters[name] = fn
	}

	// Add custom tests (allow user overrides)
	if t.tests != nil {
		for name, fn := range t.tests {
			allTests[name] = fn
			filters[name] = fn
		}
	}

	// Create __istest__ and __isnot__ filters with access to all tests
	filters["__istest__"] = createFilterIsTest(allTests)
	filters["__isnot__"] = createFilterIsNot(allTests)

	// Add custom filters (allow user overrides)
	if t.filters != nil {
		for name, fn := range t.filters {
			filters[name] = fn
		}
	}

	// Check if this template extends another template
	// Extends must be the first non-literal node
	extendsNode := t.findExtendsNode(tree)
	if extendsNode != nil {
		return t.renderWithExtends(tree, extendsNode, data, filters)
	}

	return t.renderChildren(tree, data, filters)
}

// escapeValue escapes a value for HTML output
func (t *Template) escapeValue(value any) string {
	if rawVal, ok := value.(RawValue); ok {
		return rawVal.Value
	}
	str := toString(value)
	return html.EscapeString(str)
}

// tokenize splits a template into literal text and expressions
func (t *Template) tokenize(template string) []string {
	tokens := []string{}
	i := 0
	length := len(template)
	literal := ""

	for i < length {
		// Check for comment {#
		if i < length-1 && template[i] == '{' && template[i+1] == '#' {
			// Check if this comment is on a standalone line
			lineStart := strings.LastIndex(literal, "\n")
			beforeTag := ""
			isStandaloneLine := false

			if lineStart == -1 {
				beforeTag = literal
				isStandaloneLine = literal == "" || (i == len(beforeTag) && strings.TrimSpace(beforeTag) == "")
			} else {
				beforeTag = literal[lineStart+1:]
				isStandaloneLine = strings.TrimSpace(beforeTag) == ""
			}

			// If standalone, remove just the whitespace on this line
			if isStandaloneLine && lineStart != -1 {
				literal = literal[:lineStart+1]
			} else if isStandaloneLine && lineStart == -1 {
				literal = ""
			}

			// Skip the comment - find closing #}
			i += 2
			commentEnd := i
			for commentEnd < length-1 {
				if template[commentEnd] == '#' && template[commentEnd+1] == '}' {
					commentEnd += 2
					break
				}
				commentEnd++
			}

			// If standalone line, consume trailing newline after comment
			if isStandaloneLine && commentEnd < length && template[commentEnd] == '\n' {
				commentEnd++
			} else if isStandaloneLine && commentEnd < length-1 && template[commentEnd] == '\r' && template[commentEnd+1] == '\n' {
				commentEnd += 2
			}

			i = commentEnd
			continue
		}

		// Check for control structure {%
		if i < length-1 && template[i] == '{' && template[i+1] == '%' {
			// Check if this control structure is on a standalone line
			lineStart := strings.LastIndex(literal, "\n")
			beforeTag := ""
			isStandaloneLine := false

			if lineStart == -1 {
				beforeTag = literal
				isStandaloneLine = literal == "" || (i == len(beforeTag) && strings.TrimSpace(beforeTag) == "")
			} else {
				beforeTag = literal[lineStart+1:]
				isStandaloneLine = strings.TrimSpace(beforeTag) == ""
			}

			// If standalone, remove just the whitespace on this line
			if isStandaloneLine && lineStart != -1 {
				literal = literal[:lineStart+1]
			} else if isStandaloneLine && lineStart == -1 {
				literal = ""
			}

			tokens = append(tokens, literal)
			literal = ""
			i += 2
			expr := ""
			quoted := false
			escaped := false
			for i < length-1 {
				r, size := utf8.DecodeRuneInString(template[i:])
				if !escaped {
					if r == '"' {
						quoted = !quoted
					} else if r == '\\' {
						escaped = true
					} else if !quoted && r == '%' && i+1 < length && template[i+1] == '}' {
						tokens = append(tokens, "@"+strings.TrimSpace(expr))
						i += 2

						// If standalone line, consume trailing newline
						if isStandaloneLine && i < length && template[i] == '\n' {
							i++
						} else if isStandaloneLine && i < length-1 && template[i] == '\r' && template[i+1] == '\n' {
							i += 2
						}
						break
					}
				} else {
					escaped = false
				}
				expr += string(r)
				i += size
			}
			continue
		}

		// Check for variable {{
		if i < length-1 && template[i] == '{' && template[i+1] == '{' {
			tokens = append(tokens, literal)
			literal = ""
			i += 2
			expr := ""
			quoted := false
			escaped := false
			for i < length-1 {
				r, size := utf8.DecodeRuneInString(template[i:])
				if !escaped {
					if r == '"' {
						quoted = !quoted
					} else if r == '\\' {
						escaped = true
					} else if !quoted && r == '}' && i+1 < length && template[i+1] == '}' {
						tokens = append(tokens, strings.TrimSpace(expr))
						i += 2
						break
					}
				} else {
					escaped = false
				}
				expr += string(r)
				i += size
			}
			continue
		}

		// Regular character - read full UTF-8 rune
		r, size := utf8.DecodeRuneInString(template[i:])
		literal += string(r)
		i += size
	}

	tokens = append(tokens, literal)
	return tokens
}

// explodeRespectingQuotes splits a string by separator, respecting quoted substrings
func (t *Template) explodeRespectingQuotes(separator, str string, count int) []string {
	if count == -1 {
		count = 0
	}
	tokens := []string{}
	token := ""
	quote := '"'
	escape := '\\'
	escaped := false
	quoted := false

	for i := 0; i < len(str); {
		ch, size := utf8.DecodeRuneInString(str[i:])
		if !quoted {
			if ch == quote {
				quoted = true
			} else if strings.HasPrefix(str[i:], separator) {
				// Special handling for | separator: check if it's part of || operator
				if separator == "|" && i+1 < len(str) && str[i+1] == '|' {
					// This is part of || operator, don't split
					token += "||"
					i += 2
					continue
				}
				tokens = append(tokens, token)
				if count > 0 && len(tokens) == count-1 {
					token = str[i+len(separator):]
					break
				}
				token = ""
				i += len(separator)
				continue
			}
		} else {
			if !escaped {
				if ch == quote {
					quoted = false
				} else if ch == escape {
					escaped = true
				}
			} else {
				escaped = false
			}
		}
		token += string(ch)
		i += size
	}
	tokens = append(tokens, token)
	return tokens
}

// createSyntaxTree creates an abstract syntax tree from tokens
func (t *Template) createSyntaxTree(tokens []string) *TreeNode {
	root := &TreeNode{Type: "root"}
	current := root
	stack := []*TreeNode{}

	for i, token := range tokens {
		if i%2 == 1 {
			// Control structures are prefixed with @
			isControl := strings.HasPrefix(token, "@")
			if isControl {
				token = token[1:] // Remove @ prefix
			}

			nodeType := ""
			expression := ""

			if token == "endif" {
				nodeType = "endif"
			} else if token == "endfor" {
				nodeType = "endfor"
			} else if token == "endblock" {
				nodeType = "endblock"
			} else if token == "else" {
				nodeType = "else"
			} else if strings.HasPrefix(token, "elseif ") {
				nodeType = "elseif"
				expression = strings.TrimSpace(token[7:])
			} else if strings.HasPrefix(token, "if ") {
				nodeType = "if"
				expression = strings.TrimSpace(token[3:])
			} else if strings.HasPrefix(token, "for ") {
				nodeType = "for"
				expression = strings.TrimSpace(token[4:])
			} else if strings.HasPrefix(token, "block ") {
				nodeType = "block"
				expression = strings.TrimSpace(token[6:])
			} else if strings.HasPrefix(token, "extends ") {
				nodeType = "extends"
				expression = strings.TrimSpace(token[8:])
			} else if strings.HasPrefix(token, "include ") {
				nodeType = "include"
				expression = strings.TrimSpace(token[8:])
			} else {
				nodeType = "var"
				expression = token
			}

			if nodeType == "endif" || nodeType == "endfor" || nodeType == "endblock" || nodeType == "elseif" || nodeType == "else" {
				if len(stack) > 0 {
					current = stack[len(stack)-1]
					stack = stack[:len(stack)-1]
				}
			}

			if nodeType == "var" {
				node := &TreeNode{Type: nodeType, Expression: expression}
				current.Children = append(current.Children, node)
			}

			if nodeType == "if" || nodeType == "for" || nodeType == "block" || nodeType == "elseif" || nodeType == "else" {
				node := &TreeNode{Type: nodeType, Expression: expression}
				current.Children = append(current.Children, node)
				stack = append(stack, current)
				current = node
			}

			if nodeType == "extends" || nodeType == "include" {
				node := &TreeNode{Type: nodeType, Expression: expression}
				current.Children = append(current.Children, node)
			}
		} else {
			current.Children = append(current.Children, &TreeNode{Type: "lit", Expression: token})
		}
	}

	return root
}
