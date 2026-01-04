package tqtemplate

import (
	"fmt"
	"html"
	"regexp"
	"strconv"
	"strings"
	"unicode"
)

// RawValue marks a value that should not be escaped
type RawValue struct {
	Value string
}

// ExpressionToken represents a token in an expression
type ExpressionToken struct {
	Type  string // "number", "string", "identifier", "operator", "parenthesis"
	Value string
}

// Expression represents a parsed expression with operators
type Expression struct {
	tokens []ExpressionToken
}

type operator struct {
	precedence    int
	associativity string
}

var operators = map[string]operator{
	"or":  {precedence: 1, associativity: "left"},
	"||":  {precedence: 1, associativity: "left"},
	"and": {precedence: 2, associativity: "left"},
	"&&":  {precedence: 2, associativity: "left"},
	"==":  {precedence: 3, associativity: "left"},
	"!=":  {precedence: 3, associativity: "left"},
	"<":   {precedence: 4, associativity: "left"},
	">":   {precedence: 4, associativity: "left"},
	"<=":  {precedence: 4, associativity: "left"},
	">=":  {precedence: 4, associativity: "left"},
	"+":   {precedence: 5, associativity: "left"},
	"-":   {precedence: 5, associativity: "left"},
	"*":   {precedence: 6, associativity: "left"},
	"/":   {precedence: 6, associativity: "left"},
	"%":   {precedence: 6, associativity: "left"},
	"not": {precedence: 7, associativity: "right"},
}

// NewExpression creates a new expression from a string
func NewExpression(expr string) *Expression {
	e := &Expression{}
	e.tokens = e.tokenize(expr)
	return e
}

// tokenize converts an expression string into tokens
func (e *Expression) tokenize(expr string) []ExpressionToken {
	tokens := []ExpressionToken{}
	expr = strings.TrimSpace(expr)
	i := 0
	length := len(expr)

	for i < length {
		ch := expr[i]

		// Skip whitespace
		if unicode.IsSpace(rune(ch)) {
			i++
			continue
		}

		// Handle parentheses
		if ch == '(' || ch == ')' {
			tokens = append(tokens, ExpressionToken{Type: "parenthesis", Value: string(ch)})
			i++
			continue
		}

		// Handle word-based operators (and, or, not)
		if unicode.IsLetter(rune(ch)) {
			word := ""
			start := i
			for i < length && unicode.IsLetter(rune(expr[i])) {
				word += string(expr[i])
				i++
			}
			if _, exists := operators[word]; exists {
				tokens = append(tokens, ExpressionToken{Type: "operator", Value: word})
				continue
			}
			// Not an operator, reset and handle as identifier
			i = start
		}

		// Handle two-character operators
		if i < length-1 {
			twoChar := expr[i : i+2]
			if _, exists := operators[twoChar]; exists {
				tokens = append(tokens, ExpressionToken{Type: "operator", Value: twoChar})
				i += 2
				continue
			}
		}

		// Handle single-character operators
		if _, exists := operators[string(ch)]; exists {
			tokens = append(tokens, ExpressionToken{Type: "operator", Value: string(ch)})
			i++
			continue
		}

		// Handle numbers
		if unicode.IsDigit(rune(ch)) || (ch == '.' && i < length-1 && unicode.IsDigit(rune(expr[i+1]))) {
			num := ""
			for i < length && (unicode.IsDigit(rune(expr[i])) || expr[i] == '.') {
				num += string(expr[i])
				i++
			}
			tokens = append(tokens, ExpressionToken{Type: "number", Value: num})
			continue
		}

		// Handle string literals
		if ch == '"' {
			str := ""
			i++ // Skip opening quote
			escaped := false
			for i < length {
				if escaped {
					str += string(expr[i])
					escaped = false
				} else if expr[i] == '\\' {
					escaped = true
				} else if expr[i] == '"' {
					i++ // Skip closing quote
					break
				} else {
					str += string(expr[i])
				}
				i++
			}
			tokens = append(tokens, ExpressionToken{Type: "string", Value: str})
			continue
		}

		// Handle identifiers/paths (with dots for nested access)
		if unicode.IsLetter(rune(ch)) || ch == '_' {
			ident := ""
			for i < length && (unicode.IsLetter(rune(expr[i])) || unicode.IsDigit(rune(expr[i])) || expr[i] == '_' || expr[i] == '.') {
				ident += string(expr[i])
				i++
			}
			tokens = append(tokens, ExpressionToken{Type: "identifier", Value: ident})
			continue
		}

		// Unknown character, skip it
		i++
	}

	return tokens
}

