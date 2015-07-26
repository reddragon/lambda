package lang

import (
  "errors"
  "fmt"
)

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
