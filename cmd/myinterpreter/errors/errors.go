package errors

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type RuntimeError struct {
	op  *token.Token
	msg string
}

func NewRuntimeError(op *token.Token, msg string) *RuntimeError {
	return &RuntimeError{
		op:  op,
		msg: msg,
	}
}

func (e RuntimeError) Error() string {
	return fmt.Sprintf("%s\n[line %v]", e.msg, e.op.Line)
}