// Evaluate evaluates the expression with the given data context
func (e *Expression) Evaluate(data map[string]any, resolvePath func(string, map[string]any) (any, error)) (any, error) {
	rpn := e.toReversePolishNotation()
	return e.evaluateRPN(rpn, data, resolvePath)
}

// toReversePolishNotation converts infix notation to RPN using Shunting Yard algorithm
func (e *Expression) toReversePolishNotation() []ExpressionToken {
	output := []ExpressionToken{}
	operatorStack := []ExpressionToken{}

	for _, token := range e.tokens {
		if token.Type == "number" || token.Type == "string" || token.Type == "identifier" {
			// Operand
			output = append(output, token)
		} else if token.Type == "parenthesis" && token.Value == "(" {
			operatorStack = append(operatorStack, token)
		} else if token.Type == "parenthesis" && token.Value == ")" {
			// Pop operators until we find the matching '('
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top.Type == "parenthesis" && top.Value == "(" {
					break
				}
				output = append(output, top)
				operatorStack = operatorStack[:len(operatorStack)-1]
			}
			if len(operatorStack) > 0 {
				operatorStack = operatorStack[:len(operatorStack)-1] // Remove the '('
			}
		} else if token.Type == "operator" {
			o1 := token.Value
			for len(operatorStack) > 0 {
				top := operatorStack[len(operatorStack)-1]
				if top.Type == "parenthesis" {
					break
				}
				if top.Type != "operator" {
					break
				}
				o2 := top.Value
				o1Prec := operators[o1].precedence
				o2Prec := operators[o2].precedence
				o1Assoc := operators[o1].associativity

				if (o1Assoc == "left" && o1Prec <= o2Prec) ||
					(o1Assoc == "right" && o1Prec < o2Prec) {
					output = append(output, top)
					operatorStack = operatorStack[:len(operatorStack)-1]
				} else {
					break
				}
			}
			operatorStack = append(operatorStack, token)
		}
	}

	// Pop remaining operators
	for len(operatorStack) > 0 {
		output = append(output, operatorStack[len(operatorStack)-1])
		operatorStack = operatorStack[:len(operatorStack)-1]
	}

	return output
}

