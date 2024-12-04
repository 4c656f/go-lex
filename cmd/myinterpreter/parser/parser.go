package parser

import (
	"fmt"
	"strings"

	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/expression"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/stmt"
	"github.com/codecrafters-io/interpreter-starter-go/cmd/myinterpreter/token"
)

type ASTPrinter struct {
	outString string
}

func (a *ASTPrinter) PrintProgram(st []stmt.Stmt) string {
	var b strings.Builder
	if st == nil {
		return ""
	}

	for _, s := range st {
		s.Accept(a)
		b.WriteString(a.Out())
	}
	return b.String()
}

func (a *ASTPrinter) Print(exp expression.Expression) string {
	if exp == nil {
		return ""
	}
	exp.Accept(a)
	return a.Out()
}

func (a *ASTPrinter) Out() string {
	return a.outString
}

func (a *ASTPrinter) VisitExpressionStmt(s *stmt.ExpressionStmt) {
	a.outString = a.parenthesize("stmt", s.Exp)
}

func (a *ASTPrinter) VisitBlockStmt(s *stmt.BlockStmt) {
	var inside strings.Builder
	for _, st := range s.Statements {
		st.Accept(a)
		inside.WriteString(a.outString)
	}
	a.outString = fmt.Sprintf("{ %s }", inside.String())
}

func (a *ASTPrinter) VisitPrintStmt(s *stmt.PrintStmt) {
	a.outString = a.parenthesize("print", s.Exp)
}

func (a *ASTPrinter) VisitReturnStmt(s *stmt.ReturnStmt) {
	a.outString = a.parenthesize("return", s.Exp)
}

func (a *ASTPrinter) VisitIfStmt(s *stmt.IfStmt) {
	s.ThenBranch.Accept(a)
	a.outString = fmt.Sprintf("%s, {\n%s\n}", a.parenthesize("if", s.Condition), a.Out())
}

func (a *ASTPrinter) VisitVarStmt(s *stmt.VarStmt) {
	a.outString = a.parenthesize("var = ", s.Init)
}

func (a *ASTPrinter) VisitBinary(b *expression.BinaryExpression) {
	a.outString = a.parenthesize(b.Op.Text, b.Lhs, b.Rhs)
}

func (a *ASTPrinter) VisitGrouping(g *expression.GroupingExpression) {
	a.outString = a.parenthesize("group", g.Exp)
}

func (a *ASTPrinter) VisitLiteral(l *expression.LiteralExpression) {
	if l.Val.TokenValue.Type == token.BoolValue || l.Val.TokenValue.Type == token.NullValue {
		a.outString = l.Val.Text
		return
	}
	a.outString = l.Val.TokenValue.String()
}

func (a *ASTPrinter) VisitUnary(u *expression.UnaryExpression) {
	a.outString = a.parenthesize(u.Op.Text, u.Rhs)
}

func (a *ASTPrinter) VisitVarExpression(u *expression.VarExpression) {
	a.outString = fmt.Sprintf("var %s", u.Name.Text)
}

func (a *ASTPrinter) VisitAssignmentExpression(u *expression.AssignmentExpression) {
	a.outString = fmt.Sprintf("ass %s", a.parenthesize(u.Name.Text, u.Val))
}

func (a *ASTPrinter) VisitLogicalExpression(u *expression.LogicalExpression) {
	a.outString = a.parenthesize(u.Op.Text, u.Lhs, u.Rhs)
}

func (a *ASTPrinter) VisitWhileStmt(s *stmt.WhileStmt) {
	s.Body.Accept(a)
	a.outString = fmt.Sprintf("%s, {\n%s\n}", a.parenthesize("while", s.Condition), a.Out())
}

func (a *ASTPrinter) VisitFunctionCallExpression(f *expression.FunctionCallExpression) {
	var args []expression.Expression
	args = append(args, f.Callee)
	args = append(args, f.Args...)
	a.outString = a.parenthesize("call", args...)
}

