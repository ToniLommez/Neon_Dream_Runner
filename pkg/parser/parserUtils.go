package parser

import (
	"fmt"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

func NewParser(tokens []l.Token) Parser {
	return Parser{
		Tokens:  tokens,
		Current: 0,
		Depth:   0,
	}
}

func (p Parser) String() string {
	s := ""
	for _, t := range p.Tokens {
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
	if !p.isLastToken() && !p.isAtEnd() {
		p.Current++
	}
	return p.previous()
}

func (p *Parser) hasMore() bool {
	return !p.isLastToken() && !p.isAtEnd()
}

func (p *Parser) isAtEnd() bool {
	return p.peek().Type == l.EOF
}

func (p *Parser) isLastToken() bool {
	return len(p.Tokens) <= p.Current+1
}

func (p *Parser) peek() l.Token {
	return p.Tokens[p.Current]
}

func (p *Parser) peekN(n int) (bool, l.Token) {
	if p.Current+n < len(p.Tokens) {
		return true, p.Tokens[p.Current+n]
	} else {
		return false, p.Tokens[p.Current]
	}
}

func (p *Parser) previous() l.Token {
	return p.Tokens[p.Current-1]
}

func (p *Parser) consume(expected l.TokenType) (l.Token, error) {
	if p.check(expected) {
		return p.advance(), nil
	}
	token := p.peek()
	return token, e.Error(token.Line, token.Column, token.Lexeme, e.PARSER, fmt.Sprintf("expect %s", expected))
}

// ensureNotUnterminated checks if the parser found an end and should have more after the new_line
func (p *Parser) ensureNotUnterminated() error {
	for p.peek().Type == l.NEW_LINE {
		p.consume(l.NEW_LINE)
		if !p.hasMore() {
			return e.Error(0, 0, "", e.UNTERMINATED_STATEMENT, "")
		}
	}
	return nil
}

func (p *Parser) Synchronize() {
	p.advance()

	for !p.isLastToken() && !p.isAtEnd() {
		if p.previous().Type == l.NEW_LINE {
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
		case l.PUT:
		case l.PRINT:
		case l.PRINTF:
		case l.PRINTLN:
		case l.RETURN:
			return
		}

		p.advance()
	}
}
