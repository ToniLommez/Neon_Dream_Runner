package lexer

import (
	"fmt"
	"strings"
)

type Token struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
	Line    int
	Column  int
}

func (t Token) String() (str string) {
	if t.Type != EOF && t.Type != NEW_LINE {
		lexeme := strings.Replace(t.Lexeme, "\n", "\\n", -1)
		lexeme = strings.Replace(lexeme, "\t", "\\t", -1)
		if t.Literal == nil {
			str = fmt.Sprintf("[%s, %s]", t.Type, lexeme)
		} else {
			str = fmt.Sprintf("[%s, %s, %v]", t.Type, lexeme, t.Literal)
		}
	}
	return
}
