package stmt

import (
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type Stmt interface {
	Accept(v Visitor)
}

type Visitor interface {
	VisitExpressionStmt(s *ExpressionStmt)
	VisitPrintStmt(s *PrintStmt)
	VisitVarStmt(s *VarStmt)
}

type ExpressionStmt struct {
	Exp expression.Expression
}

type PrintStmt struct {
	Exp expression.Expression
}

type VarStmt struct {
	Name *token.Token
	Init expression.Expression
}

func (s *ExpressionStmt) Accept(v Visitor) {
	v.VisitExpressionStmt(s)
}

func (s *PrintStmt) Accept(v Visitor) {
	v.VisitPrintStmt(s)
}

func (s *VarStmt) Accept(v Visitor) {
	v.VisitVarStmt(s)
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

func NewVarStmt(name *token.Token, init expression.Expression) *VarStmt {
	return &VarStmt{
		Name: name,
		Init: init,
	}
}
