package lang

import (
	// "bufio"
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

// TODO Use only err as error string
type Atom struct {
	Err    bool
	ErrMsg string
	Val    int
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

func add(operand1, operand2 Atom) Atom {
	var retVal Atom
	retVal.Val = operand1.Val + operand2.Val
	return retVal
}

func checkArgLen(operatorName string, operands []Atom, expectedArgs int) Atom {
	var retVal Atom
	if len(operands) != expectedArgs {
		var retVal Atom
		retVal.Err = true
		retVal.ErrMsg = fmt.Sprintf("For operator %s, expected %d args, but got %d.", operatorName, expectedArgs, len(operands))
		return retVal
	}
	return retVal
}

func getHandler(operator string) (int, func([]Atom) Atom) {
	if operator == Add {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Add, operands, 2)
			if argCheckErr.Err {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Val = operands[0].Val + operands[1].Val
			return retVal
		}
	} else if operator == Sub {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Sub, operands, 2)
			if argCheckErr.Err {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Val = operands[0].Val - operands[1].Val
			return retVal
		}
	} else if operator == Mul {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Mul, operands, 2)
			if argCheckErr.Err {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Val = operands[0].Val * operands[1].Val
			return retVal
		}
	} else if operator == Div {
		return 2, func(operands []Atom) Atom {
			var retVal Atom
			argCheckErr := checkArgLen(Div, operands, 2)
			if argCheckErr.Err {
				return argCheckErr
			}
			// TODO
			// Do type checks

			retVal.Val = operands[0].Val / operands[1].Val
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

func evalArgs(tokens []string) (Atom, []string) {
	var retVal Atom
	retVal.Err = false
	retVal.ErrMsg = ""

	// Expecting a value, but if a ')' occurs, its an Error
	if tokens[0] == ClosedBracket {
		retVal.ErrMsg = "Unexpected ')'"
		retVal.Err = true
		return retVal, tokens
	}

	// This could be a nested expression
	if tokens[0] == OpenBracket {
		return Eval(tokens)
	}

	// This is an atomic value
	token := ""
	token, tokens = pop(tokens)
	val, err := strconv.Atoi(token)
	if err != nil {
		retVal.ErrMsg = "Problem while converting " + token + ": " + err.Error()
		retVal.Err = true
		return retVal, tokens
	}
	retVal.Val = val
	return retVal, tokens
}

func Eval(tokens []string) (Atom, []string) {
	var retVal Atom
	retVal.Err = false
	retVal.ErrMsg = ""

	token := ""
	token, tokens = pop(tokens)
	if token != OpenBracket {
		// TODO
		// Raise an exception
		retVal.ErrMsg = "Expected a '('"
		retVal.Err = true
		return retVal, tokens
	}

	token, tokens = pop(tokens)
	argCount, handler := getHandler(token)
	if handler == nil {
		retVal.ErrMsg = "Invalid operator, " + token
		retVal.Err = true
		return retVal, tokens
	}

	// Read two args.
	// TODO change this to read as many args as the operator wants

	// TODO
	// Reading just two ints/floats here
	// Change to read two nested expressions here too
	operands := make([]Atom, 0)
	for i := 0; i < argCount; i++ {
		var operand Atom
		operand, tokens = evalArgs(tokens)
		if operand.Err {
			return operand, tokens
		}
		operands = append(operands, operand)
	}

	token, tokens = pop(tokens)
	if token != ClosedBracket {
		retVal.ErrMsg = "Expected a ')'"
		retVal.Err = true
		return retVal, tokens
	}

	return handler(operands), tokens
}
