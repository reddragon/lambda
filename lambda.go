package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	l "github.com/reddragon/lambda/lang"
	"github.com/tiborvass/uniline"
	"os"
)

func process(env *l.LangEnv, line string) {
	evalResult := l.Eval(line, env)
	tokens := evalResult.RemainingTokens
	if len(evalResult.ErrStr) > 0 {
		fmt.Printf("Error: %s\n", evalResult.ErrStr)
	} else {
		if len(evalResult.ValStr) > 0 {
			fmt.Printf("%s\n", evalResult.ValStr)
		} else {
			fmt.Printf("\n")
		}

		if len(tokens) > 0 {
			process(env, tokens)
		}
	}
}

func initREPL() {
	// Setup the language environment
	env := l.NewEnv()

	prompt := "lambda> "
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

func processScriptFile(scriptFilePath string) {
	file, err := os.Open(scriptFilePath)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	var concBuf bytes.Buffer
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		concBuf.WriteString(scanner.Text())
		concBuf.WriteString(" ")
	}

	if scanner.Err() != nil {
		panic(scanner.Err())
	}
	process(l.NewEnv(), concBuf.String())
}

func main() {
	var scriptFile = flag.String("f", "", "path of the file to read from")
	flag.Parse()

	if len(*scriptFile) > 0 {
		processScriptFile(*scriptFile)
	} else {
		initREPL()
	}
}
