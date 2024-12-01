package token

import (
	"fmt"
)

type TokenType string

const (
	// Single-character tokens.
	LEFT_PAREN  TokenType = "LEFT_PAREN"
	RIGHT_PAREN TokenType = "RIGHT_PAREN"
	LEFT_BRACE  TokenType = "LEFT_BRACE"
	RIGHT_BRACE TokenType = "RIGHT_BRACE"
	COMMA       TokenType = "COMMA"
	DOT         TokenType = "DOT"
	MINUS       TokenType = "MINUS"
	PLUS        TokenType = "PLUS"
	SEMICOLON   TokenType = "SEMICOLON"
	SLASH       TokenType = "SLASH"
	STAR        TokenType = "STAR"

	// One or two character tokens.
	BANG          TokenType = "BANG"
	BANG_EQUAL    TokenType = "BANG_EQUAL"
	EQUAL         TokenType = "EQUAL"
	EQUAL_EQUAL   TokenType = "EQUAL_EQUAL"
	GREATER       TokenType = "GREATER"
	GREATER_EQUAL TokenType = "GREATER_EQUAL"
	LESS          TokenType = "LESS"
	LESS_EQUAL    TokenType = "LESS_EQUAL"

	// Literals.
	IDENTIFIER TokenType = "IDENTIFIER"
	STRING     TokenType = "STRING"
	NUMBER     TokenType = "NUMBER"

	// Keywords.
	AND    TokenType = "AND"
	CLASS  TokenType = "CLASS"
	ELSE   TokenType = "ELSE"
	FALSE  TokenType = "FALSE"
	FUN    TokenType = "FUN"
	FOR    TokenType = "FOR"
	IF     TokenType = "IF"
	NIL    TokenType = "NIL"
	OR     TokenType = "OR"
	PRINT  TokenType = "PRINT"
	RETURN TokenType = "RETURN"
	SUPER  TokenType = "SUPER"
	THIS   TokenType = "THIS"
	TRUE   TokenType = "TRUE"
	VAR    TokenType = "VAR"
	WHILE  TokenType = "WHILE"

	EOF TokenType = "EOF"
)

type TokenValueType string

const (
	StringValue = "string"
	NumValue    = "num"
	BoolValue   = "bool"
	NullValue   = "null"
)

type TokenValue struct {
	Type        TokenValueType
	valueNum    float64
	valueString string
	valueBool   bool
}

type Token struct {
	Type       TokenType
	Line       int
	Text       string
	TokenValue *TokenValue
}

func NewToken(tType TokenType, line int, text string, value *TokenValue) *Token {
	return &Token{
		Type:       tType,
		Line:       line,
		Text:       text,
		TokenValue: value,
	}
}

func NewNumValue(num float64) *TokenValue {
	return &TokenValue{
		Type:     NumValue,
		valueNum: num,
	}
}

func NewBoolValue(b bool) *TokenValue {
	return &TokenValue{
		Type:      BoolValue,
		valueBool: b,
	}
}

func NewStringValue(str string) *TokenValue {
	return &TokenValue{
		Type:        StringValue,
		valueString: str,
	}
}

func NewNullValue() *TokenValue {
	return &TokenValue{
		Type: NullValue,
	}
}

func (v TokenValue) String() string {
	switch v.Type {
	case StringValue:
		return v.valueString
	case NullValue:
		return "null"
	case BoolValue:
		return "null"
	case NumValue:
		if v.valueNum == float64(int(v.valueNum)) {
			return fmt.Sprintf("%.1f", v.valueNum)
		}
		return fmt.Sprintf("%g", v.valueNum)
	}
	return ""
}

func (v TokenValue) GetValue() any {
	switch v.Type {
	case StringValue:
		return v.valueString
	case NullValue:
		return nil
	case NumValue:
		return v.valueNum
	case BoolValue:
		return v.valueBool
	}
	return nil
}

var stringToKeywoard = map[string]TokenType{
	"and":    AND,
	"class":  CLASS,
	"else":   ELSE,
	"false":  FALSE,
	"for":    FOR,
	"fun":    FUN,
	"if":     IF,
	"nil":    NIL,
	"or":     OR,
	"return": RETURN,
	"super":  SUPER,
	"this":   THIS,
	"true":   TRUE,
	"var":    VAR,
	"while":  WHILE,
	"print":  PRINT,
}

func MatchStringToKeywoard(s string) (TokenType, bool) {
	t, ok := stringToKeywoard[s]
	return t, ok
}

func (t Token) String() string {
	return fmt.Sprintf("%s %s %s", string(t.Type), t.Text, t.TokenValue.String())
}
