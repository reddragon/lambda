# eclisp
<img src="https://travis-ci.org/reddragon/eclisp.svg?branch=master"/>

This is a WIP!

A Lisp dialect written in Golang, which will hopefully 'eclipse' other Lisp dialects :) 

I have been amazed at the kind of things that we can achieve with simple s-expressions. This is my attempt at writing yet
another Lisp dialect. Right now what we have is a simple REPL, which can work on simple operators like +, -, * and /. It can do floating-point math, variables, comparisons such as [`eq`, `>`, `>=`, `<`, `<=`]. I will be adding other goodness pretty soon.

### How to Use
* `go get github.com/reddragon/eclisp`
* `go build $GOPATH/src/github.com/reddragon/eclisp/eclisp.go`
* `$GOPATH/bin/eclisp`

### Sample Usage
```
> ./eclisp
eclisp> (+ 1 2)
3

eclisp> (- (/ (* (+ 1 2) 3) 3) 2)
1

eclisp> (/ 22 7.0)
3.142857142857143

eclisp> (defvar pi 3.14159265359)
3.14159265359

eclisp> (defvar r 10)
10

eclisp> (* pi (* r r))
314.159265359

eclisp> (/ 1 0)
Error: Divide by zero

eclisp> (defun addSq(x y) (+ (* x x) (* y y)))
<Method: addSq>

eclisp> (addSq 3 4)
25

eclisp> ^D
Goodbye!
```
eclisp can also read from files and execute them. Currently it only executes the first well-formed expression, but you can try it out with the `-f` option, like:

```
./eclisp -f ~/path/to/my/script.el
```

### Inspiration
* [Peter Norvig's post about writing a Lisp-like language](http://norvig.com/lispy.html)
* [Build Your Own Lisp](http://www.buildyourownlisp.com/)
