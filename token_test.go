package main

import (
	"reflect"
	"testing"
)

var (
	SelectTk = &Token{
		Type:       KeyWordToken,
		KeyWordVal: Select,
	}

	StarTk = &Token{
		Type:       KeyWordToken,
		KeyWordVal: Star,
	}

	FromTk = &Token{
		Type:       KeyWordToken,
		KeyWordVal: From,
	}

	WhereTk = &Token{
		Type:       KeyWordToken,
		KeyWordVal: Where,
	}

	EqualTk = &Token{
		Type:       KeyWordToken,
		KeyWordVal: Equal,
	}
)

func unQuoteStrTk(inp string) *Token {
	return &Token{
		Type:      UnquoteStringToken,
		StringVal: inp,
	}
}

func stringTk(inp string) *Token {
	return &Token{
		Type:      StringToken,
		StringVal: inp,
	}
}

func intTk(inp int) *Token {
	return &Token{
		Type:       IntegerToken,
		IntegerVal: inp,
	}
}

func TestTokenize(t *testing.T) {
	tts := []struct {
		name   string
		input  string
		expect []*Token
	}{
		{
			"Select Statment",
			"select * from test",
			[]*Token{
				SelectTk,
				StarTk,
				FromTk,
				unQuoteStrTk("test"),
			},
		},
		{
			"Select Statment with Where Clause",
			"select * from test where id = 1",
			[]*Token{
				SelectTk,
				StarTk,
				FromTk,
				unQuoteStrTk("test"),
				WhereTk,
				unQuoteStrTk("id"),
				EqualTk,
				intTk(1),
			},
		},
		{
			"Select Statment with string and unquoteStr",
			"select * from test where name = 'test-name'",
			[]*Token{
				SelectTk,
				StarTk,
				FromTk,
				unQuoteStrTk("test"),
				WhereTk,
				unQuoteStrTk("name"),
				EqualTk,
				stringTk("test-name"),
			},
		},
	}

	for i, tt := range tts {
		tt := tt
		i := i
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := Tokenize([]rune(tt.input))
			if err != nil {
				t.Fatalf("case %d (%s) failed: %v", i, tt.name, err)
			}

			if !reflect.DeepEqual(got, tt.expect) {
				t.Fatalf("case %d (%s) failed: got(%s), expect(%s)", i,
					tt.name, got, tt.expect)
			}
			t.Logf("case %d (%s) succeed", i, tt.name)
		})
	}
}

func TestIsString(t *testing.T) {
	tts := []struct {
		name     string
		input    string
		expectTk *Token
		quoted   bool
	}{
		{
			"Quoted string",
			"'str'",
			&Token{
				Type:      StringToken,
				StringVal: "str",
			},
			true,
		},
		{
			"Quoted string",
			"str",
			nil,
			false,
		},
	}

	for i, tt := range tts {
		tt := tt
		i := i
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, quoted := isString(tt.input)
			if quoted != tt.quoted {
				t.Fatalf("case %d (%s) failed: "+
					"got quoted(%v), expect quoted(%v)",
					i, tt.name, quoted, tt.quoted)
			}

			if !reflect.DeepEqual(got, tt.expectTk) {
				t.Fatalf("case %d (%s) failed: got(%s), expect(%s)", i,
					tt.name, got, tt.expectTk)
			}
			t.Logf("case %d (%s) succeed", i, tt.name)
		})
	}
}
