package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Interpreter struct {
	out  any
	errs []error
}

func (i *Interpreter) Eval(e expression.Expression) (any, []error) {
	e.Accept(i)
	return i.out, i.errs
}

func (i Interpreter) String() string {
	if i.out == nil {
		return "nil"
	}
	return fmt.Sprintf("%v", i.out)
}

func (i *Interpreter) VisitBinary(b *expression.BinaryExpression) {
	lhs, _ := i.Eval(b.Lhs)
	rhs, _ := i.Eval(b.Rhs)

	lNum, rNum, isNumeric := matchOperandsType[float64](lhs, rhs)
	lStr, rStr, isString := matchOperandsType[string](lhs, rhs)
	switch b.Op.Type {
	case token.MINUS:
		if !isNumeric {
			i.onError(NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		if isNumeric {
			i.out = lNum - rNum
			return
		}
	case token.PLUS:
		if !isNumeric && !isString {
			i.onError(NewRuntimeError(b.Op, "Operands must be two numbers or two strings."))
		}
		if isNumeric {
			i.out = lNum + rNum
			return
		}
		if isString {
			i.out = lStr + rStr
			return
		}
	case token.SLASH:
		if !isNumeric {
			i.onError(NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum / rNum
	case token.STAR:
		if !isNumeric {
			i.onError(NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum * rNum
	case token.GREATER:
		if !isNumeric {
			i.onError(NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum > rNum
	case token.GREATER_EQUAL:
		if !isNumeric {
			i.onError(NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum >= rNum
	case token.LESS:
		if !isNumeric {
			i.onError(NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum < rNum
	case token.LESS_EQUAL:
		if !isNumeric {
			i.onError(NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum <= rNum
	}

}

func (i *Interpreter) VisitGrouping(g *expression.GroupingExpression) {
	i.Eval(g.Exp)
}

func (i *Interpreter) VisitUnary(u *expression.UnaryExpression) {
	lhs, _ := i.Eval(u.Rhs)

	switch u.Op.Type {
	case token.MINUS:
		switch v := lhs.(type) {
		case float64:
			i.out = -v
		case int:
			i.out = -v
		default:
			i.onError(NewRuntimeError(u.Op, "Operand must be a number."))
		}
		return
	case token.BANG:
		i.out = !isTrue(lhs)
		return
	}
	i.out = nil
}

func (i Interpreter) onError(e error) {

}

func (i *Interpreter) VisitLiteral(u *expression.LiteralExpression) {
	i.out = u.Val.TokenValue.GetValue()
}

func New() *Interpreter {
	return &Interpreter{}
}

func matchOperandsType[V int | float64 | string](lhs any, rhs any) (V, V, bool) {
	lhv, lOk := lhs.(V)
	rhv, rOk := rhs.(V)
	return lhv, rhv, lOk && rOk
}

func isTrue(v any) bool {
	switch t := v.(type) {
	case bool:
		return t
	case nil:
		return false
	default:
		return true
	}
}
