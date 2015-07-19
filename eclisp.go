package main

import (
	// "bufio"
	"fmt"
	// "strconv"
	// "os"
	l "github.com/reddragon/eclisp/lang"
	"github.com/tiborvass/uniline"
)

func process(line string) {
	tokens := l.Tokenize(line)
	if len(tokens) == 0 {
		fmt.Println("Nothing to evaluate")
		return
	}
	node, tokens, err := l.GetAST(tokens)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		// fmt.Println("Worked fine!")
		// fmt.Printf("Parsed AST: %s\n", l.StringifyAST(node))
		var retVal l.Atom
		retVal = l.EvalAST(node)
		if retVal.Err != nil {
			fmt.Printf("Error: %s\n", retVal.Err)
		} else {
			fmt.Printf("%d\n", retVal.Value)
		}
	}
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
