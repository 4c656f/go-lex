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
	VisitBlockStmt(s *BlockStmt)
	VisitIfStmt(s *IfStmt)
	VisitWhileStmt(s *WhileStmt)
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

type BlockStmt struct {
	Statements []Stmt
}

type IfStmt struct {
	Condition  expression.Expression
	ThenBranch Stmt
	ElseBranch Stmt
}

type WhileStmt struct {
	Condition expression.Expression
	Body      Stmt
}

func (s *WhileStmt) Accept(v Visitor) {
	v.VisitWhileStmt(s)
}

func (s *BlockStmt) Accept(v Visitor) {
	v.VisitBlockStmt(s)
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

func (s *IfStmt) Accept(v Visitor) {
	v.VisitIfStmt(s)
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

func NewBlockStmt(statements []Stmt) *BlockStmt {
	return &BlockStmt{
		Statements: statements,
	}
}

func NewIfStmt(Condition expression.Expression,
	ThenBranch Stmt,
	ElseBranch Stmt) *IfStmt {
	return &IfStmt{
		Condition:  Condition,
		ThenBranch: ThenBranch,
		ElseBranch: ElseBranch,
	}
}

func NewWhileStmt(condition expression.Expression, body Stmt) *WhileStmt {
	return &WhileStmt{
		Condition: condition,
		Body:      body,
	}
}
