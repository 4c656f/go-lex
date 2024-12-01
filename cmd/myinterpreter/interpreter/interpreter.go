package interpreter

import (
	"fmt"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/environment"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/stmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Interpreter struct {
	env  *environment.Environment
	out  any
	errs []error
}

func (i *Interpreter) Interp(program []stmt.Stmt) (any, []error) {
	for _, s := range program {
		if i.isErrorOcured() {
			break
		}
		i.exec(s)
	}
	return i.out, i.errs
}

func (i Interpreter) isErrorOcured() bool {
	return i.errs != nil
}

func (i *Interpreter) Eval(exp expression.Expression) (any, []error) {
	exp.Accept(i)
	return i.out, i.errs
}

func (i *Interpreter) exec(st stmt.Stmt) {
	st.Accept(i)
}

func (i *Interpreter) VisitExpressionStmt(s *stmt.ExpressionStmt) {
	i.Eval(s.Exp)
}

func (i *Interpreter) VisitPrintStmt(s *stmt.PrintStmt) {
	i.Eval(s.Exp)
	if i.isErrorOcured() {
		return
	}
	fmt.Println(i.String())
}

func (i *Interpreter) VisitVarStmt(s *stmt.VarStmt) {
	var value any
	if s.Init != nil {
		value, _ = i.Eval(s.Init)
	}
	i.env.Define(s.Name.Text, value)
}

func (i *Interpreter) VisitVarExpression(s *expression.VarExpression) {
	val, err := i.env.Get(s.Name)
	if err != nil {
		i.onError(err)
	}
	i.out = val
}

func (i *Interpreter) VisitAssignmentExpression(s *expression.AssignmentExpression) {
	v, errs := i.Eval(s.Val)
	if errs != nil {
		return
	}
	err := i.env.Assign(s.Name, v)
	if err != nil {
		i.onError(err)
	}
	i.out = v
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
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		if isNumeric {
			i.out = lNum - rNum
			return
		}
	case token.PLUS:
		if !isNumeric && !isString {
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be two numbers or two strings."))
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
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum / rNum
	case token.STAR:
		if !isNumeric {
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum * rNum
	case token.GREATER:
		if !isNumeric {
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum > rNum
	case token.GREATER_EQUAL:
		if !isNumeric {
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum >= rNum
	case token.LESS:
		if !isNumeric {
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum < rNum
	case token.LESS_EQUAL:
		if !isNumeric {
			i.onError(errors.NewRuntimeError(b.Op, "Operands must be numbers."))
		}
		i.out = lNum <= rNum
	case token.EQUAL_EQUAL:
		i.out = rhs == lhs
	case token.BANG_EQUAL:
		i.out = rhs != lhs
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
		default:
			i.onError(errors.NewRuntimeError(u.Op, "Operand must be a number."))
		}
		return
	case token.BANG:
		i.out = !isTrue(lhs)
		return
	}
	i.out = nil
}

func (i *Interpreter) onError(e error) {
	i.errs = append(i.errs, e)
	i.out = nil
}

func (i *Interpreter) VisitLiteral(u *expression.LiteralExpression) {
	i.out = u.Val.TokenValue.GetValue()
}

func New() *Interpreter {
	return &Interpreter{
		env: environment.New(),
	}
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
