package lang

import (
  "fmt"
  "testing"
)

type TestPair struct {
  a, b interface {}
}

func doTypeChecks(v Value, cases []TestPair, t *testing.T) {
  valueType := v.valueType()
  for _, p := range cases {
    check := v.ofType(p.a)
    if check != p.b {
      t.Errorf(fmt.Sprintf("%s.ofType(%s), expected: %t, actual: %t",
        valueType, p.a, p.b, check))
    }
  }
}

func TestIntValue(t *testing.T) {
  iv := new(IntValue)
  cases := make([]TestPair, 0)
  cases = append(cases, TestPair{"1.2", false})
  cases = append(cases, TestPair{"1", true})
  cases = append(cases, TestPair{"-1", true})
  cases = append(cases, TestPair{"foobar", false})

  doTypeChecks(iv, cases, t)
  // if iv.to(StringType)
}

func TestStringValue(t *testing.T) {
  sv := new(StringValue)
  cases := make([]TestPair, 0)
  cases = append(cases, TestPair{"", false})
  cases = append(cases, TestPair{"''", true})
  cases = append(cases, TestPair{"'abc'", true})
  cases = append(cases, TestPair{"\"\"", true})
  cases = append(cases, TestPair{"\"abc\"", true})
  cases = append(cases, TestPair{"1.2", false})

  doTypeChecks(sv, cases, t)
}
