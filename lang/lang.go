package lang

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

// An Atom is either a value, or an error
type Atom struct {
	Err error
	Val Value
}

type EvalResult struct {
	ValStr          string
	ErrStr          string
	RemainingTokens string
}

func Eval(exp string, env *LangEnv) *EvalResult {
	exp = strings.TrimSpace(exp)
	evalResult := new(EvalResult)
	astNode, tokens, err := getAST(exp)
	if err != nil {
		evalResult.ErrStr = err.Error()
		return evalResult
	}

	// This round-about way is not necessary. Gomobile trips if you return
	// structs by value.
	// TODO
	// Remove this hack, pending: https://github.com/golang/go/issues/11318
	result := evalAST(env, astNode)

	if result.Val != nil && result.Val.getValueType() == varType {
		result.Val, result.Err = getVarValue(env, result.Val)
	}

	if result.Err != nil {
		evalResult.ErrStr = result.Err.Error()
	} else if result.Val != nil {
		evalResult.ValStr = result.Val.Str()
	}

	if tokens != nil && len(tokens) > 0 {
		var buffer bytes.Buffer
		for _, s := range tokens {
			buffer.WriteString(s)
			buffer.WriteString(" ")
		}
		evalResult.RemainingTokens = strings.TrimSpace(buffer.String())
	}
	return evalResult
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
	if len(node.children) == 0 {
		retVal.Err = errors.New("Cannot evaluate an empty expression")
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
		o.Val = newASTValue(node)
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
