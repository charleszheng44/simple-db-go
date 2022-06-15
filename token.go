package main

import (
	"strconv"
	"strings"
	"unicode"

	"github.com/pkg/errors"
)

type KeyWord int

const (
	Create KeyWord = iota
	Select
	Insert
	Delete
	Into
	From
	Table
	Where
	LeftParen
	RightParen
	Comma
	Equal
)

const Invalid = KeyWord(-1)

func (kw KeyWord) String() string {
	switch kw {
	case Create:
		return "create"
	case Select:
		return "select"
	case Insert:
		return "insert"
	case Delete:
		return "delete"
	case Into:
		return "into"
	case From:
		return "from"
	case Table:
		return "table"
	case Where:
		return "where"
	case LeftParen:
		return "leftparen"
	case RightParen:
		return "rightparen"
	case Comma:
		return "comma"
	case Equal:
		return "equal"
	}
	return "invalid"
}

type empty struct{}

var null = struct{}{}

var keyWords map[string]empty = map[string]empty{
	Create.String():     null,
	Select.String():     null,
	Insert.String():     null,
	Delete.String():     null,
	Into.String():       null,
	From.String():       null,
	Table.String():      null,
	Where.String():      null,
	LeftParen.String():  null,
	RightParen.String(): null,
	Comma.String():      null,
	Equal.String():      null,
}

func StringToKeyWord(str string) (KeyWord, error) {
	switch str {
	case "create":
		return Create, nil
	case "select":
		return Select, nil
	case "insert":
		return Insert, nil
	case "delete":
		return Delete, nil
	case "into":
		return Into, nil
	case "from":
		return From, nil
	case "table":
		return Table, nil
	case "where":
		return Where, nil
	case "leftParen":
		return LeftParen, nil
	case "rightParen":
		return RightParen, nil
	case "comma":
		return Comma, nil
	case "equal":
		return Equal, nil
	}
	return Invalid, errors.New("unknown keywrds")
}

type TokenType int

const (
	KeyWordToken TokenType = iota
	IntegerToken
	FloatToken
	BoolToken
	StringToken
)

type Token struct {
	Type       TokenType
	KeyWordVal KeyWord
	IntegerVal int
	FloatVal   float64
	BoolVal    bool
	StringVal  string
}

func (tk Token) String() string {
	switch tk.Type {
	case KeyWordToken:
		return tk.KeyWordVal.String()
	case IntegerToken:
		return strconv.FormatInt(int64(tk.IntegerVal), 10)
	case FloatToken:
		return strconv.FormatFloat(tk.FloatVal, 'f', 8, 64)
	case BoolToken:
		return strconv.FormatBool(tk.BoolVal)
	case StringToken:
		return tk.StringVal
	}
	return "invalid"
}

func isKeyWord(word string) (KeyWord, bool) {
	var kw KeyWord
	_, exist := keyWords[strings.ToLower(word)]
	if !exist {
		return kw, false
	}
	kw, err := StringToKeyWord(string(word))
	if err != nil {
		// TODO(charleszheng44): print the error message
		return kw, false
	}

	return kw, true
}

func isNumber(word string) (*Token, bool) {
	var floatNum bool
	for _, r := range word {
		if !unicode.IsDigit(r) {
			if r == '.' && !floatNum {
				// meet the '.' for the first time
				floatNum = true
				continue
			}
			return nil, false
		}
	}

	if floatNum {
		fn, err := strconv.ParseFloat(word, 64)
		if err != nil {
			// TODO(charleszheng44): print the error message
			return nil, false
		}
		return &Token{
			Type:     FloatToken,
			FloatVal: fn,
		}, true
	}

	n, err := strconv.Atoi(word)
	if err != nil {
		// TODO(charleszheng44): print the error message
		return nil, false
	}
	return &Token{
		Type:       IntegerToken,
		IntegerVal: n,
	}, true
}

func isBool(word string) (*Token, bool) {
	if word == "false" {
		return &Token{
			Type:    BoolToken,
			BoolVal: false,
		}, true
	}

	if word == "true" {
		return &Token{
			Type:    BoolToken,
			BoolVal: true,
		}, true
	}
	return nil, false
}

func tokenize(word string) (*Token, error) {
	if tk, ok := isBool(word); ok {
		return tk, nil
	}

	if tk, ok := isNumber(word); ok {
		return tk, nil
	}

	return &Token{
		Type:      StringToken,
		StringVal: word,
	}, nil
}

func Tokenize(inp string) ([]*Token, error) {
	var tks []*Token
	// process the keywords first
	var currWord []rune
	for _, rn := range inp {
		// ignore the space
		if unicode.IsSpace(rn) {
			tk, err := tokenize(string(currWord))
			if err != nil {
				return tks, err
			}
			tks = append(tks, tk)
			// reset for the next word
			currWord = []rune{}
			continue
		}

		if rn == '(' {
			tk, err := tokenize(string(currWord))
			if err != nil {
				return tks, err
			}
			tks = append(tks, tk)
			// reset for the next word
			currWord = []rune{}
			tks = append(tks, &Token{
				Type:       KeyWordToken,
				KeyWordVal: LeftParen,
			})
			continue
		}

		if rn == ')' {
			tk, err := tokenize(string(currWord))
			if err != nil {
				return tks, err
			}
			tks = append(tks, tk)
			// reset for the next word
			currWord = []rune{}
			tks = append(tks, &Token{
				Type:       KeyWordToken,
				KeyWordVal: RightParen,
			})
			continue
		}

		if rn == ',' {
			tk, err := tokenize(string(currWord))
			if err != nil {
				return tks, err
			}
			tks = append(tks, tk)
			// reset for the next word
			currWord = []rune{}
			tks = append(tks, &Token{
				Type:       KeyWordToken,
				KeyWordVal: Comma,
			})
			continue
		}

		currWord = append(currWord, rn)
	}
	tk, err := tokenize(string(currWord))
	if err != nil {
		return tks, err
	}
	tks = append(tks, tk)

	return tks, nil
}
