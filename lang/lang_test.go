package lang

import (
	"testing"
)

func saneQueryTest(query, expected string, t *testing.T, env *LangEnv) {
	val := Eval(query, env)
	if val.Val == nil || val.Err != nil {
		t.Errorf("Expected %s to not be nil", query)
	} else {
		if val.Val.Str() != expected {
			t.Errorf("Expected %s to be %s, but was %s", query, expected, val.Val.Str())
		}
	}
}

func malformedQueryTest(query string, t *testing.T, env *LangEnv) {
	val := Eval(query, env)
	if val.Val != nil {
		t.Errorf("Expected value of %s to be nil, was %s", val.Val.Str())
	} else if val.Err == nil {
		t.Errorf("Expected %s to return a non-nil error")
	}
}

func TestBasicLang(t *testing.T) {
	env := new(LangEnv)
	env.Init()

	saneQueryTest("(+ 1 2)", "3", t, env)
	saneQueryTest("(- 1 2)", "-1", t, env)
	saneQueryTest("(* 1 2)", "2", t, env)
	saneQueryTest("(/ 1 2)", "0", t, env)

	saneQueryTest("(+ 1.1 2.1)", "3.2", t, env)
	saneQueryTest("(- 1.3 2.1)", "-0.8", t, env)
	saneQueryTest("(* 1.3 2)", "2.6", t, env)
	saneQueryTest("(/ 1.3 2)", "0.65", t, env)

	malformedQueryTest(")(", t, env)
	malformedQueryTest(")", t, env)
	malformedQueryTest("(", t, env)
	malformedQueryTest("]]]", t, env)

	malformedQueryTest("(/ 1 0)", t, env)
}
