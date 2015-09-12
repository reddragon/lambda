package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	l "github.com/reddragon/eclisp/lang"
	"github.com/tiborvass/uniline"
	"os"
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

func initREPL() {
	// Setup the language environment
	env := l.NewEnv()

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
