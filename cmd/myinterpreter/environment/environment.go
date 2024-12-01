package environment

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Environment struct {
	variables map[string]any
}

func (e *Environment) Define(name string, value any) {
	e.variables[name] = value
}

func (e *Environment) Get(name *token.Token) (any, error) {
	v, ok := e.variables[name.Text]

	if !ok {
		return nil, errors.NewRuntimeError(name, fmt.Sprintf("Undefined variable '%s'.", name.Text))
	}
	
	return v, nil
}

func New() *Environment {
	return &Environment{
		variables: map[string]any{},
	}
}
