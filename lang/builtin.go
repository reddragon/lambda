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
func typeCoerce(operatorName string, operands *[]Atom, expectedArgs int, typePrecendenceMap map[ValueType]int) (ValueType, error) {
	var err error
	err = checkArgLen(Add, operands, 2)
	if err != nil {
		return "", err
	}

	var typesCountMap map[ValueType]int
	typesCountMap, err = checkArgTypes(operatorName, operands, []ValueType{IntType, FloatType})
	if err != nil {
		return "", err
	}

	if len(typesCountMap) == 1 {
		valType := (*operands)[0].Val.GetValueType()
		return valType, nil
	}

	var finalType ValueType
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
		if (*operands)[i].Val.GetValueType() != finalType {
			var err error
			(*operands)[i].Val, err = (*operands)[i].Val.To(finalType)

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
			symbol:   Add,
			argCount: 2,
			handler: func(operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[ValueType]int{IntType: 1, FloatType: 2}
				var finalType ValueType
				finalType, retVal.Err = typeCoerce(Add, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case IntType:
					var finalVal IntValue
					finalVal.value = 0
					for _, o := range operands {
						v, ok := o.Val.(*IntValue)
						if ok {
							finalVal.value = finalVal.value + v.value
						} else {
							fmt.Errorf("Error while converting %s to IntValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break

				case FloatType:
					var finalVal FloatValue
					finalVal.value = 0
					for _, o := range operands {
						v, ok := o.Val.(*FloatValue)
						if ok {
							finalVal.value = finalVal.value + v.value
						} else {
							fmt.Errorf("Error while converting %s to FloatValue\n", o.Val.Str())
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
			symbol:   Sub,
			argCount: 2,
			handler: func(operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[ValueType]int{IntType: 1, FloatType: 2}
				var finalType ValueType
				finalType, retVal.Err = typeCoerce(Sub, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case IntType:
					var finalVal IntValue
					var val1, val2 *IntValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*IntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to IntValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*IntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to IntValue\n", operands[0].Val.Str())
					}
					finalVal.value = val1.value - val2.value
					retVal.Val = finalVal
					break

				case FloatType:
					var finalVal FloatValue
					var val1, val2 *FloatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*FloatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to FloatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*FloatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to FloatValue\n", operands[0].Val.Str())
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
			symbol:   Mul,
			argCount: 2,
			handler: func(operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[ValueType]int{IntType: 1, FloatType: 2}
				var finalType ValueType
				finalType, retVal.Err = typeCoerce(Mul, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case IntType:
					var finalVal IntValue
					var val1, val2 *IntValue
					finalVal.value = 1

					var ok bool
					val1, ok = operands[0].Val.(*IntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to IntValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*IntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to IntValue\n", operands[0].Val.Str())
					}
					finalVal.value = val1.value * val2.value
					retVal.Val = finalVal
					break

				case FloatType:
					var finalVal FloatValue
					var val1, val2 *FloatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*FloatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to FloatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*FloatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to FloatValue\n", operands[0].Val.Str())
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
			symbol:   Div,
			argCount: 2,
			handler: func(operands []Atom) Atom {
				var retVal Atom
				typePrecedenceMap := map[ValueType]int{IntType: 1, FloatType: 2}
				var finalType ValueType
				finalType, retVal.Err = typeCoerce(Div, &operands, 2, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case IntType:
					var finalVal IntValue
					var val1, val2 *IntValue
					finalVal.value = 1

					var ok bool
					val1, ok = operands[0].Val.(*IntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to IntValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*IntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to IntValue\n", operands[0].Val.Str())
					}
					if val2.value != 0 {
						finalVal.value = val1.value / val2.value
						retVal.Val = finalVal
					} else {
						retVal.Err = errors.New(fmt.Sprintf("Divide by zero"))
					}
					break

				case FloatType:
					var finalVal FloatValue
					var val1, val2 *FloatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(*FloatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to FloatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(*FloatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to FloatValue\n", operands[0].Val.Str())
					}
					if val2.value != 0 {
						finalVal.value = val1.value / val2.value
						retVal.Val = finalVal
					} else {
						retVal.Err = errors.New(fmt.Sprintf("Divide by zero"))
					}
					break
				}
				return retVal
			},
		},
	)
	return opMap
}

func builtinTypes() []Value {
	types := make([]Value, 0)
	types = append(types, new(StringValue))
	types = append(types, new(IntValue))
	types = append(types, new(FloatValue))
	return types
}
