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

// This method force casts all operands to the given value type.
func tryTypeCastTo(operands *[]Atom, finalType valueType) error {
	for i := 0; i < len(*operands); i++ {
		if (*operands)[i].Val.getValueType() != finalType {
			var err error
			(*operands)[i].Val, err = (*operands)[i].Val.to(finalType)

			if err != nil {
				return err
			}
		}
	}
	return nil
}

// Algorithm:
// 1. We get the type -> count mapping.
// 2. If there is only one type, there is nothing to do.
// 3. If there are multiple, pick the one with the highest precedence.
// 4. Try and cast all operand values to that type. Error out if any of them
//    resists. Because, resistance is futile.
func typeCoerce(operatorName string, operands *[]Atom, typePrecendenceMap map[valueType]int) (valueType, error) {
	var err error
	var typesCountMap map[valueType]int
	allowedTypes := make([]valueType, len(typePrecendenceMap))
	typeIdx := 0
	for vtype, _ := range typePrecendenceMap {
		allowedTypes[typeIdx] = vtype
		typeIdx++
	}

	typesCountMap, err = checkArgTypes(operatorName, operands, allowedTypes)
	if err != nil {
		return "", err
	}

	if len(typesCountMap) == 1 {
		valType := (*operands)[0].Val.getValueType()
		return valType, nil
	}

	var finalType valueType
	finalTypePrecedence := -1
	for t, c := range typesCountMap {
		if c <= 0 {
			continue
		}
		precedence := typePrecendenceMap[t]
		if precedence > finalTypePrecedence {
			finalType = t
			finalTypePrecedence = precedence
		}
	}

	err = tryTypeCastTo(operands, finalType)
	if err != nil {
		return "", err
	}
	return finalType, nil
}

// This is used when a group of value types can be used with the operator.
// This method processes the operands in order of the precendence maps.
func chainedTypeCoerce(operatorName string, operands *[]Atom, typePrecendenceMapList []map[valueType]int) (valueType, error) {
	var vtype valueType
	var err error
	for _, pMap := range typePrecendenceMapList {
		vtype, err = typeCoerce(operatorName, operands, pMap)
		if err == nil {
			return vtype, nil
		}
	}
	return "", err
}

func getASTStr(node *ASTNode) string {
	if node != nil {
		if node.isValue {
			return node.value
		}

		if node.children != nil {
			astStr := ""
			for _, child := range node.children {
				astStr = astStr + getASTStr(child) + " "
			}
			return astStr
		}
	}
	return ""
}

func printVarMap(varMap map[string]Value) {
	fmt.Println("varMap values: ")
	for k, v := range varMap {
		fmt.Printf("%s = %s\n", k, v.Str())
	}
}
