package parser

import "fmt"

func parenthesize(name string, exprs ...Expr) string {
	parts := ""
	parts += "(" + name
	for _, e := range exprs {
		parts += " " + e.String()
	}
	parts += ")"

	return parts
}

func (x Equality) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Comparison) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Term) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Factor) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Unary) String() string {
	return parenthesize(x.Operator.Lexeme, x.Right)
}

func (x Grouping) String() string {
	return parenthesize("group ", x.Expression)
}

func (x Sequence) String() string {
	return fmt.Sprintf("%v ; %v", x.Left, x.Right)
}

func (x Identifier) String() string {
	return fmt.Sprintf("%v", x.Name)
}

func (x Assign) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Pipeline) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Cast) String() string {
	return fmt.Sprintf("(%v:%s)", x.Left, x.TypeCast)
}

func (x Ternary) String() string {
	return fmt.Sprintf("(%v?%v:%v)", x.Expression, x.True, x.False)
}

func (x Range) String() string {
	return parenthesize("..", x.Left, x.Right)
}

func (x Logic) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Bitshift) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Bitwise) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Access) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Power) String() string {
	return parenthesize("**", x.Left, x.Right)
}

func (x Elvis) String() string {
	if x.ReturnZero {
		return fmt.Sprintf("(%v?: <0>)", x.Left)
	} else {
		return parenthesize("?:", x.Left, x.Right)
	}
}

func (x PositionAccess) String() string {
	return fmt.Sprintf("%v[%v]", x.Expression, x.Pos)
}

func (x Increment) String() string {
	if x.Position {
		return fmt.Sprintf("(%v%v)", x.Expression, x.Operator.Lexeme)
	} else {
		return fmt.Sprintf("(%v%v)", x.Operator.Lexeme, x.Expression)
	}
}

func (x Pointer) String() string {
	return parenthesize(x.Operator.Lexeme, x.Right)
}

func (x Check) String() string {
	if x.HaveReturn {
		return fmt.Sprintf("(%v) ? => %v", x.Left, x.Right)
	} else {
		return fmt.Sprintf("(%v)?", x.Left)
	}
}

func (x ArrayLiteral) String() string {
	return fmt.Sprintf("([%s: %v]%v)", x.Typing.Lexeme, x.Size, x.Values)
}

func (x Literal) String() string {
	return fmt.Sprintf("%v", x.Value)
}

func (x Type) String() string {
	return fmt.Sprintf("(%v)", x.Name.Lexeme)
}
