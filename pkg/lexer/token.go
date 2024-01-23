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
	if t.Type == NEW_LINE {
		str = fmt.Sprintf("[%s]", t.Type)
	} else if t.Type != EOF {
		lexeme := strings.Replace(t.Lexeme, "\n", "\\n", -1)
		lexeme = strings.Replace(lexeme, "\t", "\\t", -1)
		if t.Literal == nil {
			str = fmt.Sprintf("[%s, %s]", t.Type, lexeme)
		} else {
			str = fmt.Sprintf("[%s, %s, %q]", t.Type, lexeme, t.Literal)
		}
	}
	return
}
