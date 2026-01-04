package tqtemplate

import (
	"fmt"
	"strings"
)

// findExtendsNode finds the extends node in the tree (should be first non-whitespace element)
func (t *Template) findExtendsNode(tree *TreeNode) *TreeNode {
	for _, child := range tree.Children {
		if child.Type == "extends" {
			return child
		}
		// Skip only whitespace literals
		if child.Type == "lit" && strings.TrimSpace(child.Expression) == "" {
			continue
		}
		// If we hit any other content, stop looking
		break
	}
	return nil
}

// renderWithExtends handles template inheritance
func (t *Template) renderWithExtends(childTree *TreeNode, extendsNode *TreeNode, data map[string]any, functions map[string]any) (string, error) {
	if t.loader == nil {
		return "", fmt.Errorf("template loader not configured for extends directive")
	}

	// Get the parent template name from extends expression
	parentName := strings.Trim(extendsNode.Expression, "'\"")

	// Load parent template
	parentContent, err := t.loader(parentName)
	if err != nil {
		return "", fmt.Errorf("failed to load parent template '%s': %v", parentName, err)
	}

	// Parse parent template
	parentTokens := t.tokenize(parentContent)
	parentTree := t.createSyntaxTree(parentTokens)

	// Collect blocks from child template
	childBlocks := t.collectBlocks(childTree)

	// Render parent with child blocks overriding
	return t.renderWithBlocks(parentTree, childBlocks, data, functions)
}

// collectBlocks extracts all block definitions from a template tree
func (t *Template) collectBlocks(tree *TreeNode) map[string]*TreeNode {
	blocks := make(map[string]*TreeNode)
	var walk func(*TreeNode)
	walk = func(node *TreeNode) {
		if node.Type == "block" {
			blocks[node.Expression] = node
		}
		for _, child := range node.Children {
			walk(child)
		}
	}
	walk(tree)
	return blocks
}

// renderWithBlocks renders a tree with block overrides
func (t *Template) renderWithBlocks(tree *TreeNode, blockOverrides map[string]*TreeNode, data map[string]any, functions map[string]any) (string, error) {
	result := ""
	ifNodes := []*TreeNode{}

	for i, child := range tree.Children {
		switch child.Type {
		case "block":
			// Check if this block is overridden
			blockName := child.Expression

			// Check if the previous node is a literal with only whitespace (no newlines)
			// to preserve indentation from parent
			precedingWhitespace := ""
			if i > 0 {
				prevNode := tree.Children[i-1]
				if prevNode.Type == "lit" {
					prevContent := prevNode.Expression
					// Check if it's whitespace without newlines
					if strings.TrimSpace(prevContent) == "" && !strings.Contains(prevContent, "\n") && !strings.Contains(prevContent, "\r") {
						precedingWhitespace = prevContent
					}
				}
			}

			if override, exists := blockOverrides[blockName]; exists {
				// Add preceding whitespace before override content
				result += precedingWhitespace
				// Render the override block (with block overrides for nested blocks)
				output, err := t.renderWithBlocks(override, blockOverrides, data, functions)
				if err != nil {
					return "", err
				}
				result += output
			} else {
				// Render the default block content (with block overrides for nested blocks)
				output, err := t.renderWithBlocks(child, blockOverrides, data, functions)
				if err != nil {
					return "", err
				}
				result += output
			}
			ifNodes = []*TreeNode{}
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
			// Skip this literal if it's preceding whitespace for a block
			// (it's already been handled as part of the block rendering)
			if i+1 < len(tree.Children) && tree.Children[i+1].Type == "block" {
				// Check if this literal is whitespace-only without newlines
				if strings.TrimSpace(child.Expression) == "" && !strings.Contains(child.Expression, "\n") && !strings.Contains(child.Expression, "\r") {
					// This will be included with the block, so skip it here
					ifNodes = []*TreeNode{}
					continue
				}
			}
			result += child.Expression
			ifNodes = []*TreeNode{}
		}
	}

	return result, nil
}
