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
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
	interpreter := New()
	_, errs = interpreter.Eval(expression)
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %v", errs)
	}
	result := interpreter.String()
	expected := "-1"
	if result != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", result, expected)
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
	program, errs := p.ParseProgram()
	interpreter := New()
	_, errs = interpreter.Interp(program)
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

func TestInterpVarProgram(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		var a = 1;
		var b = 2;
		print a + b;
`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, program, parser.NewAstPrinter().PrintProgram(program))
	expected := "3\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestInterpVarAssProgram(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		var a = 1;
		a = 2;
		print a;
`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, program, parser.NewAstPrinter().PrintProgram(program))
	expected := "2\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestInterpBlockVarProgram(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		var a = 1;
		{
			a = 2;
			print a;
		}
`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, program, parser.NewAstPrinter().PrintProgram(program))
	expected := "2\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestInterpIfProgram(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		var truth = true;
		if (truth) 
			print "if";
		
		if (!truth)
			print "if";
		else print "else";
		
`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, program, parser.NewAstPrinter().PrintProgram(program))
	expected := "if\nelse\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestLogicalOperators(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`print false or "hi";`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	expected := "hi\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestWhileStmt(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		var i = 0;
		while (i < 4)
			i = i +1;
		print i;
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	expected := "4\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestForStmt(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		var i = 0;
		for (var i = 0; i < 4; i = i+1)
			print i;
		print i;
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	expected := "0\n1\n2\n3\n0\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestNativeTimeFunctionStmt(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		print clock();
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	if len(res) < 5 {
		t.Errorf("TestParser Error, got: %s, want: %v", res, 5)
	}
}

func TestFunctionDeclarationStmt(t *testing.T) {

	lex := lexer.New(`
		fun foo(){
		
		}
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
}

func TestReturnStmt(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		fun foo(variable){
			print variable;
			var i = 0;
			while (true) {
				if (i >= 4){
					return i;
				}
				i = i + 1;
			}
		}
		print foo(10);
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	expected := "10\n4\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestClosuresStmt(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		fun makeCounter() {
		  var i = 0;
		  fun count() {
		    i = i + 1;
		    print i;
		  }
		
		  return count;
		}
		var counter = makeCounter();
		counter();
		counter();
		counter();
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	expected := "1\n2\n3\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestFibStmt(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		fun returnArg(arg) {
			return arg;
		}
		
		fun returnFunCallWithArg(func, arg) {
			return returnArg(func)(arg);
		}
		
		fun printArg(arg) {
		print arg;
		}
		
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	expected := ""
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}

func TestReturnControllFlowStmt(t *testing.T) {
	// mock stdout
	rescueStdout := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	lex := lexer.New(`
		fun f() {
   			while (!true) return "ok";
      	}
       	print f();
	`)
	lex.Lex()
	tokens := lex.Tokens()
	p := parser.New(tokens)
	program, errs := p.ParseProgram()
	if errs != nil {
		t.Errorf("TestInterpreter non nil error %s", errs)
		return
	}
	interpreter := New()
	_, errs = interpreter.Interp(program)
	// demock stdout
	w.Close()
	out, _ := io.ReadAll(r)
	os.Stdout = rescueStdout
	res := string(out)
	fmt.Println(errs, parser.NewAstPrinter().PrintProgram(program))
	expected := "nil\n"
	if res != expected {
		t.Errorf("TestParser Error, got: %s, want: %s", res, expected)
	}
}
