package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func readLines(rd io.Reader, stsChan chan<- []rune) error {
	var (
		isStr  bool
		stsTks [][]rune
		curSts []rune
	)
	reader := bufio.NewReader(rd)
	fmt.Printf("@simple-db=>")
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			fmt.Printf("[ERROR] failed to read from stdin: %v\n", err)
		}
		if r == '\'' {
			// TODO(charleszheng44): explain
			isStr = !isStr
		}

		if r == ';' && !isStr {
			stsTks = append(stsTks, curSts)
			curSts = []rune{}
		}

		if r == '\n' {
			for _, stk := range stsTks {
				stsChan <- stk
			}
			stsTks = nil

			if len(curSts) == 0 {
				fmt.Println("@simple-db=> ")
				continue
			}

			if isStr {
				fmt.Println("@simple-db'> ")
				continue
			}

			fmt.Println("@simple-db-> ")
		}
		curSts = append(curSts, r)
	}
}

func main() {
	stsChan := make(chan []rune)
	go readLines(os.Stdin, stsChan)

	for sts := range stsChan {
		tks, err := Tokenize(sts)
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
