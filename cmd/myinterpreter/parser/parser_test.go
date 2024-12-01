package parser

import (
	"testing"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

func TestPrinter(t *testing.T) {
	var exp expression.Expression = expression.NewBinaryExpression(
		expression.NewLiteralExpression(
			token.NewToken(token.NUMBER, 1, "100", token.NewNumValue(100)),
		),
		token.NewToken(token.PLUS, 1, "+", token.NewNullValue()),
		expression.NewLiteralExpression(
			token.NewToken(token.NUMBER, 1, "100", token.NewNumValue(100)),
		),
	)
	printer := NewAstPrinter()
	res := printer.Print(exp)
	match := "(+ 100.0 100.0)"
	if res != match {
		t.Errorf("TestPaserPrinter Error, get: %s, want: %s", res, match)
	}
}

func TestPrinterComplex(t *testing.T) {
	// Create the expression: (- 123) * (45.67)
	expression := expression.NewBinaryExpression(
		expression.NewUnaryExpression(
			token.NewToken(token.MINUS, 1, "-", token.NewNullValue()),
			expression.NewLiteralExpression(
				token.NewToken(token.NUMBER, 1, "123", token.NewNumValue(123)),
			),
		),
		token.NewToken(token.STAR, 1, "*", token.NewNullValue()),
		expression.NewGroupingExpression(
			expression.NewLiteralExpression(
				token.NewToken(token.NUMBER, 1, "45.67", token.NewNumValue(45.67)),
			),
		),
	)

	printer := NewAstPrinter()
	result := printer.Print(expression)
	expected := "(* (- 123.0) (group 45.67))"

	if result != expected {
		t.Errorf("TestParserPrinterComplex Error, got: %s, want: %s", result, expected)
	}
}

func TestParser(t *testing.T) {
	// Create the expression: (- 123) * (45.67)
	tokens := []*token.Token{
		token.NewToken(token.NUMBER, 1, "1", token.NewNumValue(1)),
		token.NewToken(token.PLUS, 1, "+", token.NewNullValue()),
		token.NewToken(token.NUMBER, 1, "1", token.NewNumValue(1)),
		token.NewToken(token.EOF, 1, "", token.NewNullValue()),
	}
	parser := New(tokens)
	expression, errs := parser.Parse()
	printer := NewAstPrinter()
	result := printer.Print(expression)
	expected := "(+ 1.0 1.0)"

	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestParser non nil error %v", errs)
	}
}

func TestLexParser(t *testing.T) {
	lex := lexer.New("true")
	lex.Lex()
	tokens := lex.Tokens()
	parser := New(tokens)
	expression, errs := parser.Parse()
	printer := NewAstPrinter()
	result := printer.Print(expression)
	expected := "true"
	
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestParser non nil error %v", errs)
	}
}

func TestLexNilParser(t *testing.T) {
	lex := lexer.New("(nil)")
	lex.Lex()
	tokens := lex.Tokens()
	parser := New(tokens)
	expression, errs := parser.Parse()
	printer := NewAstPrinter()
	result := printer.Print(expression)
	expected := "(group nil)"
	
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestParser non nil error %v", errs)
	}
}