func (a *ASTPrinter) VisitFunctionDeclarationStmt(f *stmt.FunctionDeclarationStmt) {
	a.VisitBlockStmt(stmt.NewBlockStmt(f.Body))
	a.outString = fmt.Sprintf("fun %s () %s", f.Name.Text, a.Out())
}

// Helper function to create parenthesized expressions
func (a *ASTPrinter) parenthesize(name string, exprs ...expression.Expression) string {
	var result strings.Builder
	result.WriteString("(" + name)
	for _, expr := range exprs {
		result.WriteString(" ")
		expr.Accept(a)
		result.WriteString(a.Out())
	}
	result.WriteString(")")
	return result.String()
}

func NewAstPrinter() *ASTPrinter {
	return &ASTPrinter{}
}

type Parser struct {
	tokens []*token.Token
	errors []error
	cur    int
}

func (p *Parser) Parse() (expression.Expression, []error) {
	return p.expression(), p.errors
}

func (p *Parser) ParseProgram() ([]stmt.Stmt, []error) {
	statements := []stmt.Stmt{}
	for !p.isAtEnd() {
		statements = append(statements, p.declaration())
	}
	return statements, p.errors
}

func (p *Parser) declaration() stmt.Stmt {
	if p.match(token.FUN) {
		return p.functionDeclaration()
	}
	if p.match(token.VAR) {
		return p.varDeclaration()
	}
	return p.statement()
}

func (p *Parser) statement() stmt.Stmt {
	if p.match(token.IF) {
		return p.ifStmt()
	}
	if p.match(token.WHILE) {
		return p.whileStmt()
	}
	if p.match(token.FOR) {
		return p.forStmt()
	}
	if p.match(token.PRINT) {
		return p.printStmt()
	}
	if p.match(token.LEFT_BRACE) {
		return stmt.NewBlockStmt(p.blockStmt())
	}
	if p.match(token.RETURN) {
		return p.returnStmt()
	}
	return p.expStmt()
}

func (p *Parser) returnStmt() stmt.Stmt {
	keywoard := p.prev()
	var exp expression.Expression
	if !p.check(token.SEMICOLON) {
		exp = p.expression()
	}
	_, err := p.consume(token.SEMICOLON, "Expect ';' after return value.")
	if err != nil {
		return nil
	}
	return stmt.NewReturnStmt(keywoard, exp)
}

func (p *Parser) whileStmt() stmt.Stmt {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'while'.")
	if err != nil {
		return nil
	}
	condition := p.expression()
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after condition.")
	if err != nil {
		return nil
	}
	body := p.statement()
	return stmt.NewWhileStmt(condition, body)
}

