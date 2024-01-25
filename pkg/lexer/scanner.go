package lexer

import (
	"fmt"
	"strconv"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
)

type Scanner struct {
	source  string
	tokens  []Token
	start   int
	current int
	line    int
	column  int
}

func NewScanner(source string) Scanner {
	return Scanner{
		source:  source,
		tokens:  []Token{},
		start:   0,
		current: 0,
		line:    1,
		column:  1,
	}
}

func (s *Scanner) ScanTokens() ([]Token, error) {
	for !s.isAtEnd() {
		s.start = s.current
		if err := s.scanToken(); err != nil {
			return s.tokens, err
		}
	}

	// Adicionar token de EOF ao final
	s.tokens = append(s.tokens, Token{Type: EOF, Line: s.line})
	return s.tokens, nil
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
}

func (s *Scanner) advance() byte {
	s.current++
	s.column++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.line, Column: s.column})
}

func (s Scanner) peek() byte {
	if s.isAtEnd() {
		return byte('\x00')
	}
	return s.source[s.current]
}

func (s Scanner) peekN(n int) string {
	if s.isAtEnd() {
		return "\x00"
	}

	end := s.current + n
	if end > len(s.source) {
		end = len(s.source)
	}

	return s.source[s.current:end]
}

func (s Scanner) peekNext() byte {
	if s.current+1 >= len(s.source) {
		return '\x00'
	}
	return s.source[s.current+1]
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}

	s.current++
	s.column++
	return true
}

func (s *Scanner) matchN(expected string) bool {
	if s.isAtEnd() {
		return false
	}

	end := s.current + len(expected)
	if end > len(s.source) || s.source[s.current:end] != expected {
		return false
	}

	s.current += len(expected)
	s.column += len(expected)
	return true
}

func (s *Scanner) switchToken(tokens []string, tokenTypes []TokenType) {
	for i, t := range tokens {
		if s.matchN(t) {
			s.addToken(tokenTypes[i], nil)
			return
		}
	}

	s.addToken(tokenTypes[len(tokenTypes)-1], nil)
}

func (s *Scanner) string() error {
	for s.peek() != '"' && !s.isAtEnd() {
		if s.peek() == '\n' {
			return errutils.Error(s.line, s.column, "unterminated string.")
		}
		s.advance()
	}

	// the closing ".
	s.advance()

	// trim quotes
	str := s.source[s.start+1 : s.current-1]
	s.addToken(STRING_LITERAL, str)

	return nil
}

func (s *Scanner) number() {
	isFloat := false

	for isDigit(s.peek()) {
		s.advance()
	}

	// Look for a fractional part.
	if s.peek() == '.' && isDigit(s.peekNext()) {
		isFloat = true

		// Consume the "."
		s.advance()

		for isDigit(s.peek()) {
			s.advance()
		}
	}

	if isFloat {
		f, _ := strconv.ParseFloat(s.source[s.start:s.current], 64)
		s.addToken(FLOAT_LITERAL, f)
	} else {
		n, _ := strconv.Atoi(s.source[s.start:s.current])
		s.addToken(NUMBER_LITERAL, n)
	}
}

func (s *Scanner) identifier() {
	for isAlphaNumeric(s.peek()) {
		s.advance()
	}

	text := s.source[s.start:s.current]
	token, ok := keywords[text]
	if !ok {
		token = IDENTIFIER
	}
	s.addToken(token, text)
}

func (s *Scanner) scanToken() (err error) {
	c := s.advance()

	if t, ok := complexTokens[c]; ok {
		s.switchToken(t.patterns, t.tokenTypes)
		return
	}

	switch c {
	case '/':
		if s.match('/') {
			for s.peek() != '\n' && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for s.peekN(2) != "*/" && !s.isAtEnd() {
				s.advance()
			}
			s.advance()
			s.advance()
		} else if s.match('=') {
			s.addToken(DIV_ASSIGN, nil)
		} else {
			s.addToken(SLASH, nil)
		}
	case '\n':
		s.line++
		s.column = 1
		s.addToken(NEW_LINE, nil)
	case ' ':
	case '\r':
	case '\t':
	case '"':
		err = s.string()
	default:
		if isDigit(c) {
			s.number()
		} else if isAlpha(c) {
			s.identifier()
		} else {
			err = errutils.Error(s.line, s.column, fmt.Sprintf("unexpected '%c'.", c))
		}
	}

	return
}

var complexTokens = map[byte]struct {
	patterns   []string
	tokenTypes []TokenType
}{
	'(':  {[]string{}, []TokenType{LEFT_PAREN}},
	')':  {[]string{}, []TokenType{RIGHT_PAREN}},
	'{':  {[]string{}, []TokenType{LEFT_BRACE}},
	'}':  {[]string{}, []TokenType{RIGHT_BRACE}},
	':':  {[]string{}, []TokenType{COLON}},
	';':  {[]string{}, []TokenType{SEMICOLON}},
	'\'': {[]string{}, []TokenType{QUOTE}},
	',':  {[]string{}, []TokenType{COMMA}},
	'@':  {[]string{}, []TokenType{AT}},
	'#':  {[]string{}, []TokenType{TAG}},
	'.':  {[]string{"."}, []TokenType{RANGE_DOT, DOT}},
	'+':  {[]string{"+", "="}, []TokenType{INCREMENT, ADD_ASSIGN, PLUS}},
	'-':  {[]string{"-", "="}, []TokenType{DECREMENT, SUB_ASSIGN, MINNUS}},
	'*':  {[]string{"**=", "**", "*="}, []TokenType{POW_ASSIGN, POW, MUL_ASSIGN, STAR}},
	'%':  {[]string{"="}, []TokenType{MOD_ASSIGN, MOD}},
	'&':  {[]string{"&", "="}, []TokenType{AND_LOGIC, AND_ASSIGN, AND_BITWISE}},
	'|':  {[]string{"|", ">", "="}, []TokenType{OR_LOGIC, PIPELINE_RIGHT, OR_ASSIGN, OR_BITWISE}},
	'^':  {[]string{"="}, []TokenType{XOR_ASSIGN, XOR_BITWISE}},
	'~':  {[]string{"&=", "&", "|=", "|", "^=", "^"}, []TokenType{NAND_ASSIGN, NAND_BITWISE, NOR_ASSIGN, NOR_BITWISE, XNOR_ASSIGN, XNOR_BITWISE, NOT_BITWISE}},
	'<':  {[]string{"|", "!>", "!", "<=", "<<=", "<<", "<"}, []TokenType{PIPELINE_LEFT, GO_BI, GO_IN, BITSHIFT_LEFT_ASSIGN, ROUNDSHIFT_LEFT_ASSIGN, ROUNDSHIFT_LEFT, SHIFT_LEFT, LESS}},
	'>':  {[]string{">=", ">>=", ">>", ">"}, []TokenType{BITSHIFT_RIGHT_ASSIGN, ROUNDSHIFT_RIGHT_ASSIGN, ROUNDSHIFT_RIGHT, SHIFT_RIGHT, GREATER}},
	'!':  {[]string{".", ">", "="}, []TokenType{BANG_NAV, GO_OUT, NOT_ASSIGN, BANG}},
	'?':  {[]string{".", ":"}, []TokenType{CHECK_NAV, ELVIS, CHECK}},
	'=':  {[]string{">", "="}, []TokenType{RETURN, EQUAL, ASSIGN}},
}
