package lang

import (
	"fmt"
	"testing"
)

type TestPair struct {
	a, b interface{}
}

func doTypeChecks(v Value, cases []TestPair, t *testing.T) {
	vtype := v.GetValueType()
	for _, p := range cases {
		str, ok := p.a.(string)
		if !ok {
			t.Errorf("Error setting up the test cases.")
			return
		}

		check := v.ofType(str)
		if check != p.b {
			t.Errorf(fmt.Sprintf("%s.ofType(%s), expected: %t, actual: %t",
				vtype, p.a, p.b, check))
		}
	}
}

func doChecks(v Value, cases []TestPair, t *testing.T) {
	vtype := v.GetValueType()
	for _, p := range cases {
		_, ok := p.a.(string)
		if !ok {
			t.Errorf("Error setting up the test cases.")
			return
		}

		if p.a != p.b {
			t.Errorf(fmt.Sprintf("%s.Str(), expected: %s, actual: %s",
				vtype, p.b, p.a))
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

	strCases := make([]TestPair, 0)
	iv.value = 1
	strCases = append(strCases, TestPair{iv.Str(), "1"})
	iv.value = 0
	strCases = append(strCases, TestPair{iv.Str(), "0"})

	doChecks(iv, strCases, t)

	iv.value = 2
	conv, err := iv.To(IntType)
	if err != nil || conv.GetValueType() != IntType || conv.Str() != iv.Str() {
		t.Errorf("Could not convert from IntType to IntType")
	}

	conv, err = iv.To(FloatType)
	if err != nil || conv.GetValueType() != FloatType || conv.Str() != iv.Str() {
		t.Errorf("Could not convert from IntType to FloatType")
	}
}

func TestFloatValue(t *testing.T) {
	fv := new(FloatValue)
	cases := make([]TestPair, 0)
	cases = append(cases, TestPair{"1.2", true})
	cases = append(cases, TestPair{"1", true})
	cases = append(cases, TestPair{"-1", true})
	cases = append(cases, TestPair{"foobar", false})
	doTypeChecks(fv, cases, t)

	strCases := make([]TestPair, 0)
	fv.value = 1.0
	strCases = append(strCases, TestPair{fv.Str(), "1"})
	fv.value = -1.0
	strCases = append(strCases, TestPair{fv.Str(), "-1"})
	fv.value = -1.1
	strCases = append(strCases, TestPair{fv.Str(), "-1.1"})
	fv.value = -1.23456789
	strCases = append(strCases, TestPair{fv.Str(), "-1.23456789"})
	doChecks(fv, strCases, t)
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

	strCases := make([]TestPair, 0)
	sv.value = "\"abc\""
	strCases = append(strCases, TestPair{sv.Str(), sv.value})

	doChecks(sv, strCases, t)
}

func TestTypeInfer(t *testing.T) {
	v, e := GetValue("1")
	if v == nil || v.GetValueType() != IntType || e != nil {
		t.Errorf("Could not correctly GetValue(1)")
	}

	v, e = GetValue("1.2")
	if v == nil || v.GetValueType() != FloatType || e != nil {
		t.Errorf("Could not correctly GetValue(1)")
	}

	v, e = GetValue("'xyz'")
	if v == nil || v.GetValueType() != StringType || e != nil {
		t.Errorf("Could not correctly GetValue(1)")
	}
}
