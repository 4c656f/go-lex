package lexer

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Lexer struct {
	source string
	start  int
	end    int
	line   int
	tokens []token.Token
}

func New(src string) *Lexer {
	return &Lexer{
		source: src,
	}
}

func Lex() error {
	return nil
}

func (l Lexer) isAtEnd() bool {
	return l.end >= len(l.source)
}
