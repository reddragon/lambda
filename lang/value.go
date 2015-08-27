package lang

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
)

// Different types of values supported
type valueType interface{}

const (
	// Value type
	stringType = "stringType"
	intType    = "intType"
	floatType  = "floatType"
	varType    = "varType"
)

type Value interface {
	Str() string
	getValueType() valueType
	to(valueType) (Value, error)
	ofType(string) bool
	newValue(string) Value
}

// Algorithm
// 1. Go through all the value types, in order.
// 2. Pick the highest value type that complies.
// 3. Return that value type.
func getValue(env *LangEnv, token string) (Value, error) {
	// TODO Use env types
	types := builtinTypes()
	for _, t := range types {
		if t.ofType(token) {
			if t.getValueType() == varType {
				val := env.varMap[token]
				if val != nil {
					return val, nil
				}
			}
			return t.newValue(token), nil
		}
	}

	return nil, errors.New(fmt.Sprintf("Could not get type for token: %s", token))
}

/*
Types in eclisp:
> 1 + 1
2
> 1 + 1.0
2.0
> 1 * 1.0
2.0
> 1 / 0
NaN
>
*/

func typeConvError(from, to valueType) error {
	return errors.New(fmt.Sprintf("Cannot convert %s to %s", from, to))
}

type stringValue struct {
	value string
}

func (v stringValue) getValueType() valueType {
	return stringType
}

func (v stringValue) to(targetType valueType) (Value, error) {
	switch targetType {
	case stringType:
		return v, nil
	}
	return nil, typeConvError(v.getValueType(), targetType)
}

func (v stringValue) Str() string {
	return v.value
}

func (v stringValue) ofType(targetValue string) bool {
	valLen := len(targetValue)
	if valLen < 2 {
		return false
	}
	// TODO
	// Stricter checks for quotes inside strings, like ''' should not be valid.
	f, l := targetValue[0], targetValue[valLen-1]
	if (f == '\'' && l == '\'') || (f == '"' && l == '"') {
		return true
	}
	return false
}

func (v stringValue) newValue(str string) Value {
	val := new(stringValue)
	val.value = str
	return val
}

type intValue struct {
	value int64
}

func (v intValue) getValueType() valueType {
	return intType
}

func (v intValue) to(targetType valueType) (Value, error) {
	switch targetType {
	case intType:
		return v, nil
	case floatType:
		val := new(floatValue)
		val.value = float64(v.value)
		return val, nil
	}
	return nil, typeConvError(v.getValueType(), targetType)
}

func (v intValue) ofType(targetValue string) bool {
	_, err := strconv.ParseInt(targetValue, 0, 64)
	if err != nil {
		// fmt.Printf("Error processing %s: %s", targetValue, err)
		return false
	}
	return true
}

func (v intValue) Str() string {
	return strconv.FormatInt(v.value, 10)
}

func (v intValue) newValue(str string) Value {
	intVal, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		return nil
	}
	val := new(intValue)
	val.value = intVal
	return val
}

type floatValue struct {
	value float64
}

func (v floatValue) getValueType() valueType {
	return floatType
}

func (v floatValue) to(targetType valueType) (Value, error) {
	switch targetType {
	case floatType:
		return v, nil
	}
	return nil, typeConvError(v.getValueType(), targetType)
}

func (v floatValue) ofType(targetValue string) bool {
	_, err := strconv.ParseFloat(targetValue, 64)
	if err != nil {
		// fmt.Printf("Error processing %s: %s", targetValue, err)
		return false
	}
	return true
}

func (v floatValue) Str() string {
	return strconv.FormatFloat(v.value, 'g', -1, 64)
}

func (v floatValue) newValue(str string) Value {
	floatVal, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil
	}
	val := new(floatValue)
	val.value = floatVal
	return val
}

type varValue struct {
	value string
}

func (v varValue) getValueType() valueType {
	return varType
}

func (v varValue) to(targetType valueType) (Value, error) {
	return nil, typeConvError(v.getValueType(), targetType)
}

func (v varValue) ofType(targetValue string) bool {
	matched, err := regexp.MatchString("[a-zA-Z]+[a-zA-Z0-9]*", targetValue)
	if matched == false || err != nil {
		return false
	}
	return true
}

func (v varValue) Str() string {
	return v.value
}

func (v varValue) newValue(str string) Value {
	val := new(varValue)
	val.value = str
	return val
}
