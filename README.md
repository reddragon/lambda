# lambda
<img src="https://travis-ci.org/reddragon/lambda.svg?branch=master"/>

This is a WIP!

A Lisp dialect written in Golang, written purely for fun :-) 

I have been amazed at the kind of things that we can achieve with simple s-expressions. This is my attempt at writing yet
another Lisp dialect. Right now what we have is a simple REPL, which can work on simple operators like +, -, * and /. It can do floating-point math, variables, comparisons such as [`eq`, `>`, `>=`, `<`, `<=`]. I will be adding other goodness pretty soon.

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
