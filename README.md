# eclisp
<img src="https://travis-ci.org/reddragon/eclisp.svg?branch=master"/>

This is a WIP!

A Lisp dialect written in Golang, which will hopefully 'eclipse' other Lisp dialects :) 

I have been amazed at the kind of things that we can achieve with simple s-expressions. This is my attempt at writing yet
another Lisp dialect. Right now what we have is a simple REPL, which can work on simple operators like +, -, * and /. It ~~can't even~~ can do floating-point math.

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

eclisp> (/ 1 0)
Error: Divide by zero

eclisp> ^D
Goodbye!
```

### Inspiration
* [Peter Norvig's post about writing a Lisp-like language](http://norvig.com/lispy.html)
* [Build Your Own Lisp](http://www.buildyourownlisp.com/)
