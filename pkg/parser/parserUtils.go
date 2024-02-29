package parser

import (
	"fmt"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

func NewParser(tokens []l.Token) Parser {
	return Parser{
		tokens:  tokens,
		current: 0,
	}
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
	x := p.peek()
	return x.Type == t
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

func (p *Parser) peekN(n int) (bool, l.Token) {
	if p.current+n < len(p.tokens) {
		return true, p.tokens[p.current+n]
	} else {
		return false, p.tokens[p.current]
	}
}

func (p *Parser) previous() l.Token {
	return p.tokens[p.current-1]
}

func (p *Parser) consume(expected l.TokenType) (l.Token, error) {
	if p.check(expected) {
		return p.advance(), nil
	}
	token := p.previous()
	return token, errutils.Error(token.Line, token.Column, token.Lexeme, errutils.PARSER, fmt.Sprintf("expect %s after expression.", expected))
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
