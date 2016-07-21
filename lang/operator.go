package lang

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"math/big"
	"strings"
)

type Operator struct {
	symbol           string
	minArgCount      int
	maxArgCount      int
	doNotResolveVars bool
	passRawAST       bool
	handler          (func(*LangEnv, []Atom) Atom)
}

const (
	// Operators
	add   string = "+"
	sub   string = "-"
	mul   string = "*"
	div   string = "/"
	def   string = "defvar"
	eq    string = "="
	gt    string = ">"
	geq   string = ">="
	lt    string = "<"
	leq   string = "<="
	and   string = "and"
	or    string = "or"
	defun string = "defun"
	cond  string = "cond"
)

func addOperator(opMap map[string]*Operator, op *Operator) {
	opMap[op.symbol] = op
}

func addBuiltinOperators(opMap map[string]*Operator) {
	numValPrecedenceMap := map[valueType]int{intType: 1, bigIntType: 2, floatType: 3}
	strValPrecedenceMap := map[valueType]int{stringType: 1}
	boolValPrecedenceMap := map[valueType]int{boolType: 1}

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

			performOp:
				switch finalType {
				case intType:
					var finalVal intValue
					finalVal.value = 0
					for _, o := range operands {
						v, ok := o.Val.(intValue)
						if ok {
							// Check for overflow/underflow here.
							if (v.value > 0 && (finalVal.value > math.MaxInt64-v.value)) ||
								(v.value <= 0 && finalVal.value < math.MinInt64-v.value) {
								// There will be an overflow, so better cast to bigIntType here.
								err := tryTypeCastTo(&operands, bigIntType)
								if err != nil {
									fmt.Printf("Problem while avoiding overflow in operand %s: %s.\n", add, err)
								} else {
									finalType = bigIntType
									goto performOp
								}
							}

							finalVal.value = finalVal.value + v.value
						} else {
							fmt.Errorf("Error while converting %s to intValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break

				case bigIntType:
					var finalVal bigIntValue
					finalVal.value = new(big.Int)
					for _, o := range operands {
						v, ok := o.Val.(bigIntValue)
						if ok {
							finalVal.value = finalVal.value.Add(finalVal.value, v.value)
						} else {
							fmt.Errorf("Error while converting %s to bigIntValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break

				case floatType:
					var finalVal floatValue
					finalVal.value = 0
					for _, o := range operands {
						v, ok := o.Val.(floatValue)
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
						v, ok := o.Val.(stringValue)
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

			performOp:
				switch finalType {
				case intType:
					var finalVal intValue
					var val1, val2 intValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[1].Val.Str())
					}

					// Check for overflow/underflow here.
					if (val2.value > 0 && val1.value < math.MinInt64+val2.value) || (val2.value <= 0 && val1.value > math.MaxInt64+val2.value) {
						err := tryTypeCastTo(&operands, bigIntType)
						if err != nil {
							fmt.Printf("Problem while avoiding overflow in operand %s: %s.\n", add, err)
						} else {
							finalType = bigIntType
							goto performOp
						}
					}

					finalVal.value = val1.value - val2.value
					retVal.Val = finalVal
					break

				case bigIntType:
					var finalVal bigIntValue
					var val1, val2 bigIntValue
					finalVal.value = new(big.Int)

					var ok bool
					val1, ok = operands[0].Val.(bigIntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to bigIntValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(bigIntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to bigIntValue\n", operands[1].Val.Str())
					}
					finalVal.value.Sub(val1.value, val2.value)
					retVal.Val = finalVal
					break

				case floatType:
					var finalVal floatValue
					var val1, val2 floatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(floatValue)
					if !ok {
						fmt.Errorf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(floatValue)
					if !ok {
						fmt.Printf("It was not ok!, type: %s, rawtype: %T\n",
							operands[1].Val.getValueType(), operands[1].Val)
						fmt.Errorf("Error while converting %s to floatValue\n", operands[1].Val.Str())
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

			performOp:
				switch finalType {
				case intType:
					var finalVal intValue
					finalVal.value = 1

					bigIntOp1 := new(big.Int)
					bigIntOp2 := new(big.Int)
					for _, o := range operands {
						v, ok := o.Val.(intValue)
						if ok {
							// Check for overflow/underflow here.
							// The check for multiply turns out to be quite complex.
							// Refer to this document for the suggested implementation:
							// https://www.securecoding.cert.org/confluence/display/java/NUM00-J.+Detect+or+prevent+integer+overflow
							// In order to preserve readability, I would just check if
							// big.Int's Mul method returns the same value.
							// This would obviously make multiplications a little slow.
							bigIntOp1.SetInt64(finalVal.value)
							bigIntOp2.SetInt64(v.value)
							bigIntOp1.Mul(bigIntOp1, bigIntOp2)

							finalVal.value = finalVal.value * v.value
							if bigIntOp1.String() != finalVal.Str() {
								err := tryTypeCastTo(&operands, bigIntType)
								if err != nil {
									fmt.Printf("Problem while avoiding overflow in operand %s: %s.\n", add, err)
								} else {
									finalType = bigIntType
									goto performOp
								}
							}
						} else {
							fmt.Errorf("Error while converting %s to intValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break

				case bigIntType:
					var finalVal bigIntValue
					finalVal.value = new(big.Int)
					finalVal.value.SetInt64(1)
					for _, o := range operands {
						v, ok := o.Val.(bigIntValue)
						if ok {
							finalVal.value = finalVal.value.Mul(finalVal.value, v.value)
						} else {
							fmt.Errorf("Error while converting %s to bigIntValue\n", o.Val.Str())
						}
					}
					retVal.Val = finalVal
					break

				case floatType:
					var finalVal floatValue
					finalVal.value = 1
					for _, o := range operands {
						v, ok := o.Val.(floatValue)
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

			performOp:
				switch finalType {
				case intType:
					var finalVal intValue
					var val1, val2 intValue
					finalVal.value = 1

					var ok bool
					val1, ok = operands[0].Val.(intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(intValue)
					if !ok {
						fmt.Errorf("Error while converting %s to intValue\n", operands[1].Val.Str())
					}
					if val2.value != 0 {
						// Check for overflow/underflow here.
						if val1.value == math.MinInt64 && val2.value == -1 {
							err := tryTypeCastTo(&operands, bigIntType)
							if err != nil {
								fmt.Printf("Problem while avoiding overflow in operand %s: %s.\n", add, err)
							} else {
								finalType = bigIntType
								goto performOp
							}
						}

						finalVal.value = val1.value / val2.value
						retVal.Val = finalVal
					} else {
						retVal.Err = errors.New(fmt.Sprintf("divide by zero"))
					}
					break

				case bigIntType:
					var finalVal bigIntValue
					var val1, val2 bigIntValue
					finalVal.value = new(big.Int)
					finalVal.value.SetInt64(1)

					var ok bool
					val1, ok = operands[0].Val.(bigIntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to bigIntValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(bigIntValue)
					if !ok {
						fmt.Errorf("Error while converting %s to bigIntValue\n", operands[1].Val.Str())
					}
					if val2.value.Cmp(new(big.Int).SetInt64(0)) != 0 {
						finalVal.value.Div(val1.value, val2.value)
						retVal.Val = finalVal
					} else {
						retVal.Err = errors.New(fmt.Sprintf("divide by zero"))
					}
					break

				case floatType:
					var finalVal floatValue
					var val1, val2 floatValue
					finalVal.value = 0

					var ok bool
					val1, ok = operands[0].Val.(floatValue)
					if !ok {
						fmt.Printf("Error while converting %s to floatValue\n", operands[0].Val.Str())
					}

					val2, ok = operands[1].Val.(floatValue)
					if !ok {
						fmt.Printf("Error while converting %s to floatValue\n", operands[1].Val.Str())
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
			symbol:           def,
			minArgCount:      2,
			maxArgCount:      2,
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
					var val1, val2 intValue
					val1, _ = operands[0].Val.(intValue)
					val2, _ = operands[1].Val.(intValue)
					retVal.Val = newBoolValue(val1.value > val2.value)
					break

				case floatType:
					var val1, val2 floatValue
					val1, _ = operands[0].Val.(floatValue)
					val2, _ = operands[1].Val.(floatValue)
					retVal.Val = newBoolValue(val1.value > val2.value)
					break

				case stringType:
					var val1, val2 stringValue
					val1, _ = operands[0].Val.(stringValue)
					val2, _ = operands[1].Val.(stringValue)
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
					var val1, val2 intValue
					val1, _ = operands[0].Val.(intValue)
					val2, _ = operands[1].Val.(intValue)
					retVal.Val = newBoolValue(val1.value >= val2.value)
					break

				case floatType:
					var val1, val2 floatValue
					val1, _ = operands[0].Val.(floatValue)
					val2, _ = operands[1].Val.(floatValue)
					retVal.Val = newBoolValue(val1.value >= val2.value)
					break

				case stringType:
					var val1, val2 stringValue
					val1, _ = operands[0].Val.(stringValue)
					val2, _ = operands[1].Val.(stringValue)
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
					var val1, val2 intValue
					val1, _ = operands[0].Val.(intValue)
					val2, _ = operands[1].Val.(intValue)
					retVal.Val = newBoolValue(val1.value < val2.value)
					break

				case floatType:
					var val1, val2 floatValue
					val1, _ = operands[0].Val.(floatValue)
					val2, _ = operands[1].Val.(floatValue)
					retVal.Val = newBoolValue(val1.value < val2.value)
					break

				case stringType:
					var val1, val2 stringValue
					val1, _ = operands[0].Val.(stringValue)
					val2, _ = operands[1].Val.(stringValue)
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
					var val1, val2 intValue
					val1, _ = operands[0].Val.(intValue)
					val2, _ = operands[1].Val.(intValue)
					retVal.Val = newBoolValue(val1.value <= val2.value)
					break

				case floatType:
					var val1, val2 floatValue
					val1, _ = operands[0].Val.(floatValue)
					val2, _ = operands[1].Val.(floatValue)
					retVal.Val = newBoolValue(val1.value <= val2.value)
					break

				case stringType:
					var val1, val2 stringValue
					val1, _ = operands[0].Val.(stringValue)
					val2, _ = operands[1].Val.(stringValue)
					retVal.Val = newBoolValue(val1.value <= val2.value)
					break
				}
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      defun,
			minArgCount: 3,
			maxArgCount: 3,
			passRawAST:  true,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				astVal, ok := operands[0].Val.(astValue)
				if !ok {
					retVal.Err = errors.New(fmt.Sprintf("operand[0] has to be astValue"))
					return retVal
				}
				// Check astNode[0] is of varType, and is not registered in varMap
				if !astVal.astNodes[0].isValue {
					retVal.Err = errors.New(fmt.Sprintf("Method name not defined correctly."))
					return retVal
				}
				methodNameVal, err := getValue(env, astVal.astNodes[0].value)
				if err != nil || methodNameVal.getValueType() != varType {
					retVal.Err = errors.New(fmt.Sprintf("Expecting method name, got %s", astVal.astNodes[0].value))
					return retVal
				}

				methodName := methodNameVal.Str()
				if _, ok := env.varMap[methodNameVal.Str()]; ok {
					retVal.Err = errors.New(fmt.Sprintf("Method %s already defined as a variable", methodName))
					return retVal
				}

				if _, ok := env.opMap[methodNameVal.Str()]; ok {
					retVal.Err = errors.New(fmt.Sprintf("Method %s already defined as an operator", methodName))
					return retVal
				}

				if astVal.astNodes[1].isValue {
					retVal.Err = errors.New(fmt.Sprintf("Missing list of parameters for method %s", methodName))
					return retVal
				}

				// Check astNode[1] is a list of varTypes.
				params := make([]string, 0)
				for i, node := range astVal.astNodes[1].children {
					if !node.isValue {
						retVal.Err = errors.New(fmt.Sprintf("Malformed parameter %d in method %s.", i, methodName))
						return retVal
					}
					paramName := node.value
					val, err := getValue(env, paramName)
					if err != nil || val.getValueType() != varType {
						retVal.Err = errors.New(fmt.Sprintf("Malformed parameter %s in method %s.", paramName, methodName))
						return retVal
					}
					params = append(params, paramName)
				}

				addOperator(opMap,
					&Operator{
						symbol:      methodName,
						minArgCount: len(params),
						maxArgCount: len(params),
						handler: func(env *LangEnv, operands []Atom) Atom {
							var retVal Atom
							maxRecursionLimit := 100000
							newEnv := NewEnv()

							// Copy all the operators of the parent env.
							newEnv.opMap = make(map[string]*Operator, 0)
							for k, v := range env.opMap {
								newEnv.opMap[k] = v
							}

							// Copy all the variable values of the parent env.
							// We will favor formal arguments over previously defined variables.
							newEnv.varMap = make(map[string]Value, 0)
							for k, v := range env.varMap {
								newEnv.varMap[k] = v
							}

							// fmt.Printf("Executing the method %s with values: \n", methodName)
							for i, p := range params {
								// Check here whether operands[i] is a variable / operator.
								if op, ok := opMap[operands[i].Val.Str()]; ok {
									newEnv.opMap[p] = op
									// fmt.Printf("%s = %s (op)\n", p, operands[i].Val.Str())
								} else {
									newEnv.varMap[p] = operands[i].Val
									// fmt.Printf("%s = %s (op)\n", p, operands[i].Val.Str())
								}
							}

							// fmt.Printf(", d: %d\n", env.recursionDepth)
							newEnv.recursionDepth = env.recursionDepth + 1
							if newEnv.recursionDepth > maxRecursionLimit {
								retVal.Err = errors.New(fmt.Sprintf("Reached the recursion limit of %d. Terminating.", maxRecursionLimit))
								return retVal
							}

							// fmt.Printf("AST structure: %s\n", getASTStr(astVal.astNodes[2]))
							retVal = evalASTHelper(newEnv, astVal.astNodes[2])
							// fmt.Printf("evalASTHelper returned with %s\n", retVal.Val.Str())
							return retVal
						},
					},
				)
				var val varValue
				val.value = fmt.Sprintf("<Method: %s>", methodName)
				val.varName = methodName
				retVal.Val = val
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      and,
			minArgCount: 2,
			maxArgCount: 100,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				_, retVal.Err = typeCoerce(and, &operands, boolValPrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				result := true
				for _, o := range operands {
					v, ok := o.Val.(boolValue)
					if ok {
						result = result && v.value
					}
				}
				retVal.Val = newBoolValue(result)
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      or,
			minArgCount: 2,
			maxArgCount: 100,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				_, retVal.Err = typeCoerce(or, &operands, boolValPrecedenceMap)
				if retVal.Err != nil {
					return retVal
				}

				result := false
				for _, o := range operands {
					v, ok := o.Val.(boolValue)
					if ok {
						result = result || v.value
					}
				}
				retVal.Val = newBoolValue(result)
				return retVal
			},
		},
	)

	addOperator(opMap,
		&Operator{
			symbol:      cond,
			minArgCount: 1,
			maxArgCount: 100,
			passRawAST:  true,
			handler: func(env *LangEnv, operands []Atom) Atom {
				var retVal Atom
				astNodeVal, _ := operands[0].Val.(astValue)
				for i, astNode := range astNodeVal.astNodes {
					if len(astNode.children) != 2 {
						retVal.Err = errors.New(fmt.Sprintf(
							"Arguments for %s should be of the format `(condition value)`.",
							cond))
						return retVal
					}
					condValue := evalAST(env, astNode.children[0])
					if condValue.Err != nil {
						return condValue
					}
					if condValue.Val.getValueType() != boolType {
						retVal.Err = errors.New(fmt.Sprintf(
							"Arguments for operand %d for %s was of type %s instead of %s.",
							i+1, cond, condValue.Val.getValueType(), boolType))
					}
					condBoolValue, _ := condValue.Val.(boolValue)
					if condBoolValue.value {
						return evalAST(env, astNode.children[1])
					}

					retVal.Err = errors.New(fmt.Sprintf(
						"None of the arguments for %s evaluated to true.",
						cond))
				}
				return retVal
			},
		},
	)
}
