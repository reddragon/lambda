package lang

func addOperator(opMap map[string]*Operator, op *Operator) {
	opMap[op.symbol] = op
}

// This methods returns all the builtin operators to the environment
func builtinOperators() map[string]*Operator {
  // builtins := make(*Operator[], 0)
	opMap := make(map[string]*Operator)
  // This method adds the default
	addOperator(opMap,
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
	)
  return opMap
}

func builtinTypes() ([]Value) {
	types := make([]Value, 0)
	types = append(types, new(StringValue))
	types = append(types, new(IntValue))
	types = append(types, new(FloatValue))
	return types
}
