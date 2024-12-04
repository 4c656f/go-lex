package interpreter

import (
	"fmt"
	"time"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/environment"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/errors"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/stmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Interpreter struct {
	env           *environment.Environment
	globals       *environment.Environment
	out           any
	returnCalls   int
	functionCalls int
	errs          []error
}

type Callable interface {
	Call(interp *Interpreter, args []any) any
	Arity() int
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

func (i *Interpreter) VisitIfStmt(s *stmt.IfStmt) {
	cond, _ := i.Eval(s.Condition)
	if isTrue(cond) {
		i.exec(s.ThenBranch)
		return
	}
	if s.ElseBranch != nil {
		i.exec(s.ElseBranch)
	}
}

func (i *Interpreter) VisitWhileStmt(s *stmt.WhileStmt) {
	for {
		v, _ := i.Eval(s.Condition)
		if !isTrue(v) {
			break
		}
		i.exec(s.Body)
		if i.isErrorOcured() || i.isReturnCallOccured() {
			break
		}
	}
}

func (i Interpreter) isReturnCallOccured() bool {
	return i.returnCalls > 0
}

func (i Interpreter) isFunctionCallOccured() bool {
	return i.functionCalls > 0
}

func (i *Interpreter) VisitFunctionDeclarationStmt(s *stmt.FunctionDeclarationStmt) {
	i.env.Define(s.Name.Text, NewFunction(s, i.env))
}

func (i *Interpreter) VisitReturnStmt(s *stmt.ReturnStmt) {
	if !i.isFunctionCallOccured() {
		i.onError(NewRuntimeError(s.Keywoard, "return is not allowed outside of a function body"))
		i.out = nil
		return
	}

	var value any
	if s.Exp != nil {
		v, _ := i.Eval(s.Exp)
		if i.isErrorOcured() {
			return
		}
		value = v
	}
	i.out = value
	i.returnCalls++
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

func (i *Interpreter) VisitBlockStmt(s *stmt.BlockStmt) {
	i.executeBlock(s.Statements, environment.New(i.env))
}

func (i *Interpreter) executeBlock(stmts []stmt.Stmt, env *environment.Environment) {
	prevEnv := i.env
	i.env = env

	for _, s := range stmts {
		i.exec(s)
		if i.isErrorOcured() || i.isReturnCallOccured() {
			i.env = prevEnv
			break
		}
	}
	if !i.isReturnCallOccured(){
		i.out = nil
	}
	i.env = prevEnv
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

func (i *Interpreter) VisitLogicalExpression(s *expression.LogicalExpression) {
	left, _ := i.Eval(s.Lhs)

	if s.Op.Type == token.OR {
		if isTrue(left) {
			i.out = left
			return
		}
	}
	if s.Op.Type == token.AND {
		if !isTrue(left) {
			i.out = left
			return
		}
	}

	i.Eval(s.Rhs)
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

func (i *Interpreter) VisitFunctionCallExpression(g *expression.FunctionCallExpression) {
	calle, _ := i.Eval(g.Callee)
	argsValues := make([]any, len(g.Args))
	for idx, a := range g.Args {
		argV, err := i.Eval(a)
		if err != nil {

		}
		argsValues[idx] = argV
	}
	function, ok := calle.(Callable)
	if !ok {
		i.onError(NewRuntimeError(g.RightParan, "Can only call functions and classes."))
		return
	}
	if function.Arity() != len(argsValues) {
		i.onError(NewRuntimeError(g.RightParan, fmt.Sprintf("Expected %v arguments but got %v.", function.Arity(), len(argsValues))))
		return
	}
	i.out = function.Call(i, argsValues)
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
	globalEnv := environment.New(nil)
	defineGlobals(globalEnv)
	return &Interpreter{
		env:     globalEnv,
		globals: globalEnv,
	}
}

func defineGlobals(env *environment.Environment) {
	env.Define("clock", NewClockFc())
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

type NativeClock struct {
}

func (c *NativeClock) Call(interp *Interpreter, args []any) any {
	return float64(time.Now().Unix())
}

func (c NativeClock) Arity() int {
	return 0
}

func (c NativeClock) String() string {
	return "<native fn>"
}

func NewClockFc() *NativeClock {
	return &NativeClock{}
}

type Function struct {
	closure     *environment.Environment
	declaration *stmt.FunctionDeclarationStmt
}

func (c *Function) Call(interp *Interpreter, args []any) any {
	env := environment.New(c.closure)
	interp.functionCalls += 1
	startReturnCalls := interp.returnCalls
	for i := 0; i < len(c.declaration.Args); i++ {
		env.Define(c.declaration.Args[i].Text, args[i])
	}
	interp.executeBlock(c.declaration.Body, env)
	interp.functionCalls -= 1
	// decrement return calls only if return was called inside function
	if interp.returnCalls > startReturnCalls {
		interp.returnCalls -= 1
	}
	return interp.out
}

func (c Function) Arity() int {
	return len(c.declaration.Args)
}

func (c Function) String() string {
	return fmt.Sprintf("<fn %s>", c.declaration.Name.Text)
}

func NewFunction(declaration *stmt.FunctionDeclarationStmt, closure *environment.Environment) *Function {
	return &Function{
		closure:     closure,
		declaration: declaration,
	}
}
