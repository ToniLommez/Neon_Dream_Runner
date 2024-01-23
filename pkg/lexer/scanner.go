package lexer

import (
	"fmt"

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
		column:  -1,
	}
}

func (s Scanner) isAtEnd() bool {
	return s.current >= len(s.source)
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

func (s *Scanner) advance() byte {
	s.current++
	s.column++
	return s.source[s.current-1]
}

func (s *Scanner) addToken(tokenType TokenType, literal interface{}) {
	text := s.source[s.start:s.current]
	s.tokens = append(s.tokens, Token{Type: tokenType, Lexeme: text, Literal: literal, Line: s.line, Column: s.column})
}

func (s *Scanner) scanToken() (err error) {
	char := s.advance()

	switch char {
	case '(':
		s.addToken(LEFT_PAREN, nil)
	case ')':
		s.addToken(RIGHT_PAREN, nil)
	case '{':
		s.addToken(LEFT_BRACE, nil)
	case '}':
		s.addToken(RIGHT_BRACE, nil)
	case '.':
		if s.match('.') {
			s.addToken(RANGE_DOT, nil)
		} else {
			s.addToken(DOT, nil)
		}
	case ',':
		s.addToken(COMMA, nil)
	case '@':
		s.addToken(AT, nil)
	case '#':
		s.addToken(TAG, nil)
	case '+':
		if s.match('+') {
			s.addToken(INCREMENT, nil)
		} else if s.match('=') {
			s.addToken(ADD_ASSIGN, nil)
		} else {
			s.addToken(PLUS, nil)
		}
	case '-':
		if s.match('-') {
			s.addToken(DECREMENT, nil)
		} else if s.match('=') {
			s.addToken(SUB_ASSIGN, nil)
		} else {
			s.addToken(MINNUS, nil)
		}
	case '/':
		if s.match('/') {
			for s.peek(1) != "\n" && !s.isAtEnd() {
				s.advance()
			}
		} else if s.match('*') {
			for s.peek(2) != "*/" && !s.isAtEnd() {
				s.advance()
			}
			s.advance()
			s.advance()
		} else if s.match('=') {
			s.addToken(DIV_ASSIGN, nil)
		} else {
			s.addToken(SLASH, nil)
		}
	case '*':
		if s.match('*') {
			if s.match('=') {
				s.addToken(POW_ASSIGN, nil)
			} else {
				s.addToken(POW, nil)
			}
		} else if s.match('=') {
			s.addToken(MUL_ASSIGN, nil)
		} else {
			s.addToken(STAR, nil)
		}
	case '%':
		if s.match('=') {
			s.addToken(MOD_ASSIGN, nil)
		} else {
			s.addToken(MOD, nil)
		}
	case '&':
		if s.match('&') {
			s.addToken(AND_LOGIC, nil)
		} else if s.match('=') {
			s.addToken(AND_ASSIGN, nil)
		} else {
			s.addToken(AND_BITWISE, nil)
		}
	case '|':
		if s.match('|') {
			s.addToken(OR_LOGIC, nil)
		} else if s.match('>') {
			s.addToken(PIPELINE_RIGHT, nil)
		} else if s.match('=') {
			s.addToken(OR_ASSIGN, nil)
		} else {
			s.addToken(OR_BITWISE, nil)
		}
	case '^':
		if s.match('=') {
			s.addToken(XOR_ASSIGN, nil)
		} else {
			s.addToken(XOR_BITWISE, nil)
		}
	case '~':
		if s.match('&') {
			if s.match('=') {
				s.addToken(NAND_ASSIGN, nil)
			} else {
				s.addToken(NAND_BITWISE, nil)
			}
		} else if s.match('|') {
			if s.match('=') {
				s.addToken(NOR_ASSIGN, nil)
			} else {
				s.addToken(NOR_BITWISE, nil)
			}
		} else if s.match('^') {
			if s.match('=') {
				s.addToken(XNOR_ASSIGN, nil)
			} else {
				s.addToken(XNOR_BITWISE, nil)
			}
		} else {
			s.addToken(NOT_BITWISE, nil)
		}
	case '<':
		if s.match('|') {
			s.addToken(PIPELINE_LEFT, nil)
		} else if s.match('!') {
			if s.match('>') {
				s.addToken(GO_BI, nil)
			} else {
				s.addToken(GO_IN, nil)
			}
		} else if s.match('<') {
			if s.match('=') {
				s.addToken(BITSHIFT_LEFT_ASSIGN, nil)
			} else if s.match('<') {
				if s.match('=') {
					s.addToken(ROUNDSHIFT_LEFT_ASSIGN, nil)
				} else {
					s.addToken(ROUNDSHIFT_LEFT, nil)
				}
			} else {
				s.addToken(SHIFT_LEFT, nil)
			}
		} else {
			s.addToken(LESS, nil)
		}
	case '>':
		if s.match('>') {
			if s.match('=') {
				s.addToken(BITSHIFT_RIGHT_ASSIGN, nil)
			} else if s.match('>') {
				if s.match('=') {
					s.addToken(ROUNDSHIFT_RIGHT_ASSIGN, nil)
				} else {
					s.addToken(ROUNDSHIFT_RIGHT, nil)
				}
			} else {
				s.addToken(SHIFT_RIGHT, nil)
			}
		} else {
			s.addToken(GREATER, nil)
		}
	case '!':
		if s.match('.') {
			s.addToken(BANG_NAV, nil)
		} else if s.match('>') {
			s.addToken(GO_OUT, nil)
		} else if s.match('=') {
			s.addToken(NOT_ASSIGN, nil)
		} else {
			s.addToken(BANG, nil)
		}
	case '?':
		if s.match('.') {
			s.addToken(CHECK_NAV, nil)
		} else if s.match(':') {
			s.addToken(ELVIS, nil)
		} else {
			s.addToken(CHECK, nil)
		}
	case ':':
		s.addToken(COLON, nil)
	case ';':
		s.addToken(SEMICOLON, nil)
	case '\'':
		s.addToken(QUOTE, nil)
	case '=':
		if s.match('=') {
			s.addToken(RETURN, nil)
		} else if s.match('=') {
			s.addToken(EQUAL, nil)
		} else {
			s.addToken(ASSIGN, nil)
		}

	// Keywords
	// case 'let': s.addToken(LET, nil)
	// case 'let!': s.addToken(LET_BANG, nil)
	// case 'let?': s.addToken(LET_CHECK, nil)
	// case 'let!?': s.addToken(LET_BANG_CHECK, nil)
	// case 'fn': s.addToken(FN, nil)
	// case 'asm': s.addToken(ASM, nil)
	// case 'for': s.addToken(FOR, nil)
	// case 'loop': s.addToken(LOOP, nil)
	// case 'while': s.addToken(WHILE, nil)
	// case 'until': s.addToken(UNTIL, nil)
	// case 'do': s.addToken(DO, nil)
	// case 'in': s.addToken(IN, nil)
	// case 'pulse': s.addToken(PULSE, nil)
	// case 'before': s.addToken(BEFORE, nil)
	// case 'inside': s.addToken(INSIDE, nil)
	// case 'after': s.addToken(AFTER, nil)
	// case 'error': s.addToken(ERROR, nil)
	// case 'nil': s.addToken(NIL, nil)
	// case 'case': s.addToken(CASE, nil)
	// case 'of': s.addToken(OF, nil)
	// case 'if': s.addToken(IF, nil)
	// case 'else': s.addToken(ELSE, nil)
	// case 'elif': s.addToken(ELIF, nil)
	// case 'use': s.addToken(USE, nil)
	// case 'as': s.addToken(AS, nil)
	// case 'merge': s.addToken(MERGE, nil)
	// case 'obj': s.addToken(OBJ, nil)
	// case 'pub': s.addToken(PUB, nil)
	// case 'when': s.addToken(WHEN, nil)
	// case 'trigger': s.addToken(TRIGGER, nil)
	// case 'trait': s.addToken(TRAIT, nil)
	// case 'this': s.addToken(THIS, nil)

	// Reserved
	// case 'print': s.addToken(PRINT, nil)
	// case 'println': s.addToken(PRINTLN, nil)
	// case 'printf': s.addToken(PRINTF, nil)
	// case 'true': s.addToken(TRUE, nil)
	// case 'false': s.addToken(FALSE, nil)

	// Types
	case '"':
		err = s.string()
	// case 'bool': s.addToken(BOOL, nil)
	// case 'int': s.addToken(INT, nil)
	// case 'float': s.addToken(FLOAT, nil)
	// case 'any': s.addToken(ANY, nil)

	case '\n':
		s.line++
		s.column = 0
		s.addToken(NEW_LINE, nil)
	case ' ':
	case '\r':
	case '\t':
	default:
		err = errutils.Error(s.line, s.column, fmt.Sprintf("unexpected '%c'.", char))
	}

	return
}

func (s *Scanner) match(expected byte) bool {
	if s.isAtEnd() || s.source[s.current] != expected {
		return false
	}

	s.current++
	s.column++
	return true
}

func (s Scanner) peek(size int) string {
	if s.isAtEnd() {
		return "\x00"
	}

	end := s.current + size
	if end > len(s.source) {
		end = len(s.source)
	}

	return s.source[s.current:end]
}

func (s *Scanner) string() error {
	for s.peek(1) != "\"" && !s.isAtEnd() {
		if s.peek(1) == "\n" {
			s.line++
			s.column = 0
		}
		s.advance()
	}

	if s.isAtEnd() {
		return errutils.Error(s.line, s.column, "unterminated string.")
	}

	// the closing ".
	s.advance()

	// trim quotes
	str := s.source[s.start+1 : s.current-1]
	s.addToken(STRING, str)

	return nil
}
