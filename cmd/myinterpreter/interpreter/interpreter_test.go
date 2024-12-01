package interpreter

import (
	"fmt"
	"io"
	"os"
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

func TestInterpLiteral(t *testing.T) {
	lex := lexer.New("21")
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.Parse()
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	result := interpreter.String()
	fmt.Println(result, parser.NewAstPrinter().Print(expression))
	expected := "21"
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
}

func TestInterpBoolLiteral(t *testing.T) {
	lex := lexer.New("true")
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.Parse()
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	result := interpreter.String()
	fmt.Println(result, parser.NewAstPrinter().Print(expression))
	expected := "true"
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
}

func TestInterpNegIntLiteral(t *testing.T) {
	lex := lexer.New("-43")
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.Parse()
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	result := interpreter.String()
	fmt.Println(result, parser.NewAstPrinter().Print(expression))
	expected := "-43"
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
}

func TestInterpNilUrnaryLiteral(t *testing.T) {
	lex := lexer.New("!nil")
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.Parse()
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	result := interpreter.String()
	fmt.Println(result, parser.NewAstPrinter().Print(expression))
	expected := "true"
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
}

func TestInterpDivOp(t *testing.T) {
	lex := lexer.New("7/5")
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.Parse()
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	result := interpreter.String()
	fmt.Println(result, parser.NewAstPrinter().Print(expression))
	expected := "1.4"
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
	}
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
}

func TestInterpRuntimeErrors(t *testing.T) {
	lex := lexer.New(`-"world"`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.Parse()
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	result := interpreter.String()
	fmt.Println(result, parser.NewAstPrinter().Print(expression))
	fmt.Println(errs)
	if errs == nil {
		t.Errorf("TestParser does not had runtime Error, got: %v", errs)
	}
}

func TestInterpProgram(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w
	
	lex := lexer.New(`print true;`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	expression, errs := p.ParseProgram()
	interpreter := New()
	_, errs = interpreter.Interp(expression)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs)
	expected := "true\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}
