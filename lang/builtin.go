package lang

// This methods returns all the builtin operators to the environment
func builtinOperators() map[string]*Operator {
	opMap := make(map[string]*Operator)
	addBuiltinOperators(opMap)
	return opMap
}

func builtinTypes() []Value {
	types := make([]Value, 0)
	// Append the values in the order of prefence, i.e, more specific types
	// should be first.
	types = append(types, new(stringValue))
	types = append(types, new(intValue))
	types = append(types, new(bigIntValue))
	types = append(types, new(floatValue))
	types = append(types, new(boolValue))
	types = append(types, new(varValue))
	return types
}
