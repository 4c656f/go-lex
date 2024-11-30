package interpreter

import (
	"fmt"
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
)

func TestInterpreter(t *testing.T) {
	lex := lexer.New("(1 + 1) - 3")
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.Parse()
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	result := interpreter.String()
	expected := "-1"
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
}