func (p *Parser) forStmt() stmt.Stmt {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'for'.")
	if err != nil {
		return nil
	}
	var initializer stmt.Stmt
	var condition expression.Expression
	var incriment expression.Expression
	if p.match(token.SEMICOLON) {

	} else if p.match(token.VAR) {
		initializer = p.varDeclaration()
	} else {
		initializer = p.expStmt()
	}

	if !p.check(token.SEMICOLON) {
		condition = p.expression()
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after loop condition.")
	if err != nil {
		return nil
	}

	if !p.check(token.RIGHT_PAREN) {
		incriment = p.expression()
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after for clauses.")
	if err != nil {
		return nil
	}
	body := p.statement()

	if incriment != nil {
		body = stmt.NewBlockStmt([]stmt.Stmt{
			body,
			stmt.NewExpressionStmt(incriment),
		})
	}

	if condition == nil {
		condition = expression.NewLiteralExpression(token.NewToken(token.BoolValue, 1, "true", token.NewBoolValue(true)))
	}
	body = stmt.NewWhileStmt(condition, body)
	if initializer != nil {
		body = stmt.NewBlockStmt([]stmt.Stmt{
			initializer,
			body,
		})
	}
	return body
}

func (p *Parser) ifStmt() stmt.Stmt {
	_, err := p.consume(token.LEFT_PAREN, "Expect '(' after 'if'.")
	if err != nil {
		return nil
	}
	condition := p.expression()
	_, err = p.consume(token.RIGHT_PAREN, "Expect ')' after if condition.")
	if err != nil {
		return nil
	}
	flow := p.statement()
	var elseStmt stmt.Stmt
	if p.match(token.ELSE) {
		elseStmt = p.statement()
	}
	return stmt.NewIfStmt(condition, flow, elseStmt)
}

func (p *Parser) functionDeclaration() stmt.Stmt {
	name, err := p.consume(token.IDENTIFIER, "Expect function name.")
	if err != nil {
		return nil
	}
	_, err = p.consume(token.LEFT_PAREN, "Expect ( after function name.")
	if err != nil {
		return nil
	}
	args := []*token.Token{}
	for !p.check(token.RIGHT_PAREN) {
		arg, err := p.consume(token.IDENTIFIER, "Expect parameter name.")
		if err != nil {
			return nil
		}
		if len(args) >= 255 {
			p.onError(NewParserError(arg, "Can't have more than 255 parameters."))
			break
		}
		args = append(args, arg)
		if !p.match(token.COMMA) {
			break
		}
	}
	_, err = p.consume(token.RIGHT_PAREN, "Expect ) after function name.")
	if err != nil {
		return nil
	}
	_, err = p.consume(token.LEFT_BRACE, "Expect { after function body.")
	if err != nil {
		return nil
	}
	body := p.blockStmt()
	return stmt.NewFunctionDeclarationStmt(name, body, args)
}

func (p *Parser) varDeclaration() stmt.Stmt {
	name, err := p.consume(token.IDENTIFIER, "Expect variable name.")
	if err != nil {
		return nil
	}
	var initializer expression.Expression

	if p.match(token.EQUAL) {
		initializer = p.expression()
	}
	_, err = p.consume(token.SEMICOLON, "Expect ';' after variable declaration.")
	if err != nil {
		return nil
	}
	return stmt.NewVarStmt(name, initializer)
}

func (p *Parser) printStmt() stmt.Stmt {
	value := p.expression()
	_, err := p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil
	}
	return stmt.NewPrintStmt(value)
}

func (p *Parser) blockStmt() []stmt.Stmt {
	statemnts := []stmt.Stmt{}

	for !p.check(token.RIGHT_BRACE) && !p.isAtEnd() {
		statemnts = append(statemnts, p.declaration())
	}
	_, err := p.consume(token.RIGHT_BRACE, "Expect '}' after block.")
	if err != nil {
		return nil
	}
	return statemnts
}

func (p *Parser) expStmt() stmt.Stmt {
	exp := p.expression()
	_, err := p.consume(token.SEMICOLON, "Expect ';' after value.")
	if err != nil {
		return nil
	}
	return stmt.NewExpressionStmt(exp)
}

func (p *Parser) expression() expression.Expression {
	return p.assignment()
}

func (p *Parser) assignment() expression.Expression {
	exp := p.logicalOr()
	if p.match(token.EQUAL) {
		equals := p.prev()
		value := p.assignment()
		name, ok := exp.(*expression.VarExpression)
		if !ok {
			p.onError(NewParserError(equals, "Invalid assignment target."))
			return exp
		}
		return expression.NewAssignmentExprExpression(name.Name, value)
	}
	return exp
}

func (p *Parser) logicalOr() expression.Expression {
	exp := p.logicalAnd()

	for p.match(token.OR) {
		op := p.prev()
		rhs := p.logicalAnd()
		exp = expression.NewLogicalExpression(exp, op, rhs)
	}
	return exp
}

func (p *Parser) logicalAnd() expression.Expression {
	exp := p.equality()

	for p.match(token.AND) {
		op := p.prev()
		rhs := p.equality()
		exp = expression.NewLogicalExpression(exp, op, rhs)
	}

	return exp
}

func (p *Parser) equality() expression.Expression {
	exp := p.comparison()

	for p.match(token.BANG_EQUAL, token.EQUAL_EQUAL) {
		op := p.prev()
		rhs := p.comparison()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}

	return exp
}

func (p *Parser) comparison() expression.Expression {
	exp := p.term()

	for p.match(token.GREATER, token.GREATER_EQUAL, token.LESS_EQUAL, token.LESS) {
		op := p.prev()
		rhs := p.term()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}

	return exp
}

func (p *Parser) term() expression.Expression {
	exp := p.factor()
	for p.match(token.MINUS, token.PLUS) {
		op := p.prev()
		rhs := p.factor()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}
	return exp
}

func (p *Parser) factor() expression.Expression {
	exp := p.unary()
	for p.match(token.SLASH, token.STAR) {
		op := p.prev()
		rhs := p.unary()
		exp = expression.NewBinaryExpression(exp, op, rhs)
	}
	return exp
}

func (p *Parser) unary() expression.Expression {
	if p.match(token.BANG, token.MINUS) {
		op := p.prev()
		rhs := p.unary()
		return expression.NewUnaryExpression(op, rhs)
	}
	return p.call()
}

func (p *Parser) call() expression.Expression {
	callee := p.primary()
	for p.match(token.LEFT_PAREN) {
		args := []expression.Expression{}
		if !p.check(token.RIGHT_PAREN) {
			for {
				args = append(args, p.expression())
				if !p.match(token.COMMA) {
					break
				}
			}
		}
		if len(args) > 255 {
			p.onError(NewParserError(p.peek(), "Can't have more than 255 arguments."))
		}
		rightParan, err := p.consume(token.RIGHT_PAREN, "Expect ')' after arguments.")
		if err != nil {
			return nil
		}
		callee = expression.NewFunctionCallExpression(callee, args, rightParan)
	}
	return callee
}

func (p *Parser) primary() expression.Expression {
	if p.match(token.TRUE, token.FALSE, token.NIL, token.NUMBER, token.STRING) {
		return expression.NewLiteralExpression(p.prev())
	}
	if p.match(token.IDENTIFIER) {
		return expression.NewVarExpression(p.prev())
	}

	if p.match(token.LEFT_PAREN) {
		exp := p.expression()
		_, err := p.consume(token.RIGHT_PAREN, "Expect ')' after expression.")
		if err != nil {
			return nil
		}
		return expression.NewGroupingExpression(exp)
	}
	p.onError(NewParserError(p.peek(), "Expect expression."))
	return nil
}

func (p *Parser) consume(t token.TokenType, message string) (*token.Token, error) {
	if p.check(t) {
		return p.advance(), nil
	}
	err := NewParserError(p.peek(), message)
	p.onError(err)
	return nil, err
}

func (p *Parser) onError(err error) {
	p.errors = append(p.errors, err)
	p.sync()
}

func (p *Parser) match(tokens ...token.TokenType) bool {
	for _, t := range tokens {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) sync() {
	p.advance()

	for !p.isAtEnd() {
		if p.prev().Type == token.SEMICOLON {
			return
		}
		switch p.peek().Type {
		case token.CLASS:
			return
		case token.FUN:
			return
		case token.VAR:
			return
		case token.FOR:
			return
		case token.IF:
			return
		case token.WHILE:
			return
		case token.PRINT:
			return
		case token.RETURN:
			return
		}
		p.advance()
	}
}

func (p *Parser) check(t token.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() *token.Token {
	if !p.isAtEnd() {
		p.cur++
	}
	return p.prev()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == token.EOF
}

func (p *Parser) peek() *token.Token {
	return p.tokens[p.cur]
}

func (p *Parser) prev() *token.Token {
	if p.cur == 0 {
		return nil
	}
	return p.tokens[p.cur-1]
}

func New(tokens []*token.Token) *Parser {
	return &Parser{
		tokens: tokens,
	}
}
