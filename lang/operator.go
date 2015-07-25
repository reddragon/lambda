package lang

type Operator struct {
	symbol   string
	argCount int
	// TODO The return type has to change to Value
	handler (func([]Atom) Atom)
}

const (
	// Operators
	Add string = "+"
	Sub string = "-"
	Mul string = "*"
	Div string = "/"
)
