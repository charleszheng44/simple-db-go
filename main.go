package main

import (
	"bufio"
	"fmt"
	"os"
)

// TODO (charleszheng44): dynamic column width
func printRows(rs []*Row, cols []string) {
	for _, col := range cols {
		fmt.Printf("|%s\t", col)
	}
	fmt.Println()
	for _, r := range rs {
		for _, col := range cols {
			fmt.Printf("|%v\t", r.fields[col])
		}
		fmt.Println()
	}
}

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
					printRows(result.rows, result.cols)
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
