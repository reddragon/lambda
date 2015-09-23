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

func Eval(exp string, env *LangEnv) (*Atom, []string) {
	astNode, tokens, err := getAST(exp)
	if err != nil {
		retVal := new(Atom)
		retVal.Err = err
		return retVal, nil
	}

	// This round-about way is not necessary. Gomobile trips if you return
	// structs by value.
	// TODO
	// Remove this hack, pending: https://github.com/golang/go/issues/11318
	retVal := new(Atom)
	result := evalAST(env, astNode)
	retVal.Val = result.Val
	retVal.Err = result.Err
	return retVal, tokens
}

func evalAST(env *LangEnv, node *ASTNode) Atom {
	var retVal Atom
	retVal.Err = nil

	if node.isValue {
		value, err := getValue(env, node.value)
		if err != nil {
			retVal.Err = errors.New(fmt.Sprintf("%s %s", errStr("value", node.value), err))
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

	if operator.minArgCount == operator.maxArgCount {
		argCount := operator.minArgCount
		if len(node.children)-1 != argCount {
			retVal.Err = errors.New(
				fmt.Sprintf("Received %d arguments for operator %s, expected: %d",
					len(node.children)-1, symbol, argCount))
			return retVal
		}
	} else {
		if len(node.children)-1 < operator.minArgCount {
			retVal.Err = errors.New(
				fmt.Sprintf("Received %d arguments for operator %s, minimum expected arguments: %d",
					len(node.children)-1, symbol, operator.minArgCount))
			return retVal
		} else if len(node.children)-1 > operator.maxArgCount {
			retVal.Err = errors.New(
				fmt.Sprintf("Received %d arguments for operator %s, maximum expected arguments: %d",
					len(node.children)-1, symbol, operator.maxArgCount))
			return retVal
		}
	}

	operands := make([]Atom, 0)
	if operator.passRawAST {
		var o Atom
		o.Val = newASTValue(node.children[1:])
		operands = append(operands, o)
	} else {
		for i := 1; i < len(node.children); i++ {
			v := evalAST(env, node.children[i])
			if v.Err != nil {
				return v
			}
			if !operator.doNotResolveVars && v.Val.getValueType() == varType {
				v.Val, v.Err = getVarValue(env, v.Val)
				if v.Err != nil {
					return v
				}
			}
			operands = append(operands, v)
		}
	}
	v := operator.handler(env, operands)
	if v.Err != nil {
		return v
	}
	retVal.Val = v.Val
	return retVal
}
