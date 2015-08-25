package lang

import (
	"errors"
	"fmt"
)

// An Atom is either a value, or an error
type Atom struct {
	Err error
	Val Value
}

func Eval(exp string, env *LangEnv) *Atom {
	astNode, _, err := getAST(exp)
	if err != nil {
		retVal := new(Atom)
		retVal.Err = err
		return retVal
	}

	// This round-about way is not necessary. Gomobile trips if you return
	// structs by value.
	// TODO
	// Remove this hack, pending: https://github.com/golang/go/issues/11318
	retVal := new(Atom)
	result := evalAST(env, astNode)
	retVal.Val = result.Val
	retVal.Err = result.Err
	return retVal
}

func evalAST(env *LangEnv, node *ASTNode) Atom {
	var retVal Atom
	retVal.Err = nil

	if node.isValue {
		value, err := getValue(node.value)
		if err != nil {
			retVal.Err = errStr("value", node.value)
		} else {
			retVal.Val = value
		}
		return retVal
	}
	if len(node.children) == 1 {
		return evalAST(env, node.children[0])
	}

	// Assuming that the first child is an operand
	symbol := node.children[0].value
	operator := env.getOperator(symbol)
	if operator == nil {
		retVal.Err = errors.New(fmt.Sprintf("Unknown operator '%s'", symbol))
		return retVal
	}

	if len(node.children)-1 != operator.argCount {
		retVal.Err = errors.New(
			fmt.Sprintf("Received %d arguments for operator %s, expected: %d",
				len(node.children)-1, symbol, operator.argCount))
		return retVal
	}

	operands := make([]Atom, 0)
	for i := 1; i < len(node.children); i++ {
		v := evalAST(env, node.children[i])
		if v.Err != nil {
			return v
		}
		operands = append(operands, v)
	}
	v := operator.handler(operands)
	if v.Err != nil {
		return v
	}
	retVal.Val = v.Val
	return retVal
}
