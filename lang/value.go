package lang

import (
	"errors"
	"fmt"
	"strconv"
)

// Different types of values supported
type ValueType string

const (
	// Value type
	StringType = "StringType"
	IntType    = "IntType"
	FloatType  = "FloatType"
)

type Value interface {
	valueType() ValueType
	to(ValueType) (Value, error)
	str() string
	ofType(string) bool
	newValue(string) Value
	// TODO
	// Add other methods
}

// Algorithm
// 1. Go through all the value types, in order.
// 2. Pick the highest value type that complies.
// 3. Return that value type.
func GetValue(token string) (Value, error) {
	// TODO Use env types
	types := builtinTypes()
	for _, t := range types {
		if t.ofType(token) {
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

func typeConvError(from, to ValueType) error {
	return errors.New(fmt.Sprintf("Cannot convert %s to %s", from, to))
}

type StringValue struct {
	value string
}

func (v StringValue) valueType() ValueType {
	return StringType
}

func (v StringValue) to(targetType ValueType) (Value, error) {
	switch targetType {
	case StringType:
		return v, nil
	}
	return nil, typeConvError(v.valueType(), targetType)
}

func (v StringValue) str() string {
	return v.value
}

func (v StringValue) ofType(targetValue string) bool {
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

func (v StringValue) newValue(str string) Value {
	val := new(StringValue)
	val.value = str
	return val
}

type IntValue struct {
	value int64
}

func (v IntValue) valueType() ValueType {
	return IntType
}

func (v IntValue) to(targetType ValueType) (Value, error) {
	switch targetType {
	case IntType:
		return v, nil
	case FloatType:
		return FloatValue{float64(v.value)}, nil
	}
	return nil, typeConvError(v.valueType(), targetType)
}

func (v IntValue) ofType(targetValue string) bool {
	_, err := strconv.ParseInt(targetValue, 0, 64)
	if err != nil {
		// fmt.Printf("Error processing %s: %s", targetValue, err)
		return false
	}
	return true
}

func (v IntValue) str() string {
	return strconv.FormatInt(v.value, 10)
}

func (v IntValue) newValue(str string) Value {
	intVal, err := strconv.ParseInt(str, 0, 64)
	if err != nil {
		return nil
	}
	val := new(IntValue)
	val.value = intVal
	return val
}

type FloatValue struct {
	value float64
}

func (v FloatValue) valueType() ValueType {
	return FloatType
}

func (v FloatValue) to(targetType ValueType) (Value, error) {
	switch targetType {
	case FloatType:
		return v, nil
	}
	return nil, typeConvError(v.valueType(), targetType)
}

func (v FloatValue) ofType(targetValue string) bool {
	_, err := strconv.ParseFloat(targetValue, 64)
	if err != nil {
		// fmt.Printf("Error processing %s: %s", targetValue, err)
		return false
	}
	return true
}

func (v FloatValue) str() string {
	return strconv.FormatFloat(v.value, 'g', -1, 64)
}

func (v FloatValue) newValue(str string) Value {
	floatVal, err := strconv.ParseFloat(str, 64)
	if err != nil {
		return nil
	}
	val := new(FloatValue)
	val.value = floatVal
	return val
}
