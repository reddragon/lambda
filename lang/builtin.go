package lang

import (
	"fmt"
)

func addOperator(opMap map[string]*Operator, op *Operator) {
	opMap[op.symbol] = op
}

// Algorithm:
// 1. We already know the type -> count mapping.
// 2. If there is only one type, there is nothing to do.
// 3. If there are multiple, pick the one with the highest precedence.
// 4. Try and cast all operand values to that type. Error out if any of them
//    resists. Because, resistance is futile.
//
// TODO
// Mix typeCoerce() and checkArgTypes()
func typeCoerce(operands *[]Atom, typesCountMap map[ValueType]int, typePrecendenceMap map[ValueType]int) (ValueType, error) {
	if len(typesCountMap) == 1 {
		fmt.Printf("Value initially was: %s\n", (*operands)[0].Val.Str())
		valType := (*operands)[0].Val.GetValueType()
		fmt.Printf("Only one type (%s) in the typesCountMap, so nothing to do.\n", valType)
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

	for _, o := range *operands {
		if o.Val.GetValueType() != finalType {
			var err error
			o.Val, err = o.Val.To(finalType)
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
				retVal.Err = checkArgLen(Add, operands, 2)
				if retVal.Err != nil {
					return retVal
				}

				var typesCountMap map[ValueType]int
				typesCountMap, retVal.Err = checkArgTypes(Add, operands, []ValueType{IntType, FloatType})
				if retVal.Err != nil {
					return retVal
				}

				// Type Coercion
				typePrecedenceMap := map[ValueType]int{IntType: 1, FloatType: 2}
				var finalType ValueType
				finalType, retVal.Err = typeCoerce(&operands, typesCountMap, typePrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}
				fmt.Printf("Did type-coercion, final type was: %s\n", finalType)

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
					break;

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
					break;
				}
				return retVal
			},
		},
	)

	/*
		addOperator(opMap,
			&Operator{
				symbol:   Sub,
				argCount: 2,
				handler: func(operands []Atom) Atom {
					var retVal Atom
					argCheckErr := checkArgLen(Sub, operands, 2)
					if argCheckErr.Err != nil {
						return argCheckErr
					}
					// TODO
					// Do type checks

					retVal.Val = operands[0].Val - operands[1].Val
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
					argCheckErr := checkArgLen(Mul, operands, 2)
					if argCheckErr.Err != nil {
						return argCheckErr
					}
					// TODO
					// Do type checks

					retVal.Val = operands[0].Val * operands[1].Val
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
					argCheckErr := checkArgLen(Div, operands, 2)
					if argCheckErr.Err != nil {
						return argCheckErr
					}
					// TODO
					// Do type checks

					retVal.Val = operands[0].Val / operands[1].Val
					return retVal
				},
			},
		) */
	return opMap
}

func builtinTypes() []Value {
	types := make([]Value, 0)
	types = append(types, new(StringValue))
	types = append(types, new(IntValue))
	types = append(types, new(FloatValue))
	return types
}
