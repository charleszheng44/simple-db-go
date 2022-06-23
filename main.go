package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	var (
		isStr  bool
		stsTks [][]rune
		curSts []rune
	)
	db := NewDatabase()
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("@simple-db=> ")
	for {
		r, _, err := reader.ReadRune()
		if err != nil {
			fmt.Printf("[ERROR] failed to read from stdin: %v\n", err)
		}
		if r == '\'' {
			isStr = !isStr
		}

		if r == ';' && !isStr {
			stsTks = append(stsTks, curSts)
			curSts = []rune{}
			continue
		}

		if r == '\n' {
			for _, stk := range stsTks {
				tks, err := Tokenize(stk)
				if err != nil {
					fmt.Printf("[ERROR] failed to tokenize the input string: %v\n", err)
					continue
				}

				sts, err := parse(tks)
				if err != nil {
					fmt.Printf("[ERROR] failed to parse the input tokens: %v\n", err)
					continue
				}

				result := db.Interpret(sts)
				if result.err != nil {
					fmt.Printf("[ERROR] failed to interpret "+
						"the statement: %v\n", err)
					continue
				}
				if len(result.message) != 0 {
					fmt.Println(result.message)
				}
				if len(result.rows) != 0 {
					fmt.Println(result.rows)
				}
			}
			stsTks = nil

			if len(curSts) == 0 {
				fmt.Printf("@simple-db=> ")
				continue
			}

			if isStr {
				fmt.Printf("@simple-db'> ")
				continue
			}

			fmt.Printf("@simple-db-> ")
		}
		curSts = append(curSts, r)
	}
}
