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
	boolType   = "boolType"
	astType    = "astType"
)

type Value interface {
	Str() string
	getValueType() valueType
	to(valueType) (Value, error)
	ofType(string) bool
	newValue(string) Value
}

func getVarValue(env *LangEnv, varValue Value) (Value, error) {
	if varValue != nil && varValue.getValueType() == varType {
		val := env.varMap[varValue.Str()]
		if val != nil {
			return val, nil
		}
		return nil, errors.New(fmt.Sprintf("Undefined variable: %s", varValue.Str()))
	}
	return nil, errors.New(fmt.Sprintf("Error while resolving variable."))
}

// Algorithm
// 1. Go through all the value types, in order.
// 2. Pick the highest value type that complies.
// 3. Return that value type.
func getValue(env *LangEnv, token string) (Value, error) {
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
	var val stringValue
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
		var val floatValue
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
	var val intValue
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
	var val floatValue
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
	var val varValue
	val.value = str
	return val
}

type boolValue struct {
	value bool
}

func (v boolValue) getValueType() valueType {
	return boolType
}

func (v boolValue) to(targetType valueType) (Value, error) {
	return nil, typeConvError(v.getValueType(), targetType)
}

func (v boolValue) ofType(targetValue string) bool {
	if targetValue == "true" || targetValue == "false" {
		return true
	}
	return false
}

func (v boolValue) Str() string {
	if v.value == true {
		return "true"
	} else {
		return "false"
	}
}

func (v boolValue) newValue(str string) Value {
	var val boolValue
	if str == "true" {
		val.value = true
	} else {
		val.value = false
	}
	return val
}

func newBoolValue(b bool) Value {
	var val boolValue
	val.value = b
	return val
}

type astValue struct {
	astNodes []*ASTNode
}

func (v astValue) getValueType() valueType {
	return astType
}

func (v astValue) to(targetType valueType) (Value, error) {
	return nil, typeConvError(v.getValueType(), targetType)
}

func (v astValue) ofType(targetValue string) bool {
	return false
}

func (v astValue) Str() string {
	return ""
}

func (v astValue) newValue(str string) Value {
	return nil
}

func newASTValue(astNodes []*ASTNode) Value {
	var val astValue
	val.astNodes = astNodes
	return val
}

type method struct {
	methodName string
	params     []string
	ast        *ASTNode
}
