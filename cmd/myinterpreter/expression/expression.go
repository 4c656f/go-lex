package expression

import "github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"

type Visitor interface {
	VisitBinary(b *BinaryExpression)
	VisitGrouping(g *GroupingExpression)
	VisitUnary(u *UnaryExpression)
	VisitLiteral(u *LiteralExpression)
	VisitVarExpression(u *VarExpression)
	VisitAssignmentExpression(u *AssignmentExpression)
	VisitOrExpression(u *OrExpression)
	VisitAndExpression(u *AndExpression)
}

type Expression interface {
	Accept(v Visitor)
}

type BinaryExpression struct {
	Lhs Expression
	Op  *token.Token
	Rhs Expression
}

type GroupingExpression struct {
	Exp Expression
}

type LiteralExpression struct {
	Val *token.Token
}

type UnaryExpression struct {
	Op  *token.Token
	Rhs Expression
}

type VarExpression struct {
	Name *token.Token
}

type AssignmentExpression struct {
	Name *token.Token
	Val  Expression
}

type OrExpression struct {
	Lhs Expression
	Rhs Expression
}

type AndExpression struct {
	Lhs Expression
	Rhs Expression
}

func (this *BinaryExpression) Accept(v Visitor) {
	v.VisitBinary(this)
}

func (this *AssignmentExpression) Accept(v Visitor) {
	v.VisitAssignmentExpression(this)
}

func (this *GroupingExpression) Accept(v Visitor) {
	v.VisitGrouping(this)
}

func (this *LiteralExpression) Accept(v Visitor) {
	v.VisitLiteral(this)
}

func (this *UnaryExpression) Accept(v Visitor) {
	v.VisitUnary(this)
}

func (this *VarExpression) Accept(v Visitor) {
	v.VisitVarExpression(this)
}

func (this *OrExpression) Accept(v Visitor) {
	v.VisitOrExpression(this)
}

func (this *AndExpression) Accept(v Visitor) {
	v.VisitAndExpression(this)
}

func NewBinaryExpression(
	lhs Expression,
	op *token.Token,
	rhs Expression,
) *BinaryExpression {
	return &BinaryExpression{
		Lhs: lhs,
		Op:  op,
		Rhs: rhs,
	}
}

func NewGroupingExpression(exp Expression) *GroupingExpression {
	return &GroupingExpression{exp}
}

func NewLiteralExpression(val *token.Token) *LiteralExpression {
	return &LiteralExpression{val}
}

func NewUnaryExpression(
	op *token.Token,
	exp Expression,
) *UnaryExpression {
	return &UnaryExpression{
		Op:  op,
		Rhs: exp,
	}
}

func NewVarExpression(
	name *token.Token,
) *VarExpression {
	return &VarExpression{
		Name: name,
	}
}

func NewAssignmentExprExpression(
	name *token.Token,
	value Expression,
) *AssignmentExpression {
	return &AssignmentExpression{
		Name: name,
		Val:  value,
	}
}

func NewOrExpression(
	lhs Expression,
	rhs Expression,
) *OrExpression {
	return &OrExpression{
		Lhs: lhs,
		Rhs: rhs,
	}
}

func NewAndExpression(
	lhs Expression,
	rhs Expression,
) *AndExpression {
	return &AndExpression{
		Lhs: lhs,
		Rhs: rhs,
	}
}
