# lambda
<img src="https://travis-ci.org/reddragon/lambda.svg?branch=master"/>

A ~~WIP~~ Lisp dialect written in Golang, written purely for fun :-)

### Introduction

I have been amazed at the kind of things that we can achieve with simple s-expressions. s-expressions or symbolic exressions, are nothing but expressions of this format: `(operator operand1 operand2 ...)`, where we pass the operands to the operator for evaluation. For example, `(+ 1 2 3)` is the same as `1 + 2 + 3`. You can nest several such expressions like this. For example, `(* (+ 1 2) (+ 4 5))`, which is the same as `((1 + 2) * (4 + 5))`. In general, s-expressions make it easy to write tree-structured / recursive code and data.

Lisp, is a family of programming languages that have popularized the use of s-expressions. I found it interesting, and this is my attempt at writing yet another Lisp dialect. This is a purely academic pursuit, so I would not recommend using this in production. The crux of the work lies in the `lang` directory, but I have provided a simple REPL (Read-Eval-Print-Loop) to try out the language.

#### What works so far
* Integer, floating point and string types
* Mathematical operators (`+`, `-`, `*`, `/`)
* Comparison operators (`=`, `>`, `>=`, `<`, `<=`)
* Logical operators (`or`, `and`)
* Conditional (`cond`)
* Defining variables (`defvar`)
* Defining methods (`defun`) (Can't define multi-expressions methods yet)
* Methods as first-class citizens
* Support for Big Int calculations

#### What might come*
* Full Support for anonymous methods
* Multi-expression methods
* Support for comments

**Update**: I am going to to shift my attention to other projects as of July 2016. If you feel strongly about a particular feature, either feel free to implement it and send a pull request (I can help with giving pointers), or let me know and I will try to prioritize it.

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

lambda> (defun add-sq (x y) (+ (* x x) (* y y)))
<Method: add-sq>

lambda> (add-sq 3 4)
25

lambda> (defun fact (x) (cond ((= x 0) 1) (true (* x (fact (- x 1))))))
<Method: fact>

lambda> (fact 30)
265252859812191058636308480000000

lambda> (defun add-one (x) (+ x 1))
<Method: add-one>

lambda> (defun twice (f x) (f (f x)))
<Method: twice>

lambda> (twice add-one 2)
4

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
