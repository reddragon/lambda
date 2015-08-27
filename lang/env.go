package lang

// Data required for interpretation of the language.
// We start with the default environment, and build on top of it, over time.
type LangEnv struct {
	opMap  map[string]*Operator
	types  []Value
	varMap map[string]Value
}

// Initialize the environment
func (e *LangEnv) Init() {
	// e.opMap = make(map[string]*Operator)
	e.opMap = builtinOperators()
	e.types = builtinTypes()
	e.varMap = make(map[string]Value)
}

func (e *LangEnv) getOperator(sym string) *Operator {
	return e.opMap[sym]
}

func (e *LangEnv) getValue(sym string) Value {
	return e.varMap[sym]
}
