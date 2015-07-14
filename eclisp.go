package main

import (
	// "bufio"
	"fmt"
	"strconv"
	// "os"
	"github.com/tiborvass/uniline"
	l "github.com/reddragon/eclisp/lang"
)

func process(line string) {
	tokens := l.Tokenize(line)
	if len(tokens) == 0 {
		fmt.Println("Nothing to evaluate")
		return
	}
	// fmt.Printf("%q\n", tokens)
	var retVal l.Atom
	retVal, tokens = l.Eval(tokens)

	if len(tokens) != 0 {
		retVal.ErrMsg = strconv.Itoa(len(tokens)) + " extra token(s) in the string."
		retVal.Err = true
	}
	if retVal.Err {
		fmt.Printf("Error: %s\n", retVal.ErrMsg)
	} else {
		// fmt.Println("All worked fine!")
		fmt.Printf("%d\n", retVal.Val)
	}
	fmt.Printf("\n")
}

func main() {
	prompt := "eclisp> "
	scanner := uniline.DefaultScanner()
	for scanner.Scan(prompt) {
		line := scanner.Text()
		if len(line) > 0 {
			scanner.AddToHistory(line)
			process(line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	} else {
		fmt.Println("Goodbye!")
	}
}
