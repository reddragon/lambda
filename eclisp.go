package main

import (
	// "bufio"
	"fmt"
	"strings"
	// "os"
	"github.com/tiborvass/uniline"
);

func tokenize(line string) []string {
	tokens := strings.Split(strings.Replace(strings.Replace(line, "(", " ( ", -1), ")", " ) ", -1), " ")
	noWSTokens := make([]string, 0)
	for _, token := range tokens {
		if token != "" {
			noWSTokens = append(noWSTokens, token)
		}
	}
	return noWSTokens
}

func process(line string) {
	tokens := tokenize(line);
	if len(tokens) == 0 {
		fmt.Println("Nothing to evaluate");
		return;
	}
	fmt.Printf("%q\n", tokens)
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
