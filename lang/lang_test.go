package lang

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func saneExprTest(query string, t *testing.T, env *LangEnv) {
	val, _ := Eval(query, env)
	if val.Val == nil || val.Err != nil {
		t.Errorf("Expected %s to not be nil. Err: %s", query, val.Err)
	}
}

func checkExprResultTest(query, expected string, t *testing.T, env *LangEnv) {
	val, _ := Eval(query, env)
	if val.Val == nil || val.Err != nil {
		t.Errorf("Expected %s to not be nil. Err: %s", query, val.Err)
	} else {
		if val.Val.Str() != expected {
			t.Errorf("Expected %s to be %s, but was %s", query, expected, val.Val.Str())
		}
	}
}

func malformedExprTest(query string, t *testing.T, env *LangEnv) {
	val, _ := Eval(query, env)
	if val.Val != nil {
		t.Errorf("Expected value of %s to be nil, was %s", val.Val.Str())
	} else if val.Err == nil {
		t.Errorf("Expected %s to return a non-nil error")
	}
}

// Generates an s-expression out of the basic {+, -, *, /} operators.
func genSimpleOperatorsExpr(terminationProb float32, r *rand.Rand, depth int) string {
	p := r.Float32()
	// fmt.Printf("Got %f, terminationProb was %f, depth: %d\n", p, terminationProb, depth)
	if p < terminationProb || depth > 1000 {
		return fmt.Sprintf("%f", r.Float32())
	}
	operands := []string{"+", "-", "*", "/"}
	return fmt.Sprintf("(%s %s %s)", operands[r.Intn(len(operands))], genSimpleOperatorsExpr(terminationProb, r, depth+1), genSimpleOperatorsExpr(terminationProb, r, depth+1))
}

// This will run a lot of random smoke tests with simple operators.
func runRandomSmokeTests(t *testing.T, env *LangEnv) {
	seed := time.Now().Unix()
	// fmt.Printf("Using the seed %d. Use it to reproduce test failures.\n", seed)
	r := rand.New(rand.NewSource(seed))
	for i := 0; i < 100; i++ {
		expr := genSimpleOperatorsExpr(0.5, r, 0)
		saneExprTest(expr, t, env)
	}
}

func TestBasicLang(t *testing.T) {
	env := new(LangEnv)
	env.Init()

	checkExprResultTest("(+ 1 2)", "3", t, env)
	checkExprResultTest("(+ 1 2 3 4 5)", "15", t, env)
	checkExprResultTest("(+ \"Hello\" \",\" \"World!\")", "\"Hello,World!\"", t, env)
	checkExprResultTest("(- 1 2)", "-1", t, env)
	checkExprResultTest("(* 1 2)", "2", t, env)
	checkExprResultTest("(* 1 2 3 4 5)", "120", t, env)
	checkExprResultTest("(/ 1 2)", "0", t, env)

	checkExprResultTest("(+ 1.1 2.1)", "3.2", t, env)
	checkExprResultTest("(- 1.3 2.1)", "-0.8", t, env)
	checkExprResultTest("(* 1.3 2)", "2.6", t, env)
	checkExprResultTest("(/ 1.3 2)", "0.65", t, env)

	checkExprResultTest("(- 1 (/ 0.5509 0.5698))", "0.033169533169533194", t, env)
	checkExprResultTest("(- 1 (/ 6 3))", "-1", t, env)

	checkExprResultTest("(defvar x 2.0)", "2", t, env)
	checkExprResultTest("(+ x 2.0)", "4", t, env)
	checkExprResultTest("(defvar y 1.9)", "1.9", t, env)
	checkExprResultTest("(* x y)", "3.8", t, env)
	checkExprResultTest("(defvar i 5.0)", "5", t, env)
	checkExprResultTest("(defvar i 6.0)", "6", t, env)
	malformedExprTest("(+ i j)", t, env)
	checkExprResultTest("(defvar j 3.1)", "3.1", t, env)
	checkExprResultTest("(+ i j)", "9.1", t, env)

	checkExprResultTest("(> 3 2)", "true", t, env)
	checkExprResultTest("(> 2 3)", "false", t, env)
	checkExprResultTest("(< 3 2)", "false", t, env)
	checkExprResultTest("(< 2 3)", "true", t, env)
	checkExprResultTest("(>= 3 2)", "true", t, env)
	checkExprResultTest("(>= 3 3)", "true", t, env)
	checkExprResultTest("(<= 3 3)", "true", t, env)
	checkExprResultTest("(<= 3 2)", "false", t, env)

	checkExprResultTest("(> \"a\" \"b\")", "false", t, env)
	checkExprResultTest("(> \"b\" \"a\")", "true", t, env)
	checkExprResultTest("(< \"a\" \"b\")", "true", t, env)
	checkExprResultTest("(< \"b\" \"a\")", "false", t, env)

	checkExprResultTest("(and true false)", "false", t, env)
	checkExprResultTest("(and true true)", "true", t, env)
	checkExprResultTest("(and false false)", "false", t, env)
	checkExprResultTest("(and true true true true)", "true", t, env)
	checkExprResultTest("(and true true true false)", "false", t, env)

	checkExprResultTest("(or true false)", "true", t, env)
	checkExprResultTest("(or true true)", "true", t, env)
	checkExprResultTest("(or false false)", "false", t, env)
	checkExprResultTest("(or true true true true)", "true", t, env)
	checkExprResultTest("(or false false false true)", "true", t, env)

	malformedExprTest(")(", t, env)
	malformedExprTest(")", t, env)
	malformedExprTest("(", t, env)
	malformedExprTest("]]]", t, env)

	malformedExprTest("(/ 1 0)", t, env)

	runRandomSmokeTests(t, env)
}

func TestMethods(t *testing.T) {
	env := new(LangEnv)
	env.Init()

	saneExprTest("(defun foo (x) (+ 1 x))", t, env)
	checkExprResultTest("(foo 4)", "5", t, env)
	// TODO: This should fail
	// See https://github.com/reddragon/lambda/issues/10
	// malformedExprTest("(foo)", t, env)
	malformedExprTest("(foo 4 5)", t, env)
}
