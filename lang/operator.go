package lang

import (
	"bytes"
	"errors"
	"fmt"
	"strings"
)

type Operator struct {
	symbol      string
	minArgCount int
	maxArgCount int
	doNotResolveVars bool
	handler     (func(*LangEnv, []Atom) Atom)
}

const (
	// Operators
	add string = "+"
	sub string = "-"
	mul string = "*"
	div string = "/"
	def string = "defvar"
	eq  string = "eq"
	gt  string = ">"
	geq string = ">="
	lt  string = "<"
	leq string = "<="
	and string = "and"
	or  string = "or"
)

func addOperator(opMap map[string]*Operator, op *Operator) {
	opMap[op.symbol] = op
}

func addBuiltinOperators(opMap map[string]*Operator) {
	numValPrecedenceMap := map[valueType]int{intType: 1, floatType: 2}
	strValPrecedenceMap := map[valueType]int{stringType: 1}

	addOperator(opMap,
		&Operator{
			symbol:      add,
			minArgCount: 2,
			maxArgCount: 100,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = chainedTypeCoerce(add, &operands, []map[valueType]int{numValPrecedenceMap, strValPrecedenceMap})
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

				case stringType:
					var buffer bytes.Buffer
					var finalVal stringValue
					for _, o := range operands {
						v, ok := o.Val.(*stringValue)
						if ok {
							buffer.WriteString(strings.Split(v.value, "\"")[1])
						}
					}

					retVal.Val = finalVal.newValue(fmt.Sprintf("\"%s\"", buffer.String()))
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      sub,
			minArgCount: 2,
			maxArgCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = typeCoerce(sub, &operands, numValPrecedenceMap)
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
			symbol:      mul,
			minArgCount: 2,
			maxArgCount: 100,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = typeCoerce(mul, &operands, numValPrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var finalVal intValue
					finalVal.value = 1
					for _, o := range operands {
						v, ok := o.Val.(*intValue)
						if ok {
							finalVal.value = finalVal.value * v.value
						} else {
							fmt.Errorf("Error while converting %s to intValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break

				case floatType:
					var finalVal floatValue
					finalVal.value = 1
					for _, o := range operands {
						v, ok := o.Val.(*floatValue)
						if ok {
							finalVal.value = finalVal.value * v.value
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
			symbol:      div,
			minArgCount: 2,
			maxArgCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = typeCoerce(div, &operands, numValPrecedenceMap)
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
			symbol:      def,
			minArgCount: 2,
			maxArgCount: 2,
			doNotResolveVars: true,
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

	addOperator(opMap,
		&Operator{
			symbol:      eq,
			minArgCount: 2,
			maxArgCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				vtype1 := operands[0].Val.getValueType()
				vtype2 := operands[1].Val.getValueType()

				if vtype1 != vtype2 {
					retVal.Err = errors.New(fmt.Sprintf("Cannot use %s operator for two different types %s and %s", eq, vtype1, vtype2))
					return retVal
				}

				retVal.Val = newBoolValue(operands[0].Val.Str() == operands[1].Val.Str())
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      gt,
			minArgCount: 2,
			maxArgCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = chainedTypeCoerce(gt, &operands, []map[valueType]int{numValPrecedenceMap, strValPrecedenceMap})
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var val1, val2 *intValue
					val1, _  = operands[0].Val.(*intValue)
					val2, _ = operands[1].Val.(*intValue)
					retVal.Val = newBoolValue(val1.value > val2.value)
					break

				case floatType:
					var val1, val2 *floatValue
					val1, _  = operands[0].Val.(*floatValue)
					val2, _ = operands[1].Val.(*floatValue)
					retVal.Val = newBoolValue(val1.value > val2.value)
					break

				case stringType:
					var val1, val2 *stringValue
					val1, _  = operands[0].Val.(*stringValue)
					val2, _ = operands[1].Val.(*stringValue)
					retVal.Val = newBoolValue(val1.value > val2.value)
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      geq,
			minArgCount: 2,
			maxArgCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = chainedTypeCoerce(gt, &operands, []map[valueType]int{numValPrecedenceMap, strValPrecedenceMap})
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var val1, val2 *intValue
					val1, _  = operands[0].Val.(*intValue)
					val2, _ = operands[1].Val.(*intValue)
					retVal.Val = newBoolValue(val1.value >= val2.value)
					break

				case floatType:
					var val1, val2 *floatValue
					val1, _  = operands[0].Val.(*floatValue)
					val2, _ = operands[1].Val.(*floatValue)
					retVal.Val = newBoolValue(val1.value >= val2.value)
					break

				case stringType:
					var val1, val2 *stringValue
					val1, _  = operands[0].Val.(*stringValue)
					val2, _ = operands[1].Val.(*stringValue)
					retVal.Val = newBoolValue(val1.value >= val2.value)
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      lt,
			minArgCount: 2,
			maxArgCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = chainedTypeCoerce(gt, &operands, []map[valueType]int{numValPrecedenceMap, strValPrecedenceMap})
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var val1, val2 *intValue
					val1, _  = operands[0].Val.(*intValue)
					val2, _ = operands[1].Val.(*intValue)
					retVal.Val = newBoolValue(val1.value < val2.value)
					break

				case floatType:
					var val1, val2 *floatValue
					val1, _  = operands[0].Val.(*floatValue)
					val2, _ = operands[1].Val.(*floatValue)
					retVal.Val = newBoolValue(val1.value < val2.value)
					break

				case stringType:
					var val1, val2 *stringValue
					val1, _  = operands[0].Val.(*stringValue)
					val2, _ = operands[1].Val.(*stringValue)
					retVal.Val = newBoolValue(val1.value < val2.value)
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      leq,
			minArgCount: 2,
			maxArgCount: 2,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				var finalType valueType
				finalType, retVal.Err = chainedTypeCoerce(gt, &operands, []map[valueType]int{numValPrecedenceMap, strValPrecedenceMap})
				if retVal.Err != nil {
					return retVal
				}

				switch finalType {
				case intType:
					var val1, val2 *intValue
					val1, _  = operands[0].Val.(*intValue)
					val2, _ = operands[1].Val.(*intValue)
					retVal.Val = newBoolValue(val1.value <= val2.value)
					break

				case floatType:
					var val1, val2 *floatValue
					val1, _  = operands[0].Val.(*floatValue)
					val2, _ = operands[1].Val.(*floatValue)
					retVal.Val = newBoolValue(val1.value <= val2.value)
					break

				case stringType:
					var val1, val2 *stringValue
					val1, _  = operands[0].Val.(*stringValue)
					val2, _ = operands[1].Val.(*stringValue)
					retVal.Val = newBoolValue(val1.value <= val2.value)
					break
				}
				return retVal
			},
		},
	)
}
