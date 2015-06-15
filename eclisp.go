package main

import (
	// "bufio"
	"fmt"
	"strings"
	"strconv"
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

// TODO Use only err as error string
type Atom struct {
	err bool
	errMsg string
	val int
}

const (
	// Delimiters
	OpenBracket string = "("
	ClosedBracket string = ")"

	// Operators
	Add string = "+"
	Sub string = "-"
	Mul string = "*"
	Div string = "/"
)

func add(operand1, operand2 Atom) Atom {
	var retVal Atom
	retVal.val = operand1.val + operand2.val
	return retVal
}

func isValidOperator(operator string) bool {
	if (operator == Add) {
		return true
	}
	return false
}

func pop(tokens []string) (string, []string) {
	if len(tokens) == 0 {
		return "", tokens
	}
	return tokens[0], tokens[1:]
}

func evalArgs(tokens []string) (Atom, []string) {
	var retVal Atom
	retVal.err = false
	retVal.errMsg = ""
	
	// Expecting a value, but if a ')' occurs, its an error
	if tokens[0] == ClosedBracket {
		retVal.errMsg = "Unexpected ')'"
		retVal.err = true
		return retVal, tokens
	}
	
	// This could be a nested expression
	if tokens[0] == OpenBracket {
		return eval(tokens)
	}

	// This is an atomic value
	token := ""
	token, tokens = pop(tokens)
	val, err := strconv.Atoi(token)
	if err != nil {
		retVal.errMsg = "Problem while converting " + token + ": " + err.Error()
		retVal.err = true
		return retVal, tokens
	}
	retVal.val = val
	fmt.Printf("Tokens after: %q\n", tokens)
	return retVal, tokens
}

func eval(tokens []string) (Atom, []string) {
	var retVal Atom
	retVal.err = false
	retVal.errMsg = ""

	token := ""
	token, tokens = pop(tokens)
	if token != OpenBracket {
		// TODO
		// Raise an exception
		retVal.errMsg = "Expected a '('";
		retVal.err = true
		return retVal, tokens
	}
	
	token, tokens = pop(tokens)
	validOperator := isValidOperator(token)
	if !validOperator {
		retVal.errMsg = "Invalid operator, " + token
		retVal.err = true
		return retVal, tokens
	}
	
	// Read two args. 
	// TODO change this to read as many args as the operator wants

	// TODO
	// Reading just two ints/floats here
	// Change to read two nested expressions here too
	
	var operand1, operand2 Atom
	operand1, tokens = evalArgs(tokens)
	if operand1.err {
		return operand1, tokens
	}

	operand2, tokens = evalArgs(tokens)
	if operand2.err {
		return operand2, tokens
	}

	fmt.Printf("Values of operands: %d, %d\n", operand1.val, operand2.val);
	
	token, tokens = pop(tokens)
	if token != ClosedBracket {
		retVal.errMsg = "Expected a ')'";
		retVal.err = true
		return retVal, tokens
	}
	
	retVal.val = add(operand1, operand2).val
	return retVal, tokens
}

func process(line string) {
	tokens := tokenize(line);
	if len(tokens) == 0 {
		fmt.Println("Nothing to evaluate");
		return;
	}
	fmt.Printf("%q\n", tokens)
	var retVal Atom
	retVal, tokens = eval(tokens)

	if len(tokens) != 0 {
		retVal.errMsg = strconv.Itoa(len(tokens)) + " extra token(s) in the string."
		retVal.err = true
	}
	if retVal.err {
		fmt.Printf("Error: %s\n", retVal.errMsg)
	} else {
		fmt.Println("All worked fine!")
		fmt.Printf("Result: %d\n", retVal.val)
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
