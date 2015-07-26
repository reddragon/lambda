package lang

// Data required for interpretation of the language.
// We start with the default environment, and build on top of it, over time.
type LangEnv struct {
	opMap map[string]*Operator
}

// Initialize the environment
func (e *LangEnv) Init() {
	// e.opMap = make(map[string]*Operator)
	e.opMap = builtinOperators()
}

func (e *LangEnv) getOperator(sym string) *Operator {
	return e.opMap[sym]
}
