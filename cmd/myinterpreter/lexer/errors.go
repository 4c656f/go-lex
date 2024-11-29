package lexer

import "fmt"

type LexError struct {
	line int
	where string
	message string
}

func NewLexError(line int, where string, message string)error{
	return &LexError{
		line: line,
		where: where,
		message: message,
	}
}

func (e LexError)Error()string{
	if e.message == ""{
		return fmt.Sprintf("[line %v] Error: %s", e.line, e.where)
	}
	return fmt.Sprintf("[line %v] Error: %s: %s", e.line, e.where, e.message)
}
