package lang

import (
	"fmt"
	"testing"
)

type TestPair struct {
	a, b interface{}
}

func doTypeChecks(v Value, cases []TestPair, t *testing.T) {
	vtype := v.getValueType()
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
	vtype := v.getValueType()
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

func TestintValue(t *testing.T) {
	iv := new(intValue)
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
	conv, err := iv.to(intType)
	if err != nil || conv.getValueType() != intType || conv.Str() != iv.Str() {
		t.Errorf("Could not convert from intType to intType")
	}

	conv, err = iv.to(floatType)
	if err != nil || conv.getValueType() != floatType || conv.Str() != iv.Str() {
		t.Errorf("Could not convert from intType to floatType")
	}
}

func TestfloatValue(t *testing.T) {
	fv := new(floatValue)
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

func TeststringValue(t *testing.T) {
	sv := new(stringValue)
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
	env := new(LangEnv)
	env.Init()
	v, e := getValue(env, "1")
	if v == nil || v.getValueType() != intType || e != nil {
		t.Errorf("Could not correctly getValue(1)")
	}

	v, e = getValue(env, "1.2")
	if v == nil || v.getValueType() != floatType || e != nil {
		t.Errorf("Could not correctly getValue(1)")
	}

	v, e = getValue(env, "'xyz'")
	if v == nil || v.getValueType() != stringType || e != nil {
		t.Errorf("Could not correctly getValue(1)")
	}

	v, e = getValue(env, "true")
	if v == nil || v.getValueType() != boolType || e != nil {
		t.Errorf("Could not correctly getValue(1)")
	}

	v, e = getValue(env, "false")
	if v == nil || v.getValueType() != boolType || e != nil {
		t.Errorf("Could not correctly getValue(1)")
	}
}
