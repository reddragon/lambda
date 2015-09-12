package lang

import (
	"errors"
	"fmt"
)

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
					operatorName, operand.Val.Str(), operand.Val.getValueType()))
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
