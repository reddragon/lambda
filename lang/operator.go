package lang

type Operator struct {
	symbol   string
	argCount int
	handler  (func([]Atom) Atom)
}

const (
	// Operators
	add string = "+"
	sub string = "-"
	mul string = "*"
	div string = "/"
)
