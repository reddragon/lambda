package lang

import (
	"errors"
	"fmt"
	"strings"
)

// An AstNode either has a value, or has children.
// isValue = 1, if its a value, otherwise false.
type ASTNode struct {
	isValue  bool
	value    string
	children []*ASTNode
}

const (
	openBracket   string = "("
	closedBracket string = ")"
)

func errStr(expected, found string) error {
	return errors.New(fmt.Sprintf("Expected %s, got %s.", expected, found))
}

// This method gets you the AST of a given expression.
func getAST(exp string) (*ASTNode, []string, error) {
	tokens := tokenize(exp)
	if len(tokens) == 0 {
		return nil, nil, errors.New("Nothing to evaluate")
	}
	return getASTOfTokens(tokens)
}

func tokenize(exp string) []string {
	tokens := strings.Split(
		strings.Replace(strings.Replace(exp, "(", " ( ", -1), ")", " ) ", -1),
		" ",
	)
	noWSTokens := make([]string, 0)
	for _, token := range tokens {
		if token != "" {
			noWSTokens = append(noWSTokens, token)
		}
	}
	return noWSTokens
}

// This method does the heavy-lifting of building an AST, once an expression
// is tokenized.
func getASTOfTokens(tokens []string) (*ASTNode, []string, error) {
	var token = ""
	tokensLen := len(tokens)
	// If it is an empty list of tokens, the AST is a nil node
	if tokensLen == 0 {
		return nil, tokens, nil
	} else if tokensLen == 1 {
		// TODO Check that this token is a value.
		//      A proxy for now is checking if this is not a ( or )
		token, tokens = pop(tokens)
		if token == openBracket || token == closedBracket {
			return nil, tokens, errStr("value", token)
		}

		node := new(ASTNode)
		node.isValue = true
		node.value = token
		node.children = nil
		return node, tokens, nil
	} else {
		token, tokens = pop(tokens)
		if token != openBracket {
			return nil, tokens, errStr(openBracket, token)
		}

		node := new(ASTNode)
		node.isValue = false
		// Create a slice with 0 length initially.
		node.children = make([]*ASTNode, 0)

		tokensLen = len(tokens)
		for len(tokens) != 0 && tokens[0] != closedBracket {
			var childNode *ASTNode = nil
			var err error = nil
			// If this is not an open brace, this is a single value
			if tokens[0] != openBracket {
				token, tokens = pop(tokens)
				childNode, _, err = getASTOfTokens([]string{token})
			} else {
				childNode, tokens, err = getASTOfTokens(tokens)
			}
			if err != nil {
				return nil, tokens, err
			}
			node.children = append(node.children, childNode)
		}
		if len(tokens) == 0 {
			return nil, tokens, errStr(closedBracket, "nil")
		}

		token, tokens = pop(tokens)
		if token != closedBracket {
			return nil, tokens, errStr(token, closedBracket)
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
