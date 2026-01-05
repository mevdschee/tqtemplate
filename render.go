package tqtemplate

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// renderChildren renders all child nodes of a given node
func (t *Template) renderChildren(node *TreeNode, data map[string]any, filters map[string]any) (string, error) {
	result := ""
	ifNodes := []*TreeNode{}

	for _, child := range node.Children {
		switch child.Type {
		case "block":
			// Render block content directly when not in extends context
			output, err := t.renderChildren(child, data, filters)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{}
		case "if":
			output, err := t.renderIfNode(child, data, filters)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{child}
		case "elseif":
			output, err := t.renderElseIfNode(child, ifNodes, data, filters)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = append(ifNodes, child)
		case "else":
			output, err := t.renderElseNode(child, ifNodes, data, filters)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{}
		case "for":
			output, err := t.renderForNode(child, data, filters)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{}
		case "var":
			output, err := t.renderVarNode(child, data, filters)
			if err != nil {
				return "", err
			}
			result += output
			ifNodes = []*TreeNode{}
		case "include":
			output, err := t.renderIncludeNode(child, data, filters)
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
func (t *Template) renderIfNode(node *TreeNode, data map[string]any, filters map[string]any) (string, error) {
	expressionStr := node.Expression

	// Preprocess "is" tests
	exprPart, testFilter := processIsTests(expressionStr)

	parts := t.explodeRespectingQuotes("|", exprPart, -1)
	actualExpr := parts[0]
	filterParts := parts[1:]

	// Add test filter if present
	if testFilter != "" {
		filterParts = append(filterParts, testFilter)
	}

	expr := NewExpression(actualExpr)
	value, err := expr.Evaluate(data, t.resolvePath)

	// Special handling for "defined" and "undefined" tests
	// If we have an error and the test is for defined/undefined, handle it specially
	if err != nil && (strings.Contains(testFilter, "__istest__(\"defined\")") ||
		strings.Contains(testFilter, "__istest__(\"undefined\")") ||
		strings.Contains(testFilter, "__isnot__(\"defined\")") ||
		strings.Contains(testFilter, "__isnot__(\"undefined\")")) {
		// For defined/undefined tests, use sentinel value to indicate undefined
		value = undefinedValue
		err = nil
	}

	if err != nil {
		return t.escapeValue("{% if " + expressionStr + "!!" + err.Error() + " %}"), nil
	}

	value, err = t.applyfilters(value, filterParts, filters, data)
	if err != nil {
		return t.escapeValue("{% if " + expressionStr + "!!" + err.Error() + " %}"), nil
	}

	result := ""
	if toBool(value) {
		output, err := t.renderChildren(node, data, filters)
		if err != nil {
			return "", err
		}
		result += output
	}
	node.Value = toBool(value)
	return result, nil
}

// renderElseIfNode renders an 'elseif' conditional node
func (t *Template) renderElseIfNode(node *TreeNode, ifNodes []*TreeNode, data map[string]any, filters map[string]any) (string, error) {
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

		// Preprocess "is" tests
		exprPart, testFilter := processIsTests(expressionStr)

		parts := t.explodeRespectingQuotes("|", exprPart, -1)
		actualExpr := parts[0]
		filterParts := parts[1:]

		// Add test filter if present
		if testFilter != "" {
			filterParts = append(filterParts, testFilter)
		}

		expr := NewExpression(actualExpr)
		value, err := expr.Evaluate(data, t.resolvePath)

		// Special handling for "defined" and "undefined" tests
		// If we have an error and the test is for defined/undefined, handle it specially
		if err != nil && (strings.Contains(testFilter, "__istest__(\"defined\")") ||
			strings.Contains(testFilter, "__istest__(\"undefined\")") ||
			strings.Contains(testFilter, "__isnot__(\"defined\")") ||
			strings.Contains(testFilter, "__isnot__(\"undefined\")")) {
			// For defined/undefined tests, use sentinel value to indicate undefined
			value = undefinedValue
			err = nil
		}

		if err != nil {
			return t.escapeValue("{% elseif " + expressionStr + "!!" + err.Error() + " %}"), nil
		}

		value, err = t.applyfilters(value, filterParts, filters, data)
		if err != nil {
			return t.escapeValue("{% elseif " + expressionStr + "!!" + err.Error() + " %}"), nil
		}

		if toBool(value) {
			output, err := t.renderChildren(node, data, filters)
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
func (t *Template) renderElseNode(node *TreeNode, ifNodes []*TreeNode, data map[string]any, filters map[string]any) (string, error) {
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
		output, err := t.renderChildren(node, data, filters)
		if err != nil {
			return "", err
		}
		result += output
	}

	return result, nil
}

// renderForNode renders a 'for' loop node
func (t *Template) renderForNode(node *TreeNode, data map[string]any, filters map[string]any) (string, error) {
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

	value, err = t.applyfilters(value, filterParts, filters, data)
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
		output, err := t.renderChildren(node, newData, filters)
		if err != nil {
			return "", err
		}
		result += output
	}

	return result, nil
}

// renderVarNode renders a variable interpolation node
func (t *Template) renderVarNode(node *TreeNode, data map[string]any, filters map[string]any) (string, error) {
	expressionStr := node.Expression

	// Preprocess "is" tests
	exprPart, testFilter := processIsTests(expressionStr)

	parts := t.explodeRespectingQuotes("|", exprPart, -1)
	actualExpr := parts[0]
	filterParts := parts[1:]

	// Add test filter if present
	if testFilter != "" {
		filterParts = append(filterParts, testFilter)
	}

	expr := NewExpression(actualExpr)
	value, err := expr.Evaluate(data, t.resolvePath)
	if err != nil {
		return t.escapeValue("{{" + expressionStr + "!!" + err.Error() + "}}"), nil
	}

	value, err = t.applyfilters(value, filterParts, filters, data)
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

// applyfilters applies a chain of filter filters to a value
func (t *Template) applyfilters(value any, parts []string, filters map[string]any, data map[string]any) (any, error) {
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
				} else if argStr == "true" {
					// Boolean literal - true
					arguments = append(arguments, true)
				} else if argStr == "false" {
					// Boolean literal - false
					arguments = append(arguments, false)
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
		if fn, exists := filters[funcName]; exists {
			result, err := callFunction(fn, allArgs)
			if err != nil {
				return nil, err
			}
			value = result
		} else {
			return nil, fmt.Errorf("filter `%s` not found", funcName)
		}
	}

	return value, nil
}

// renderIncludeNode renders an 'include' node by loading and rendering another template
func (t *Template) renderIncludeNode(node *TreeNode, data map[string]any, filters map[string]any) (string, error) {
	if t.loader == nil {
		return "", fmt.Errorf("template loader not configured for include directive")
	}

	// Get the template name from include expression
	templateName := strings.Trim(node.Expression, "'\"")

	// Load the included template
	templateContent, err := t.loader(templateName)
	if err != nil {
		return "", fmt.Errorf("failed to load included template '%s': %v", templateName, err)
	}

	// Parse and render the included template
	tokens := t.tokenize(templateContent)
	tree := t.createSyntaxTree(tokens)

	// Render the included template with the same data and filters
	return t.renderChildren(tree, data, filters)
}
