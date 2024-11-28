package token

import (
	"fmt"
	"strconv"
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
	IntValue    = "int"
	FloatValue  = "float"
	NullValue   = "null"
)

type TokenValue struct {
	Type        TokenValueType
	valueInt    int
	valueString string
	valueFloat  float64
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

func NewIntValue(num int) *TokenValue {
	return &TokenValue{
		Type:     IntValue,
		valueInt: num,
	}
}

func NewStringValue(str string) *TokenValue {
	return &TokenValue{
		Type:        StringValue,
		valueString: str,
	}
}

func NewFloatValue(floatNum float64) *TokenValue {
	return &TokenValue{
		Type:       FloatValue,
		valueFloat: floatNum,
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
	case IntValue:
		return strconv.Itoa(v.valueInt)
	}
	return ""
}

func (t Token) String() string {
	switch t.Type {
	case NUMBER:
		return fmt.Sprint("%s %s %s", t.Type, t.Text, t.Text)
	default:
		return fmt.Sprint("%s %s %s", t.Type, t.Text, t.TokenValue.String())
	}
}