// evaluateRPN evaluates an expression in Reverse Polish Notation
func (e *Expression) evaluateRPN(rpn []ExpressionToken, data map[string]any, resolvePath func(string, map[string]any) (any, error)) (any, error) {
	stack := []any{}

	for _, token := range rpn {
		if token.Type == "number" || token.Type == "string" || token.Type == "identifier" {
			// Operand
			if token.Type == "number" {
				if strings.Contains(token.Value, ".") {
					val, _ := strconv.ParseFloat(token.Value, 64)
					stack = append(stack, val)
				} else {
					val, _ := strconv.Atoi(token.Value)
					stack = append(stack, val)
				}
			} else if token.Type == "string" {
				stack = append(stack, token.Value)
			} else if token.Type == "identifier" {
				val, err := resolvePath(token.Value, data)
				if err != nil {
					return nil, err
				}
				stack = append(stack, val)
			}
		} else if token.Type == "operator" {
			op := token.Value
			if op == "not" {
				// Unary operator
				if len(stack) == 0 {
					return nil, fmt.Errorf("not enough operands for 'not'")
				}
				operand := stack[len(stack)-1]
				stack = stack[:len(stack)-1]
				stack = append(stack, !toBool(operand))
			} else {
				// Binary operator
				if len(stack) < 2 {
					return nil, fmt.Errorf("not enough operands for '%s'", op)
				}
				right := stack[len(stack)-1]
				left := stack[len(stack)-2]
				stack = stack[:len(stack)-2]

				result, err := e.applyOperator(op, left, right)
				if err != nil {
					return nil, err
				}
				stack = append(stack, result)
			}
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("malformed expression")
	}

	return stack[0], nil
}

// applyOperator applies a binary operator to two operands
func (e *Expression) applyOperator(op string, left, right any) (any, error) {
	switch op {
	case "or", "||":
		return toBool(left) || toBool(right), nil
	case "and", "&&":
		return toBool(left) && toBool(right), nil
	case "==":
		return compare(left, right) == 0, nil
	case "!=":
		return compare(left, right) != 0, nil
	case "<":
		return compare(left, right) < 0, nil
	case ">":
		return compare(left, right) > 0, nil
	case "<=":
		return compare(left, right) <= 0, nil
	case ">=":
		return compare(left, right) >= 0, nil
	case "+":
		// String concatenation or numeric addition
		leftNum, leftIsNum := toNumber(left)
		rightNum, rightIsNum := toNumber(right)
		if leftIsNum && rightIsNum {
			return leftNum + rightNum, nil
		}
		return toString(left) + toString(right), nil
	case "-":
		leftNum, _ := toNumber(left)
		rightNum, _ := toNumber(right)
		return leftNum - rightNum, nil
	case "*":
		leftNum, _ := toNumber(left)
		rightNum, _ := toNumber(right)
		return leftNum * rightNum, nil
	case "/":
		leftNum, _ := toNumber(left)
		rightNum, rightIsNum := toNumber(right)
		if !rightIsNum || rightNum == 0 {
			return nil, fmt.Errorf("division by zero")
		}
		return leftNum / rightNum, nil
	case "%":
		leftNum, _ := toNumber(left)
		rightNum, rightIsNum := toNumber(right)
		if !rightIsNum || rightNum == 0 {
			return nil, fmt.Errorf("modulo by zero")
		}
		return int(leftNum) % int(rightNum), nil
	default:
		return nil, fmt.Errorf("unknown operator: %s", op)
	}
}

// TreeNode represents a node in the template syntax tree
type TreeNode struct {
	Type       string
	Expression string
	Children   []*TreeNode
	Value      any
}

// Template is the main template engine
type Template struct {
	escape string
}

// NewTemplate creates a new template engine with the specified escape type
func NewTemplate(escape string) *Template {
	return &Template{escape: escape}
}

// Render renders a template string with the provided data and custom functions
func (t *Template) Render(template string, data map[string]any, functions map[string]any) (string, error) {
	tokens := t.tokenize(template)
	tree := t.createSyntaxTree(tokens)

	// Add built-in 'raw' filter
	if functions == nil {
		functions = make(map[string]any)
	}
	functions["raw"] = func(value string) RawValue {
		return RawValue{Value: value}
	}

	return t.renderChildren(tree, data, functions)
}

// escapeValue escapes a value based on the escape type
func (t *Template) escapeValue(value any) string {
	if rawVal, ok := value.(RawValue); ok {
		return rawVal.Value
	}
	str := toString(value)
	if t.escape == "html" {
		return html.EscapeString(str)
	}
	return str
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
				ch := template[i]
				if !escaped {
					if ch == '"' {
						quoted = !quoted
					} else if ch == '\\' {
						escaped = true
					} else if !quoted && ch == '%' && template[i+1] == '}' {
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
				expr += string(ch)
				i++
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
				ch := template[i]
				if !escaped {
					if ch == '"' {
						quoted = !quoted
					} else if ch == '\\' {
						escaped = true
					} else if !quoted && ch == '}' && template[i+1] == '}' {
						tokens = append(tokens, strings.TrimSpace(expr))
						i += 2
						break
					}
				} else {
					escaped = false
				}
				expr += string(ch)
				i++
			}
			continue
		}

		// Regular character
		literal += string(template[i])
		i++
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

	for i := 0; i < len(str); i++ {
		ch := rune(str[i])
		if !quoted {
			if ch == quote {
				quoted = true
			} else if strings.HasPrefix(str[i:], separator) {
				// Special handling for | separator: check if it's part of || operator
				if separator == "|" && i+1 < len(str) && str[i+1] == '|' {
					// This is part of || operator, don't split
					token += "||"
					i++
					continue
				}
				tokens = append(tokens, token)
				if count > 0 && len(tokens) == count-1 {
					token = str[i+len(separator):]
					break
				}
				token = ""
				i += len(separator) - 1
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
			} else {
				nodeType = "var"
				expression = token
			}

			if nodeType == "endif" || nodeType == "endfor" || nodeType == "elseif" || nodeType == "else" {
				if len(stack) > 0 {
					current = stack[len(stack)-1]
					stack = stack[:len(stack)-1]
				}
			}

			if nodeType == "var" {
				node := &TreeNode{Type: nodeType, Expression: expression}
				current.Children = append(current.Children, node)
			}

			if nodeType == "if" || nodeType == "for" || nodeType == "elseif" || nodeType == "else" {
				node := &TreeNode{Type: nodeType, Expression: expression}
				current.Children = append(current.Children, node)
				stack = append(stack, current)
				current = node
			}
		} else {
			current.Children = append(current.Children, &TreeNode{Type: "lit", Expression: token})
		}
	}

	return root
}

// renderChildren renders all child nodes of a given node
func (t *Template) renderChildren(node *TreeNode, data map[string]any, functions map[string]any) (string, error) {
	result := ""
	ifNodes := []*TreeNode{}

	for _, child := range node.Children {
		switch child.Type {
		case "if":
			output, err := t.renderIfNode(child, data, functions)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{child}
		case "elseif":
			output, err := t.renderElseIfNode(child, ifNodes, data, functions)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = append(ifNodes, child)
		case "else":
			output, err := t.renderElseNode(child, ifNodes, data, functions)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{}
		case "for":
			output, err := t.renderForNode(child, data, functions)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{}
		case "var":
			output, err := t.renderVarNode(child, data, functions)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{}
		case "lit":
			result += child.Expression
			ifNodes = []*TreeNode{}
		}
	}

	return result, nil
}

// renderIfNode renders an 'if' conditional node
func (t *Template) renderIfNode(node *TreeNode, data map[string]any, functions map[string]any) (string, error) {
	expressionStr := node.Expression
	parts := t.explodeRespectingQuotes("|", expressionStr, -1)
	exprPart := parts[0]
	filterParts := parts[1:]

	expr := NewExpression(exprPart)
	value, err := expr.Evaluate(data, t.resolvePath)
	if err != nil {
		return t.escapeValue("{% if " + expressionStr + "!!" + err.Error() + " %}"), nil
	}

	value, err = t.applyFunctions(value, filterParts, functions, data)
	if err != nil {
		return t.escapeValue("{% if " + expressionStr + "!!" + err.Error() + " %}"), nil
	}

	result := ""
	if toBool(value) {
		output, err := t.renderChildren(node, data, functions)
		if err != nil {
			return "", err
		}
		result += output
	}
	node.Value = toBool(value)
	return result, nil
}

// renderElseIfNode renders an 'elseif' conditional node
func (t *Template) renderElseIfNode(node *TreeNode, ifNodes []*TreeNode, data map[string]any, functions map[string]any) (string, error) {
	if len(ifNodes) < 1 || ifNodes[0].Type != "if" {
		return t.escapeValue("{% elseif !!could not find matching `if` %}"), nil
	}

	result := ""
	anyTrue := false
	for _, ifNode := range ifNodes {
		if val, ok := ifNode.Value.(bool); ok && val {
			anyTrue = true
			break
		}
	}

	if !anyTrue {
		expressionStr := node.Expression
		parts := t.explodeRespectingQuotes("|", expressionStr, -1)
		exprPart := parts[0]
		filterParts := parts[1:]

		expr := NewExpression(exprPart)
		value, err := expr.Evaluate(data, t.resolvePath)
		if err != nil {
			return t.escapeValue("{% elseif " + expressionStr + "!!" + err.Error() + " %}"), nil
		}

		value, err = t.applyFunctions(value, filterParts, functions, data)
		if err != nil {
			return t.escapeValue("{% elseif " + expressionStr + "!!" + err.Error() + " %}"), nil
		}

		if toBool(value) {
			output, err := t.renderChildren(node, data, functions)
			if err != nil {
				return "", err
			}
			result += output
		}
		node.Value = toBool(value)
	} else {
		node.Value = false
	}

	return result, nil
}

// renderElseNode renders an 'else' node
func (t *Template) renderElseNode(node *TreeNode, ifNodes []*TreeNode, data map[string]any, functions map[string]any) (string, error) {
	if len(ifNodes) < 1 || ifNodes[0].Type != "if" {
		return t.escapeValue("{% else !!could not find matching `if` %}"), nil
	}

	result := ""
	anyTrue := false
	for _, ifNode := range ifNodes {
		if val, ok := ifNode.Value.(bool); ok && val {
			anyTrue = true
			break
		}
	}

	if !anyTrue {
		output, err := t.renderChildren(node, data, functions)
		if err != nil {
			return "", err
		}
		result += output
	}

	return result, nil
}

// renderForNode renders a 'for' loop node
func (t *Template) renderForNode(node *TreeNode, data map[string]any, functions map[string]any) (string, error) {
	expressionStr := node.Expression

	// Parse "for key, value in array" or "for value in array"
	re := regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_]*(?:\s*,\s*[a-zA-Z_][a-zA-Z0-9_]*)?)\s+in\s+(.+)$`)
	matches := re.FindStringSubmatch(expressionStr)
	if matches == nil {
		return t.escapeValue(`{% for ` + expressionStr + `!!invalid syntax, expected "item in array" or "key, value in array" %}`), nil
	}

	vars := matches[1]
	arrayExpr := matches[2]

	// Check if we have "key, value" or just "value"
	var key, varName string
	hasKey := strings.Contains(vars, ",")
	if hasKey {
		varParts := strings.Split(vars, ",")
		key = strings.TrimSpace(varParts[0])
		varName = strings.TrimSpace(varParts[1])
	} else {
		varName = strings.TrimSpace(vars)
	}

	// Parse filters from array expression
	parts := t.explodeRespectingQuotes("|", arrayExpr, -1)
	path := strings.TrimSpace(parts[0])
	filterParts := parts[1:]

	value, err := t.resolvePath(path, data)
	if err != nil {
		return t.escapeValue("{% for " + expressionStr + "!!" + err.Error() + " %}"), nil
	}

	value, err = t.applyFunctions(value, filterParts, functions, data)
	if err != nil {
		return t.escapeValue("{% for " + expressionStr + "!!" + err.Error() + " %}"), nil
	}

	// Convert to slice
	var items []any
	var keys []any

	switch v := value.(type) {
	case []any:
		items = v
		for i := range items {
			keys = append(keys, i)
		}
	case map[string]any:
		for k, val := range v {
			keys = append(keys, k)
			items = append(items, val)
		}
	default:
		return t.escapeValue("{% for " + expressionStr + "!!expression must evaluate to an array %}"), nil
	}

	result := ""
	for i, item := range items {
		newData := make(map[string]any)
		for k, v := range data {
			newData[k] = v
		}
		if hasKey {
			newData[key] = keys[i]
			newData[varName] = item
		} else {
			newData[varName] = item
		}
		output, err := t.renderChildren(node, newData, functions)
		if err != nil {
			return "", err
		}
		result += output
	}

	return result, nil
}

// renderVarNode renders a variable interpolation node
func (t *Template) renderVarNode(node *TreeNode, data map[string]any, functions map[string]any) (string, error) {
	expressionStr := node.Expression
	parts := t.explodeRespectingQuotes("|", expressionStr, -1)
	exprPart := parts[0]
	filterParts := parts[1:]

	expr := NewExpression(exprPart)
	value, err := expr.Evaluate(data, t.resolvePath)
	if err != nil {
		return t.escapeValue("{{" + expressionStr + "!!" + err.Error() + "}}"), nil
	}

	value, err = t.applyFunctions(value, filterParts, functions, data)
	if err != nil {
		return t.escapeValue("{{" + expressionStr + "!!" + err.Error() + "}}"), nil
	}

	if rawVal, ok := value.(RawValue); ok {
		return rawVal.Value, nil
	}

	return t.escapeValue(value), nil
}

// resolvePath resolves a dot-notation path to retrieve a value from data
func (t *Template) resolvePath(path string, data map[string]any) (any, error) {
	parts := t.explodeRespectingQuotes(".", path, -1)
	current := any(data)

	for _, part := range parts {
		if m, ok := current.(map[string]any); ok {
			if val, exists := m[part]; exists {
				current = val
			} else {
				return nil, fmt.Errorf("path `%s` not found", part)
			}
		} else {
			return nil, fmt.Errorf("path `%s` not found", part)
		}
	}

	return current, nil
}

// applyFunctions applies a chain of filter functions to a value
func (t *Template) applyFunctions(value any, parts []string, functions map[string]any, data map[string]any) (any, error) {
	for _, part := range parts {
		funcParts := t.explodeRespectingQuotes("(", strings.TrimSuffix(part, ")"), 2)
		funcName := funcParts[0]
		var arguments []any

		if len(funcParts) > 1 {
			argStrs := t.explodeRespectingQuotes(",", funcParts[1], -1)
			for _, argStr := range argStrs {
				argStr = strings.TrimSpace(argStr)
				argLen := len(argStr)
				if argLen > 1 && argStr[0] == '"' && argStr[argLen-1] == '"' {
					// String literal - unescape
					unescaped := argStr[1 : argLen-1]
					unescaped = strings.ReplaceAll(unescaped, "\\n", "\n")
					unescaped = strings.ReplaceAll(unescaped, "\\t", "\t")
					unescaped = strings.ReplaceAll(unescaped, "\\\"", "\"")
					unescaped = strings.ReplaceAll(unescaped, "\\\\", "\\")
					arguments = append(arguments, unescaped)
				} else if num, err := strconv.ParseFloat(argStr, 64); err == nil {
					// Numeric literal - convert to appropriate numeric type
					if strings.Contains(argStr, ".") {
						arguments = append(arguments, num) // float64
					} else {
						arguments = append(arguments, int(num)) // int
					}
				} else {
					// Path reference
					val, err := t.resolvePath(argStr, data)
					if err != nil {
						return nil, err
					}
					arguments = append(arguments, val)
				}
			}
		}

		// Prepend the value as the first argument
		allArgs := append([]any{value}, arguments...)

		// Call the function
		if fn, exists := functions[funcName]; exists {
			result, err := callFunction(fn, allArgs)
			if err != nil {
				return nil, err
			}
			value = result
		} else {
			return nil, fmt.Errorf("function `%s` not found", funcName)
		}
	}

	return value, nil
}

// Helper functions

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

func callFunction(fn any, args []any) (any, error) {
	switch f := fn.(type) {
	case func(string) RawValue:
		if len(args) > 0 {
			if str, ok := args[0].(string); ok {
				return f(str), nil
			}
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
	case func(any, any) bool:
		if len(args) >= 2 {
			return f(args[0], args[1]), nil
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
