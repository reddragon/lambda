package lang

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func saneExprTest(query string, t *testing.T, env *LangEnv) {
	val := Eval(query, env)
	if len(val.ValStr) == 0 || len(val.ErrStr) > 0 {
		t.Errorf("Expected %s to not be nil. Err: %s", query, val.ErrStr)
	}
}

func checkExprResultTest(query, expected string, t *testing.T, env *LangEnv) {
	val := Eval(query, env)
	if len(val.ValStr) == 0 || len(val.ErrStr) > 0 {
		t.Errorf("Expected %s to not be nil. Err: %s", query, val.ErrStr)
	} else {
		if val.ValStr != expected {
			t.Errorf("Expected %s to be %s, but was %s", query, expected, val.ValStr)
		}
	}
}

func malformedExprTest(query string, t *testing.T, env *LangEnv) {
	val := Eval(query, env)
	if len(val.ValStr) != 0 {
		t.Errorf("Expected value of %s to be nil, was %s", val.ValStr)
	} else if len(val.ErrStr) == 0 {
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

	checkExprResultTest("(cond (true 1) (false 2))", "1", t, env)
	checkExprResultTest("(cond (false 1) (true 2))", "2", t, env)
	checkExprResultTest("(cond (false 1) (true 2) (false 3))", "2", t, env)
	checkExprResultTest("(cond ((> 2 3) 1) ((= 3 3) 2))", "2", t, env)

	malformedExprTest(")(", t, env)
	malformedExprTest(")", t, env)
	malformedExprTest("(", t, env)
	malformedExprTest("]]]", t, env)

	malformedExprTest("(/ 1 0)", t, env)

	malformedExprTest("(cond (1 2))", t, env)
	malformedExprTest("(cond (false 1) (false 2))", t, env)

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

	saneExprTest("(defun fact (x) (cond ((= x 0) 1) (true (* x (fact (- x 1))))))", t, env)
	checkExprResultTest("(fact 10)", "3628800", t, env)
	saneExprTest("(defun fib (x) (cond ((= x 0) 0) (true (+ x (fib (- x 1))))))", t, env)
	checkExprResultTest("(fib 10)", "55", t, env)

	saneExprTest("(defvar varDefinedOutside 10)", t, env)
	saneExprTest("(defun add (formalArg) (+ varDefinedOutside formalArg))", t, env)
	checkExprResultTest("(add 11)", "21", t, env)

	// Checking that we correctly override previously defined vars of same name
	saneExprTest("(defvar formalArg -10)", t, env)
	checkExprResultTest("(add 11)", "21", t, env)

	// Check the magic number method.
	saneExprTest("(defun magic (x) (cond ((<= x 0) 1) (true (+ (magic (- x 1)) (* 2 (magic (- x 3)))))))", t, env)
	checkExprResultTest("(magic -10)", "1", t, env)
	checkExprResultTest("(magic 0)", "1", t, env)
	checkExprResultTest("(magic 10)", "309", t, env)
}
