package main

import (
	"fmt"
	"os"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/cli"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/interpreter"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/lexer"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/parser"
)

func main() {
	arg, err := cli.ParseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
	switch arg.Type {
	case cli.Tokenize:
		tokenize(arg.FilePath)
	case cli.Parse:
		parse(arg.FilePath)
	case cli.Eval:
		eval(arg.FilePath)
	}
}

func tokenize(fileName string) {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	lex := lexer.New(string(fileContents))
	errs := lex.Lex()
	if errs != nil {
		for _, err := range errs {
			lexer.Report(err)
		}
	}
	fmt.Println(lex.String())
	if errs != nil {
		os.Exit(65)
	}
}

func parse(fileName string) {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	lex := lexer.New(string(fileContents))
	errs := lex.Lex()
	if errs != nil {
		for _, err := range errs {
			lexer.Report(err)
		}
	}
	par := parser.New(lex.Tokens())
	exp, errs := par.Parse()
	printer := parser.NewAstPrinter()
	fmt.Println(printer.Print(exp))
	if errs != nil {
		for _, err := range errs {
			lexer.Report(err)
		}
		os.Exit(65)
	}
}

func eval(fileName string) {
	fileContents, err := os.ReadFile(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}
	lex := lexer.New(string(fileContents))
	errs := lex.Lex()
	if errs != nil {
		for _, err := range errs {
			lexer.Report(err)
		}
	}
	par := parser.New(lex.Tokens())
	exp, errs := par.Parse()
	if errs != nil {
		for _, err := range errs {
			lexer.Report(err)
		}
		os.Exit(70)
	}
	fmt.Println(parser.NewAstPrinter().Print(exp))
	interp := interpreter.New()
	_, errs = interp.Eval(exp)
	fmt.Println(interp.String())
	if errs != nil {
		for _, err := range errs {
			lexer.Report(err)
		}
		os.Exit(70)
	}
}
