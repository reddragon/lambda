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

type Atom struct {
	Err   error
	Value int
}

const (
	// Delimiters
	OpenBracket   string = "("
	ClosedBracket string = ")"

	// Operators
	Add string = "+"
	Sub string = "-"
	Mul string = "*"
	Div string = "/"
)

func checkArgLen(operatorName string, operands []Atom, expectedArgs int) Atom {
	var retVal Atom
	if len(operands) != expectedArgs {
		var retVal Atom

		retVal.Err = errors.New(fmt.Sprintf("For operator %s, expected %d args, but got %d.", operatorName, expectedArgs, len(operands)))
		return retVal
	}
	return retVal
}

func getHandler(operator string) (int, func([]Atom) Atom) {
	if operator == Add {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Add, operands, 2)
			if argCheckErr.Err != nil {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Value = operands[0].Value + operands[1].Value
			return retVal
		}
	} else if operator == Sub {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Sub, operands, 2)
			if argCheckErr.Err != nil {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Value = operands[0].Value - operands[1].Value
			return retVal
		}
	} else if operator == Mul {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Mul, operands, 2)
			if argCheckErr.Err != nil {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Value = operands[0].Value * operands[1].Value
			return retVal
		}
	} else if operator == Div {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Div, operands, 2)
			if argCheckErr.Err != nil {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Value = operands[0].Value / operands[1].Value
			return retVal
		}
	}
	return 0, nil
}

func pop(tokens []string) (string, []string) {
	if len(tokens) == 0 {
		return "", tokens
	}
	return tokens[0], tokens[1:]
}

func EvalAST(node *ASTNode) Atom {
	var retVal Atom
	retVal.Err = nil
	retVal.Value = 0

	if node.isValue {
		// TODO Check here if the value is a proper value
		value, err := strconv.Atoi(node.value)
		if err != nil {
			retVal.Err = errStr("value", node.value)
		} else {
			retVal.Value = value
		}
		return retVal
	}
	if len(node.children) == 1 {
		return EvalAST(node.children[0])
	}

	// Assuming that the first child is an operand
	operator := node.children[0].value
	// fmt.Printf("Will work on the operator %s\n", operator)
	argCount, handler := getHandler(operator)
	// fmt.Printf("ArgCount: %d\n", argCount)

	if len(node.children)-1 != argCount {
		retVal.Err = errors.New(
			fmt.Sprintf("Received %d arguments for operator %s, expected: %d",
				len(node.children)-1, operator, argCount))
		return retVal
	}

	operands := make([]Atom, 0)
	for i := 1; i < len(node.children); i++ {
		v := EvalAST(node.children[i])
		if v.Err != nil {
			return v
		}
		// fmt.Printf("Pushing value: %d\n", v.Value)
		operands = append(operands, v)
	}
	// fmt.Printf("Len of operands: %d\n", len(operands))
	v := handler(operands)
	if v.Err != nil {
		return v
	}
	retVal.Value = v.Value
	return retVal
}
