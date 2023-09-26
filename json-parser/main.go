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
		return token
	}

	if l.content[l.cursor] == BEGIN_OBJECT {
		l.cursor += 1
		token.kind = "BEGIN_OBJECT"
		return token
	}

	if l.content[l.cursor] == END_OBJECT {
		l.cursor += 1
		token.kind = "END_OBJECT"
		return token
	}

	if l.content[l.cursor] == QUOTATION_MARK {
		l.cursor += 1
		token.kind = "QUOTATION_MARK"
		return token
	}

	if l.content[l.cursor] == NAME_SEPERATOR {
		l.cursor += 1
		token.kind = "NAME_SEPERATOR"
		return token
	}

	if l.content[l.cursor] == VALUE_SEPERATOR {
		l.cursor += 1
		token.kind = "VALUE_SEPERATOR"
		return token
	}

	if l.content[l.cursor] == BEGIN_ARRAY {
		l.cursor += 1
		token.kind = "BEGIN_ARRAY"
		return token
	}

	if l.content[l.cursor] == END_ARRAY {
		l.cursor += 1
		token.kind = "END_ARRAY"
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
		if lexer.cursor == 1 && token.kind != "BEGIN_OBJECT" {
			fmt.Println("Invalid json!")
			os.Exit(1)
		}
		if token.kind == "END_OBJECT" {
			break
		}

		fmt.Println("Token", token.kind, "value", string(token.text))
	}
	fmt.Println("Hello lexer", filename)
}
