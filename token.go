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
	Primary
	Key
	KeyWordTable
	Where
	LeftParen
	RightParen
	SingleQuote
	Comma
	Equal
	Star
)

var (
	TokenWhere = Token{
		Type:       KeyWordToken,
		KeyWordVal: Where,
	}

	TokenFrom = Token{
		Type:       KeyWordToken,
		KeyWordVal: From,
	}

	TokenPrimary = Token{
		Type:       KeyWordToken,
		KeyWordVal: Primary,
	}

	TokenKey = Token{
		Type:       KeyWordToken,
		KeyWordVal: Key,
	}

	TokenComma = Token{
		Type:       KeyWordToken,
		KeyWordVal: Comma,
	}

	TokenEqual = Token{
		Type:       KeyWordToken,
		KeyWordVal: Equal,
	}

	TokenTable = Token{
		Type:       KeyWordToken,
		KeyWordVal: KeyWordTable,
	}

	TokenLeftParen = Token{
		Type:       KeyWordToken,
		KeyWordVal: LeftParen,
	}

	TokenRightParen = Token{
		Type:       KeyWordToken,
		KeyWordVal: RightParen,
	}
)

func isUnquoteStringToken(token *Token) bool {
	return token.Type == UnquoteStringToken
}

func cmpTks(tk1, tk2 Token) bool {
	if tk1.Type != tk2.Type {
		return false
	}
	switch tk1.Type {
	case KeyWordToken:
		return tk1.KeyWordVal == tk2.KeyWordVal
	case UnquoteStringToken:
		return tk1.StringVal == tk2.StringVal
	case IntegerToken:
		return tk1.IntegerVal == tk2.IntegerVal
	case FloatToken:
		return tk1.FloatVal == tk2.FloatVal
	case BoolToken:
		return tk1.BoolVal == tk2.BoolVal
	case StringToken:
		return tk1.StringVal == tk2.StringVal
	default:
		// TODO(charleszheng44): print error message
		return false
	}
}

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
	case Primary:
		return "primary"
	case Key:
		return "key"
	case KeyWordTable:
		return "table"
	case Where:
		return "where"
	case LeftParen:
		return "("
	case RightParen:
		return ")"
	case SingleQuote:
		return "'"
	case Comma:
		return ","
	case Equal:
		return "="
	case Star:
		return "*"
	}
	return "invalid"
}

type empty struct{}

var null = struct{}{}

var keyWords map[string]empty = map[string]empty{
	Create.String():       null,
	Select.String():       null,
	Insert.String():       null,
	Delete.String():       null,
	Into.String():         null,
	From.String():         null,
	Primary.String():      null,
	Key.String():          null,
	KeyWordTable.String(): null,
	Where.String():        null,
	LeftParen.String():    null,
	RightParen.String():   null,
	SingleQuote.String():  null,
	Comma.String():        null,
	Equal.String():        null,
	Star.String():         null,
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
	case "primary":
		return Primary, nil
	case "key":
		return Key, nil
	case "table":
		return KeyWordTable, nil
	case "where":
		return Where, nil
	case "(":
		return LeftParen, nil
	case ")":
		return RightParen, nil
	case "'":
		return SingleQuote, nil
	case ",":
		return Comma, nil
	case "=":
		return Equal, nil
	case "*":
		return Star, nil
	}
	return Invalid, errors.New("unknown keywrds")
}

type TokenType int

const (
	KeyWordToken TokenType = iota
	UnquoteStringToken
	IntegerToken
	FloatToken
	BoolToken
	StringToken
)

func (tt TokenType) String() string {
	switch tt {
	case KeyWordToken:
		return "KeyWord"
	case UnquoteStringToken:
		return "UnquoteString"
	case IntegerToken:
		return "Integer"
	case FloatToken:
		return "Float"
	case BoolToken:
		return "Bool"
	case StringToken:
		return "String"
	default:
		return "invalid"
	}
}

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
	case UnquoteStringToken:
		return tk.StringVal
	case IntegerToken:
		return strconv.FormatInt(int64(tk.IntegerVal), 10)
	case FloatToken:
		return strconv.FormatFloat(tk.FloatVal, 'f', 8, 64)
	case BoolToken:
		return strconv.FormatBool(tk.BoolVal)
	case StringToken:
		return "'" + tk.StringVal + "'"
	}
	return "invalid"
}

func isKeyWord(word string) (*Token, bool) {
	_, exist := keyWords[strings.ToLower(word)]
	if !exist {
		return nil, false
	}
	kw, err := StringToKeyWord(string(word))
	if err != nil {
		// TODO(charleszheng44): print the error message
		return nil, false
	}

	return &Token{
		Type:       KeyWordToken,
		KeyWordVal: kw,
	}, true
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

// isString checks if the input `word` is a quoted string, i.e., 'xxxx'.
func isString(word string) (*Token, bool) {
	rs := []rune(word)
	if rs[0] == '\'' && rs[len(rs)-1] == '\'' {
		return &Token{
			Type:      StringToken,
			StringVal: string(rs[1 : len(rs)-1]),
		}, true
	}
	return nil, false
}

func tokenize(word string) (*Token, error) {
	if tk, ok := isKeyWord(word); ok {
		return tk, nil
	}

	if tk, ok := isBool(word); ok {
		return tk, nil
	}

	if tk, ok := isNumber(word); ok {
		return tk, nil
	}

	if tk, ok := isString(word); ok {
		return tk, nil
	}

	return &Token{
		Type:      UnquoteStringToken,
		StringVal: word,
	}, nil
}

func Tokenize(inp []rune) ([]*Token, error) {
	var tks []*Token
	// process the keywords first
	var currWord []rune
	var isStr bool
	for _, rn := range inp {
		// ignore the space
		if unicode.IsSpace(rn) {
			if len(currWord) == 0 {
				continue
			}
			tk, err := tokenize(string(currWord))
			if err != nil {
				return tks, err
			}
			tks = append(tks, tk)
			// reset for the next word
			currWord = []rune{}
			continue
		}

		if rn == '(' && !isStr {
			if len(currWord) == 0 {
				tks = append(tks, &Token{
					Type:       KeyWordToken,
					KeyWordVal: LeftParen,
				})
				continue
			}
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

		if rn == ')' && !isStr {
			if len(currWord) == 0 {
				tks = append(tks, &Token{
					Type:       KeyWordToken,
					KeyWordVal: RightParen,
				})
				continue
			}
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

		if rn == ',' && !isStr {
			if len(currWord) == 0 {
				tks = append(tks, &Token{
					Type:       KeyWordToken,
					KeyWordVal: Comma,
				})
				continue
			}
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

		if rn == '\'' {
			// content within two single quotes is a string
			isStr = !isStr
		}

		currWord = append(currWord, rn)
	}
	if len(currWord) != 0 {
		tk, err := tokenize(string(currWord))
		if err != nil {
			return tks, err
		}
		tks = append(tks, tk)
	}

	return tks, nil
}
