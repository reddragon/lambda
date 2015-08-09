package lang

// Different types of values supported
type ValueType int

const (
	// Value type
  StringType = iota
	IntType
	FloatType
)

type Value interface {
	valueType() ValueType
	to(ValueType) (Value, error)
	str() string
  ofType(string) bool
	// TODO
	// Add other methods
}

func GetType(token string) ValueType {
  // Algorithm
  // 1. Go through all the value types, in order.
  // 2. Pick the highest value type that complies.
  // 3. Return that value type.
}

/*
Types in eclisp:
> 1 + 1
2
> 1 + 1.0
2.0
> 1 * 1.0
2.0
> 1 / 0
NaN
> 
*/

type StringValue struct {
  value string
}

func (v StringValue) valueType() ValueType {
  return StringType
}

func (v StringValue) to(targetType ValueType) (Value, error) {
  // TODO
  switch targetType {
  case StringType: return value
  // TODO Implement others
  }
}

func (v StringValue) str() string {
  return v.to(StringType)
}

func (v StringType) ofType(targetValue string) bool {
  return true
}

/*
type FloatValue struct {
  value float64
}

type IntValue struct {
  value int
}

func (v IntValue) valueType() ValueType {
  return IntType
}

func (v IntValue) to(vtype ValueType) (Value, error) {
  switch vtype {
    case IntType:
    case FloatType:
    case StringType:
  }
  return nil, nil

  var r FloatValue
  return r, nil
}
*/
