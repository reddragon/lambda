package main

import (
	"fmt"
	l "github.com/reddragon/eclisp/lang"
	"github.com/tiborvass/uniline"
)

func process(env *l.LangEnv, line string) {
	retVal := l.Eval(line, env)
	if retVal.Err != nil {
		fmt.Printf("Error: %s\n", retVal.Err)
	} else {
		if retVal.Val != nil {
			fmt.Printf("%s\n", retVal.Val.Str())
		} else {
			fmt.Printf("\n")
		}
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
			// printType(line)
			process(env, line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	} else {
		fmt.Println("Goodbye!")
	}
}
