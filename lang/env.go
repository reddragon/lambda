package lang

// Data required for interpretation of the language.
// We start with the default environment, and build on top of it, over time.
type LangEnv struct {
	opMap map[string]*Operator
}

// Initialize the environment
func (e *LangEnv) Init() {
	e.opMap = make(map[string]*Operator)
	e.builtin()
}

func (e *LangEnv) getOperator(sym string) *Operator {
	return e.opMap[sym]
}

func (e *LangEnv) addOperator(op *Operator) {
	e.opMap[op.symbol] = op
}

// This methods adds all the builtin operators to the environment
func (e *LangEnv) builtin() {
	// This method adds the default
	e.addOperator(
		&Operator{
			symbol:   Add,
			argCount: 2,
			handler: func(operands []Atom) Atom {
				var retVal Atom
				argCheckErr := checkArgLen(Add, operands, 2)
				if argCheckErr.Err != nil {
					return argCheckErr
				}
				// TODO
				// Do type checks

				retVal.Val = operands[0].Val + operands[1].Val
				return retVal
			},
		},
	)

	e.addOperator(
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

	e.addOperator(
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

	e.addOperator(
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
	)
}
