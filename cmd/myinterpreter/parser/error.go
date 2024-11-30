package parser

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type ParserError struct {
	t       *token.Token
	message string
}

func NewParserError(t *token.Token, message string) *ParserError {
	return &ParserError{
		t:       t,
		message: message,
	}
}

func (e ParserError) Error() string {
	if e.t.Type == token.EOF {
		return fmt.Sprintf("%v at end %s", e.t.Line, e.message)
	}
	return fmt.Sprintf("%v at '%s'%s", e.t.Line, e.t.Text, e.message)
}
