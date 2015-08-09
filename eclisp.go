package main

import (
	"fmt"
	l "github.com/reddragon/eclisp/lang"
	"github.com/tiborvass/uniline"
)

func process(env *l.LangEnv, line string) {
	node, _, err := l.GetAST(line)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		var retVal l.Atom
		retVal = l.EvalAST(env, node)
		if retVal.Err != nil {
			fmt.Printf("Error: %s\n", retVal.Err)
		} else {
			fmt.Printf("%d\n", retVal.Val)
		}
	}
}

func printType(line string) {
	inferredType, err := l.GetType(line)
	if err != nil {
		fmt.Printf("Error: %s\n", err)
	} else {
		fmt.Printf("Type: %s", inferredType)
	}
}

func main() {
	// Setup the language environment
	env := new(l.LangEnv)
	env.Init()

	prompt := "eclisp> "
	scanner := uniline.DefaultScanner()
	for scanner.Scan(prompt) {
		line := scanner.Text()
		if len(line) > 0 {
			scanner.AddToHistory(line)
			process(env, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	} else {
		fmt.Println("Goodbye!")
	}
}
