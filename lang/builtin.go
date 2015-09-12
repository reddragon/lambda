package lang

import (
	"errors"
	"fmt"
)

func addOperator(opMap map[string]*Operator, op *Operator) {
	opMap[op.symbol] = op
}

// Algorithm:
// 1. We get the type -> count mapping.
// 2. If there is only one type, there is nothing to do.
// 3. If there are multiple, pick the one with the highest precedence.
// 4. Try and cast all operand values to that type. Error out if any of them
//    resists. Because, resistance is futile.
func typeCoerce(operatorName string, operands *[]Atom, expectedArgs int, typePrecendenceMap map[valueType]int) (valueType, error) {
	var err error
	var typesCountMap map[valueType]int
	typesCountMap, err = checkArgTypes(operatorName, operands, []valueType{intType, floatType})
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

	for i := 0; i < len(*operands); i++ {
		if (*operands)[i].Val.getValueType() != finalType {
			var err error
			(*operands)[i].Val, err = (*operands)[i].Val.to(finalType)

			if err != nil {
				return "", err
			}
		}
	}
	return finalType, nil
}

// This methods returns all the builtin operators to the environment
func builtinOperators() map[string]*Operator {
	opMap := make(map[string]*Operator)

	addOperator(opMap,
		&Operator{
			symbol:   add,
			argCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[valueType]int{intType: 1, floatType: 2}
				var finalType valueType
				finalType, retVal.Err = typeCoerce(add, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var finalVal intValue
					finalVal.value = 0
					for _, o := range operands {
						v, ok := o.Val.(*intValue)
						if ok {
							finalVal.value = finalVal.value + v.value
						} else {
							fmt.Errorf("Error while converting %s to intValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break

				case floatType:
					var finalVal floatValue
					finalVal.value = 0
					for _, o := range operands {
						v, ok := o.Val.(*floatValue)
						if ok {
							finalVal.value = finalVal.value + v.value
						} else {
							fmt.Errorf("Error while converting %s to floatValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:   sub,
			argCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[valueType]int{intType: 1, floatType: 2}
				var finalType valueType
				finalType, retVal.Err = typeCoerce(sub, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var finalVal intValue
					var val1, val2 *intValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}
					finalVal.value = val1.value - val2.value
					retVal.Val = finalVal
					break

				case floatType:
					var finalVal floatValue
					var val1, val2 *floatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*floatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*floatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}
					finalVal.value = val1.value - val2.value
					retVal.Val = finalVal
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:   mul,
			argCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[valueType]int{intType: 1, floatType: 2}
				var finalType valueType
				finalType, retVal.Err = typeCoerce(mul, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var finalVal intValue
					var val1, val2 *intValue
					finalVal.value = 1

					var ok bool
					val1, ok = operands[0].Val.(*intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}
					finalVal.value = val1.value * val2.value
					retVal.Val = finalVal
					break

				case floatType:
					var finalVal floatValue
					var val1, val2 *floatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*floatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*floatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}
					finalVal.value = val1.value * val2.value
					retVal.Val = finalVal
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:   div,
			argCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[valueType]int{intType: 1, floatType: 2}
				var finalType valueType
				finalType, retVal.Err = typeCoerce(div, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var finalVal intValue
					var val1, val2 *intValue
					finalVal.value = 1

					var ok bool
					val1, ok = operands[0].Val.(*intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}
					if val2.value != 0 {
						finalVal.value = val1.value / val2.value
						retVal.Val = finalVal
					} else {
						retVal.Err = errors.New(fmt.Sprintf("divide by zero"))
					}
					break

				case floatType:
					var finalVal floatValue
					var val1, val2 *floatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*floatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*floatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}
					if val2.value != 0 {
						finalVal.value = val1.value / val2.value
						retVal.Val = finalVal
					} else {
						retVal.Err = errors.New(fmt.Sprintf("divide by zero"))
					}
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:   def,
			argCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				vtype1 := operands[0].Val.getValueType()
				vtype2 := operands[1].Val.getValueType()

				if vtype1 != varType {
					retVal.Err = errors.New(fmt.Sprintf("For %s, expected %s to be %s, but was %s", def, operands[0].Val.Str(), varType, vtype1))
					return retVal
				}

				if vtype2 == varType {
					retVal.Err = errors.New(fmt.Sprintf("For %s, expected %s to not be %s, but was.", def, operands[1].Val.Str(), varType))
					return retVal
				}

				sym := operands[0].Val.Str()
				if env.getOperator(sym) != nil {
					retVal.Err = errors.New(fmt.Sprintf("Cannot use %s as a variable, as it is defined as an operator.", sym))
					return retVal
				}

				env.varMap[sym] = operands[1].Val
				retVal.Val = operands[1].Val
				return retVal
			},
		},
	)
	return opMap
}

func builtinTypes() []Value {
	types := make([]Value, 0)
	// Append the values in the order of prefence, i.e, more specific types
	// should be first.
	types = append(types, new(stringValue))
	types = append(types, new(intValue))
	types = append(types, new(floatValue))
	types = append(types, new(boolValue))
	types = append(types, new(varValue))
	return types
}
