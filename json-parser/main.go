package main

import (
	"fmt"
	"os"
	"unicode"
)

const (
	BEGIN_ARRAY     = '['
	BEGIN_OBJECT    = '{'
	END_ARRAY       = ']'
	END_OBJECT      = '}'
	NAME_SEPERATOR  = ':'
	VALUE_SEPERATOR = ','
	QUOTATION_MARK  = '"'
)

var LITERAL_TOKENS = map[rune]string{
	BEGIN_ARRAY:     "BEGIN_ARRAY",
	BEGIN_OBJECT:    "BEGIN_OBJECT",
	END_ARRAY:       "END_ARRAY",
	END_OBJECT:      "END_OBJECT",
	NAME_SEPERATOR:  "NAME_SEPERATOR",
	VALUE_SEPERATOR: "VALUE_SEPERATOR",
	QUOTATION_MARK:  "QUOTATION_MARK",
}

type Token struct {
	text []byte
	kind string
}

type Lexer struct {
	content     []byte
	content_len int
	cursor      int
}

func NewLexer(filename string) *Lexer {
	content, err := os.ReadFile(filename)
	if err != nil {
		panic(err)
	}

	return &Lexer{
		content:     content,
		content_len: len(content),
		cursor:      0,
	}
}

func (l *Lexer) Next() *Token {
	l.trimLeft()
	text := []byte{l.content[l.cursor]}
	token := &Token{
		text: text,
		kind: "SYMBOL",
	}
	if l.cursor >= l.content_len {
		l.cursor += 1
		return token
	}

	if token_kind, ok := LITERAL_TOKENS[rune(l.content[l.cursor])]; ok {
		l.cursor += 1
		token.kind = token_kind
		return token
	}

	if l.IsAlpha(rune(l.content[l.cursor])) {
		l.cursor += 1
		token.kind = "SYMBOL"
		for {
			if l.cursor > l.content_len || !l.IsAlpha(rune(l.content[l.cursor])) {
				break
			}
			token.text = append(token.text, l.content[l.cursor])
			l.cursor += 1
		}
		return token

	}

	l.cursor += 1
	token.kind = "INVALID_TOKEN"

	return token
}

func (l *Lexer) trimLeft() {
	for l.cursor < l.content_len && unicode.IsSpace(rune(l.content[l.cursor])) {
		l.cursor++
	}
}

func (l *Lexer) IsAlpha(r rune) bool {
	return unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-'
}

func main() {
	filename := os.Args[1]

	lexer := NewLexer(filename)
	if lexer.content_len == 0 {
		fmt.Println("Empty file!")
		os.Exit(1)
	}
	fmt.Println("parsing", lexer.content)
	for {
		token := lexer.Next()
		fmt.Println("Token", token.kind, "value", string(token.text))
		if lexer.cursor == 1 && token.kind != "BEGIN_OBJECT" {
			fmt.Println("Invalid json!")
			os.Exit(1)
		}
		if token.kind == "END_OBJECT" {
			break
		}

	}
	fmt.Println("Hello lexer", filename)
}
