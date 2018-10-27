package parser

import (
	"fmt"

	"github.com/Fr0stb1t3/go-vtwin/lexer"
	"github.com/Fr0stb1t3/go-vtwin/token"
)

type Parser struct {
	l         *lexer.Lexer
	errors    []string
	topScope  Scope
	curToken  token.Token
	peekToken token.Token
}

func New(l *lexer.Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}

	p.topScope = Scope{Objects: make(map[string]*LetStatement)}
	// Read two tokens, so curToken and peekToken are both set
	p.nextToken()
	p.nextToken()

	return p
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) parseUnaryExpr() Expression {
	operator := token.NewToken(token.ADD, '+')
	switch p.curToken.Type {
	case token.ADD, token.SUBT, token.NOT, token.XOR, token.AND:
		operator = p.curToken
		p.nextToken()
	}
	return UnaryExpression{
		Operator: operator,
		Operand:  p.curToken,
	}
}
func (p *Parser) parseParenExpr() Expression {
	var parenCounter int
	pExpr := ParenExpression{}
	for p.tokenIs(token.LPAREN) {
		parenCounter++
		pExpr.Lparen = p.curToken
		p.nextToken()
	}
	pExpr.Expr = p.parseExpression(token.RPAREN)
	for p.peekTokenIs(token.RPAREN) {
		parenCounter--
		p.nextToken()
	}
	pExpr.Rparen = p.curToken
	parenCounter--
	if parenCounter > 0 {
		panic("Paren mismatch")
	}

	return pExpr
}
func (p *Parser) parseBinaryExpr(endToken token.Type, rhs bool, prec int) Expression {
	expression := BinaryExpression{}

	for !p.tokenIs(endToken) {
		if expression.completeNode() {
			if rhs {
				return expression
			} else if !p.tokenIs(endToken) {
				expression = expression.shiftNode()
			}
		}
		switch {
		case p.tokenIs(token.LPAREN):
			expr := p.parseParenExpr()
			expression.addSubnode(expr)

		case p.curToken.Type.IsOpertor() && !expression.Operator.Type.IsOpertor():
			expression.Operator = p.curToken

		case expression.emptyNode():
			expr := p.parseUnaryExpr()
			expression.addSubnode(expr)

		case rhs && expression.Operator.Type.IsOpertor() && p.peekPrecedence() > prec:
			expr := p.parseUnaryExpr()
			expression.addSubnode(expr)
			expression = expression.shiftNode()

		case expression.Operator.Type.IsOpertor() && p.peekPrecedence() > expression.Operator.Type.Precedence():
			expr := p.parseBinaryExpr(endToken, true, expression.Operator.Type.Precedence())
			expression.addSubnode(expr)

		default:
			expr := p.parseUnaryExpr()
			expression.addSubnode(expr)
		}

		if !p.peekTokenIs(token.EOF) && !(rhs && expression.completeNode()) {
			p.nextToken()
		}
	}
	return expression
}

func (p *Parser) parseExpression(endToken token.Type) Expression {
	switch {
	case p.tokenIs(token.LPAREN), p.peekToken.Type.IsOpertor():
		return p.parseBinaryExpr(endToken, false, 0)
	default:
		return p.parseUnaryExpr()
	}
}

func (p *Parser) parseLetStatement() LetStatement {
	assignment := LetStatement{
		Token: p.curToken,
	}
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		assignment.Name = &Identifier{
			Token: p.curToken,
			Value: p.curToken.Literal,
		}
	}
	if !p.peekTokenIs(token.ASSIGN) {
		panic("Invalid Let assignment")
	}
	p.nextToken() // TODO
	p.nextToken()

	assignment.Expr = p.parseExpression(token.SEMICOLON)
	return assignment
}
func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case token.CONST:
		fmt.Printf("parse as immutable assignment %v\n", p.curToken.Literal)
	case token.LET:
		stmt := p.parseLetStatement()
		p.topScope.Objects[stmt.Name.Value] = &stmt
		return stmt
	case token.RETURN:
		fmt.Printf("parse as return statement %v\n", p.curToken.Literal)
	case token.LPAREN, token.INT:
		start := p.curToken
		expression := p.parseExpression(token.SEMICOLON)
		return ExpressionStatement{
			Token: start,
			Expr:  expression,
		}
	}
	return nil
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}
	return program
}
func (p *Parser) peekPrecedence() int {
	return p.peekToken.Type.Precedence()
}
func (p *Parser) tokenIs(t token.Type) bool {
	return p.curToken.Type == t
}
func (p *Parser) peekTokenIs(t token.Type) bool {
	return p.peekToken.Type == t
}
