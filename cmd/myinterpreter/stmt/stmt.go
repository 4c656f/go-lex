package stmt

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"

type Stmt interface {
	Accept(v Visitor)
}

type Visitor interface {
	VisitExpressionStmt(s *ExpressionStmt)
	VisitPrintStmt(s *PrintStmt)
}

type ExpressionStmt struct {
	Exp expression.Expression
}

type PrintStmt struct {
	Exp expression.Expression
}

func (s *ExpressionStmt) Accept(v Visitor) {
	v.VisitExpressionStmt(s)
}

func (s *PrintStmt) Accept(v Visitor) {
	v.VisitPrintStmt(s)
}

func NewExpressionStmt(exp expression.Expression) *ExpressionStmt {
	return &ExpressionStmt{
		Exp: exp,
	}
}

func NewPrintStmt(exp expression.Expression) *PrintStmt {
	return &PrintStmt{
		Exp: exp,
	}
}
