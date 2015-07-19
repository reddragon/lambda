package lang

import (
	"errors"
	"fmt"
)

// An AstNode either has a value, or has children.
// isValue = 1, if its a value, otherwise false.
type ASTNode struct {
	isValue  bool
	value    string
	children []*ASTNode
}

const (
	oB string = "("
	cB string = ")"
)

func errStr(expected, found string) error {
	return errors.New(fmt.Sprintf("Expected %s, got %s.", expected, found))
}

// This method gets you the AST of the expression.
func GetAST(tokens []string) (*ASTNode, []string, error) {
	var token = ""
	tokensLen := len(tokens)
	// fmt.Printf("Processing tokens: %q, len: %d\n", tokens, tokensLen)
	// If it is an empty list of tokens, the AST is a nil node
	if tokensLen == 0 {
		return nil, tokens, nil
	} else if tokensLen == 1 {
		// TODO Check that this token is a value.
		//      A proxy for now is checking if this is not a ( or )
		token, tokens = pop(tokens)
		if token == oB || token == cB {
			return nil, tokens, errStr("value", token)
		}

		node := new(ASTNode)
		node.isValue = true
		node.value = token
		node.children = nil
		return node, tokens, nil
	} else {
		token, tokens = pop(tokens)
		if token != oB {
			return nil, tokens, errStr(oB, token)
		}

		node := new(ASTNode)
		node.isValue = false
		// Create a slice with 0 length initially.
		node.children = make([]*ASTNode, 0)

		tokensLen = len(tokens)
		for len(tokens) != 0 && tokens[0] != cB {
			var childNode *ASTNode = nil
			var err error = nil
			// If this is not an open brace, this is a single value
			if tokens[0] != oB {
				token, tokens = pop(tokens)
				// fmt.Printf("Processing val: %s\n", token)
				childNode, _, err = GetAST([]string{token})
			} else {
				// fmt.Printf("Processing inner AST: %q\n", tokens)
				childNode, tokens, err = GetAST(tokens)
			}
			if err != nil {
				return nil, tokens, err
			}
			node.children = append(node.children, childNode)
		}
		// fmt.Printf("Done with inner operands, tokens: %q\n", tokens)
		if len(tokens) == 0 {
			return nil, tokens, errStr(cB, "nil")
		}

		token, tokens = pop(tokens)
		if token != cB {
			return nil, tokens, errStr(token, cB)
		}
		return node, tokens, nil
	}
}

func StringifyAST(node *ASTNode) string {
	if node == nil {
		return ""
	}
	if node.isValue {
		return node.value
	}
	r := "("
	if node.children != nil {
		for i := 0; i < len(node.children); i++ {
			r += StringifyAST(node.children[i])
			if i+1 < len(node.children) {
				r += " "
			}
		}
	}
	r += ")"
	return r
}
