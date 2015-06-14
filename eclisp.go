package main

import (
	"bufio"
	"fmt"
	"os"
);

func Readline() []byte {
	fmt.Printf("eclisp> ")
	bio := bufio.NewReader(os.Stdin)
	line, _, _ := bio.ReadLine()
	// TODO
	// handle error
	return line
}

func main() {
	Readline()
}
