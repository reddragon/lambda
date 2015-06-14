package main

import (
	// "bufio"
	"fmt"
	// "os"
	"github.com/tiborvass/uniline"
);


func main() {
	prompt := "eclisp> "
	scanner := uniline.DefaultScanner()
	for scanner.Scan(prompt) {
		line := scanner.Text()
		if len(line) > 0 {
			scanner.AddToHistory(line)
			fmt.Println(line)
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}
