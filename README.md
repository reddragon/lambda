# lambda
<img src="https://travis-ci.org/reddragon/lambda.svg?branch=master"/>

A WIP Lisp dialect written in Golang, written purely for fun :-) 

### Introduction

I have been amazed at the kind of things that we can achieve with simple s-expressions. s-expressions or symbolic exressions, are nothing but expressions of this format: `(operator operand1 operand2 ...)`, where we pass the operands to the operator for evaluation. For example, `(+ 1 2 3)` is the same as `1 + 2 + 3`. You can nest several such expressions like this. For example, `(* (+ 1 2) (+ 4 5))`, which is the same as `((1 + 2) * (4 + 5))`. In general, s-expressions make it easy for tree-structured code and data easy.

Lisp, is a family of programming languages that have popularized the use of s-expressions. I found it interesting, and this is my attempt at writing yet another Lisp dialect. This is a purely academic pursuit, so don't use this in production. The crux of the work lies in the `lang` directory, but I have provided a simple REPL (Read-Eval-Print-Loop) to try out the language.

#### What works so far
* Integer, floating point and string types
* Mathematical operators (`+`, `-`, `*`, `/`)
* Comparison operators (`=`, `>`, `>=`, `<`, `<=`)
* Logical operators (`or`, `and`)
* Conditional (`cond`)
* Defining variables (`defvar`)
* Defining methods (`defun`) (Can't define multi-expressions methods yet)

#### What's coming
* Support for Lambdas
* Multi-expression methods
* Support for comments

You can check the 'Sample Usage' section for a walk-through. I'll add more thorough documentation in the near future.

### How to Use
* `go get github.com/reddragon/lambda`
* `go build $GOPATH/src/github.com/reddragon/lambda/lambda.go`
* `$GOPATH/bin/lambda`

### Sample Usage
```
> ./lambda
lambda> (+ 1 2)
3

lambda> (- (/ (* (+ 1 2) 3) 3) 2)
1

lambda> (/ 22 7.0)
3.142857142857143

lambda> (defvar pi 3.14159265359)
3.14159265359

lambda> (defvar r 10)
10

lambda> (* pi (* r r))
314.159265359

lambda> (/ 1 0)
Error: Divide by zero

lambda> (defun addSq(x y) (+ (* x x) (* y y)))
<Method: addSq>

lambda> (addSq 3 4)
25

lambda> (defun fact (x) (cond (= x 0) 1 (* x (fact (- x 1)))))
<Method: fact>

lambda> (fact 10)
3628800

lambda> ^D
Goodbye!
```
lambda can also read from files and execute them. You can try it out with the `-f` option, like:

```
./lambda -f ~/path/to/my/script.l
```

### Inspiration
* [Peter Norvig's post about writing a Lisp-like language](http://norvig.com/lispy.html)
* [Build Your Own Lisp](http://www.buildyourownlisp.com/)
