package parser

import "fmt"

func (b Binary) String() string {
	return parenthesize(b.Operator.Lexeme, b.Left, b.Right)
}

func (u Unary) String() string {
	return parenthesize(u.Operator.Lexeme, u.Right)
}

func (g Grouping) String() string {
	return parenthesize("group ", g.Expression)
}

func (l Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}

func parenthesize(name string, exprs ...Expr) string {
	parts := ""
	parts += "(" + name
	for _, e := range exprs {
		parts += " " + e.String()
	}
	parts += ")"

	return parts
}
