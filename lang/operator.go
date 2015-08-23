package lang

type Operator struct {
	symbol   string
	argCount int
	// TODO The return type has to change to Value
	handler (func([]Atom) Atom)
}

const (
	// Operators
	add string = "+"
	sub string = "-"
	mul string = "*"
	div string = "/"
)
