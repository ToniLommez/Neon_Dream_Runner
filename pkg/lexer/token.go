package lexer

import "fmt"

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
	Column  int
}

func (t Token) String() (str string) {
	if t.Type != EOF {
		if t.Literal == nil {
			str = fmt.Sprintf("[%s, %s]", t.Type, t.Lexeme)
		} else {
			str = fmt.Sprintf("[%s, %s, %+v]", t.Type, t.Lexeme, t.Literal)
		}
	}
	return
}
