package tqtemplate

import (
	"fmt"
	"strconv"
	"strings"
	"unicode"
	"unicode/utf8"
)

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
		// Decode the current rune properly for UTF-8
		ch, chSize := utf8.DecodeRuneInString(expr[i:])

		// Skip whitespace
		if unicode.IsSpace(ch) {
			i += chSize
			continue
		}

		// Handle parentheses
		if ch == '(' || ch == ')' {
			tokens = append(tokens, ExpressionToken{Type: "parenthesis", Value: string(ch)})
			i += chSize
			continue
		}

		// Handle word-based operators (and, or, not) - only ASCII letters
		if ch < utf8.RuneSelf && unicode.IsLetter(ch) {
			word := ""
			start := i
			for i < length {
				r, size := utf8.DecodeRuneInString(expr[i:])
				if r < utf8.RuneSelf && unicode.IsLetter(r) {
					word += string(r)
					i += size
				} else {
					break
				}
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
			i += chSize
			continue
		}

		// Handle numbers
		if unicode.IsDigit(ch) || (ch == '.' && i < length-1 && unicode.IsDigit(rune(expr[i+1]))) {
			num := ""
			for i < length {
				r, size := utf8.DecodeRuneInString(expr[i:])
				if unicode.IsDigit(r) || r == '.' {
					num += string(r)
					i += size
				} else {
					break
				}
			}
			tokens = append(tokens, ExpressionToken{Type: "number", Value: num})
			continue
		}

		// Handle string literals
		if ch == '"' {
			str := ""
			i += chSize // Skip opening quote
			escaped := false
			for i < length {
				r, size := utf8.DecodeRuneInString(expr[i:])
				if escaped {
					str += string(r)
					escaped = false
					i += size
				} else if r == '\\' {
					escaped = true
					i += size
				} else if r == '"' {
					i += size // Skip closing quote
					break
				} else {
					str += string(r)
					i += size
				}
			}
			tokens = append(tokens, ExpressionToken{Type: "string", Value: str})
			continue
		}

		// Handle identifiers/paths (with dots for nested access)
		if unicode.IsLetter(ch) || ch == '_' {
			ident := ""
			for i < length {
				r, size := utf8.DecodeRuneInString(expr[i:])
				if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' || r == '.' {
					ident += string(r)
					i += size
				} else {
					break
				}
			}
			tokens = append(tokens, ExpressionToken{Type: "identifier", Value: ident})
			continue
		}

		// Unknown character, skip it
		i += chSize
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
