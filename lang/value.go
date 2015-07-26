package lang

// Different types of values supported
type ValueType int

const (
	// Value type
	IntType = iota
	FloatType
	StringType
)

type Value interface {
	valueType() ValueType
	to(ValueType) (Value, error)
	str() string
	// TODO
	// Add other methods
}
