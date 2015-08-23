package lang

import (
	"errors"
	"fmt"
)

func checkArgLen(operatorName string, operands *[]Atom, expectedArgs int) error {
	if len(*operands) != expectedArgs {
		return errors.New(
			fmt.Sprintf("For operator %s, expected %d args, but got %d.",
				operatorName, expectedArgs, len(*operands)))
	}
	return nil
}

func checkArgTypes(operatorName string, operands *[]Atom, allowedTypes []valueType) (map[valueType]int, error) {
	typesFound := make(map[valueType]int, 0)
	for _, operand := range *operands {
		exists := false
		typesFound[operand.Val.getValueType()]++
		for _, t := range allowedTypes {
			if operand.Val.getValueType() == t {
				exists = true
				break
			}
		}
		if !exists {
			return nil, errors.New(
				fmt.Sprintf("For operator %s, operand %s is of unexpected type: %s.",
					operatorName, operand, operand.Val.getValueType()))
		}
	}
	return typesFound, nil
}

func pop(tokens []string) (string, []string) {
	if len(tokens) == 0 {
		return "", tokens
	}
	return tokens[0], tokens[1:]
}
