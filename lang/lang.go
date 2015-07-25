package lang

import (
	// "bufio"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

func Tokenize(line string) []string {
	tokens := strings.Split(strings.Replace(strings.Replace(line, "(", " ( ", -1), ")", " ) ", -1), " ")
	noWSTokens := make([]string, 0)
	for _, token := range tokens {
		if token != "" {
			noWSTokens = append(noWSTokens, token)
		}
	}
	return noWSTokens
}

// An Atom is either a value, or an error
type Atom struct {
	Err error
	// TODO Value should be of type Value
	Val int
}

// Different types of values supported
type ValueType int

const (
	// Value type
	IntType = iota
	FloatType
	StringType
)

type Value interface {
	valueType() ValueType
	to(ValueType) (Value, error)
	str() string
	// TODO
	// Add other methods
}

func checkArgLen(operatorName string, operands []Atom, expectedArgs int) Atom {
	var retVal Atom
	if len(operands) != expectedArgs {
		var retVal Atom

		retVal.Err = errors.New(
			fmt.Sprintf("For operator %s, expected %d args, but got %d.",
				operatorName, expectedArgs, len(operands)))
		return retVal
	}
	return retVal
}

func pop(tokens []string) (string, []string) {
	if len(tokens) == 0 {
		return "", tokens
	}
	return tokens[0], tokens[1:]
}

func EvalAST(env *LangEnv, node *ASTNode) Atom {
	var retVal Atom
	retVal.Err = nil
	retVal.Val = 0

	if node.isValue {
		// TODO Check here if the value is a proper value
		value, err := strconv.Atoi(node.value)
		if err != nil {
			retVal.Err = errStr("value", node.value)
		} else {
			retVal.Val = value
		}
		return retVal
	}
	if len(node.children) == 1 {
		return EvalAST(env, node.children[0])
	}

	// Assuming that the first child is an operand
	symbol := node.children[0].value
	// fmt.Printf("Will work on the operator %s\n", operator)
	operator := env.getOperator(symbol)
	if operator == nil {
		retVal.Err = errors.New(fmt.Sprintf("Unknown operator '%s'", symbol))
		return retVal
	}

	// fmt.Printf("ArgCount: %d\n", argCount)

	if len(node.children)-1 != operator.argCount {
		retVal.Err = errors.New(
			fmt.Sprintf("Received %d arguments for operator %s, expected: %d",
				len(node.children)-1, symbol, operator.argCount))
		return retVal
	}

	operands := make([]Atom, 0)
	for i := 1; i < len(node.children); i++ {
		v := EvalAST(env, node.children[i])
		if v.Err != nil {
			return v
		}
		// fmt.Printf("Pushing value: %d\n", v.Val)
		operands = append(operands, v)
	}
	// fmt.Printf("Len of operands: %d\n", len(operands))
	v := operator.handler(operands)
	if v.Err != nil {
		return v
	}
	retVal.Val = v.Val
	return retVal
}
