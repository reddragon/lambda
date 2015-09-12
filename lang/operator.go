package lang

type Operator struct {
	symbol   string
	argCount int
	handler  (func(*LangEnv, []Atom) Atom)
}

const (
	// Operators
	add string = "+"
	sub string = "-"
	mul string = "*"
	div string = "/"
	def string = "defvar"
	eq  string = "eq"
)
