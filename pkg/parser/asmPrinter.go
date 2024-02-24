package parser

import "fmt"

func (x Binary) String() string {
	return parenthesize(x.Operator.Lexeme, x.Left, x.Right)
}

func (x Unary) String() string {
	return parenthesize(x.Operator.Lexeme, x.Right)
}

func (x Grouping) String() string {
	return parenthesize("group ", x.Expression)
}

func (x Literal) String() string {
	return fmt.Sprintf("%v", x.Value)
}

func (x Sequence) String() string {
	return fmt.Sprintf("%v ; %v", x.Left, x.Right)
}

func (x Identifier) String() string {
	return fmt.Sprintf("%v", x.Name.Lexeme)
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

func (x Type) String() string {
	return x.Name
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
	return fmt.Sprintf("([%s: %v]%v)", x.Typing.Type, x.Size, x.Values)
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

// Type

func (x Binary) Type() string {
	return "Binary"
}

func (x Unary) Type() string {
	return "Unary"
}

func (x Grouping) Type() string {
	return "Grouping"
}

func (x Literal) Type() string {
	return "Literal"
}

func (x Sequence) Type() string {
	return "Sequence"
}

func (x Identifier) Type() string {
	return "Identifier"
}

func (x Assign) Type() string {
	return "Assign"
}

func (x Pipeline) Type() string {
	return "Pipeline"
}

func (x Cast) Type() string {
	return "Cast"
}

func (x Type) Type() string {
	return "Type"
}

func (x Ternary) Type() string {
	return "Ternary"
}

func (x Range) Type() string {
	return "Range"
}

func (x Logic) Type() string {
	return "Logic"
}

func (x Bitshift) Type() string {
	return "Bitshift"
}

func (x Bitwise) Type() string {
	return "Bitwise"
}

func (x Power) Type() string {
	return "Power"
}

func (x Increment) Type() string {
	return "Increment"
}

func (x Pointer) Type() string {
	return "Pointer"
}

func (x Access) Type() string {
	return "Access"
}

func (x PositionAccess) Type() string {
	return "PositionAccess"
}

func (x Elvis) Type() string {
	return "Elvis"
}

func (x Check) Type() string {
	return "Check"
}

func (x ArrayLiteral) Type() string {
	return "ArrayLiteral"
}
