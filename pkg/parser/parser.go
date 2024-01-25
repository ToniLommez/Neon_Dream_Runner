package parser

import (
	"fmt"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Parser struct {
	tokens  []l.Token
	current int
}

func NewParser(tokens []l.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
}

func (p *Parser) Parse() Expr {
	expr, err := p.expression()
	if err != nil {
		fmt.Println(err)
		return nil
	}
	return expr
}

func (p Parser) String() string {
	s := ""
	for _, t := range p.tokens {
		s += fmt.Sprintf("%s ", t)
		if t.Type == l.NEW_LINE {
			s += "\n"
		}
	}
	s += "\n"
	return s
}

func (p *Parser) match(types ...l.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

func (p *Parser) check(t l.TokenType) bool {
	if p.isAtEnd() {
		return false
	}
	return p.peek().Type == t
}

func (p *Parser) advance() l.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == l.EOF
}

func (p *Parser) peek() l.Token {
	return p.tokens[p.current]
}

func (p *Parser) previous() l.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) expression() (Expr, error) {
	return p.equality()
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return expr, err
	}

	for p.match(l.NOT_EQUAL, l.EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		expr = Binary{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return expr, err
	}

	for p.match(l.GREATER, l.GREATER_EQUAL, l.LESS, l.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.term()
		expr = Binary{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return expr, err
	}

	for p.match(l.MINUS, l.PLUS) {
		operator := p.previous()
		right, err := p.factor()
		expr = Binary{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.unary()
	if err != nil {
		return expr, err
	}

	for p.match(l.SLASH, l.STAR) {
		operator := p.previous()
		right, err := p.unary()
		expr = Binary{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) unary() (Expr, error) {
	if p.match(l.BANG, l.MINUS) {
		operator := p.previous()
		right, err := p.unary()
		return Unary{operator, right}, err
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
	if p.match(l.FALSE) {
		return Literal{false}, nil
	}
	if p.match(l.TRUE) {
		return Literal{true}, nil
	}
	if p.match(l.NIL) {
		return Literal{nil}, nil
	}

	if p.match(l.NUMBER_LITERAL, l.FLOAT_LITERAL, l.STRING_LITERAL) {
		return Literal{p.previous().Literal}, nil
	}

	if p.match(l.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return expr, err
		}
		if _, err := p.consume(l.RIGHT_PAREN); err != nil {
			return nil, err
		}
		return Grouping{expr}, nil
	}

	token := p.tokens[p.current]
	return nil, errutils.Error(token.Line, token.Column, "expect expression.")
}

func (p *Parser) consume(expected l.TokenType) (l.Token, error) {
	if p.check(expected) {
		return p.advance(), nil
	}
	token := p.previous()
	return token, errutils.Error(token.Line, token.Column, fmt.Sprintf("expect %s after expression.", expected))
}

func (p *Parser) Synchronize() {
	p.advance()

	for !p.isAtEnd() {
		if p.previous().Type == l.SEMICOLON {
			return
		}

		switch p.peek().Type {
		case l.FN:
		case l.FN_BANG:
		case l.ASM:
		case l.LET:
		case l.LET_BANG:
		case l.LET_CHECK:
		case l.LET_BANG_CHECK:
		case l.OBJ:
		case l.FOR:
		case l.LOOP:
		case l.WHILE:
		case l.UNTIL:
		case l.DO:
		case l.PULSE:
		case l.CASE:
		case l.IF:
		case l.ELIF:
		case l.ELSE:
		case l.WHEN:
		case l.TRAIT:
		case l.PRINT:
		case l.PRINTF:
		case l.PRINTLN:
		case l.RETURN:
			return
		}

		p.advance()
	}
}
