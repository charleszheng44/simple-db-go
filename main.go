package main

import (
	"bufio"
	"fmt"
	"os"
)

const delimiter = byte(';')

func main() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Printf("@simple-db=> ")
		inpSts, err := reader.ReadBytes(delimiter)
		if err != nil {
			fmt.Printf("[ERROR] failed to read from stdin: %v\n", err)
			continue
		}
		tks, err := Tokenize(string(inpSts))
		if err != nil {
			fmt.Printf("[ERROR] failed to tokenize the input string: %v\n", err)
			continue
		}
		sts, err := parse(tks)
		if err != nil {
			fmt.Printf("[ERROR] failed to parse the input tokens: %v\n", err)
			continue
		}
		result, err := sts.Interpret()
		if err != nil {
			fmt.Printf("[ERROR] failed to interpret the statement: %v\n", err)
			continue
		}
		if result != "" {
			fmt.Println(result)
		}
	}
}
