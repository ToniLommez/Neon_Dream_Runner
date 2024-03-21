package parser

import (
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Expr interface {
	String() string
}

type Sequence struct {
	Left  Expr
	Right Expr
}

type Assign struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Pipeline struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Ternary struct {
	Expression Expr
	True       Expr
	False      Expr
}

type Range struct {
	Left  Expr
	Right Expr
}

type Logic struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Equality struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Comparison struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Bitshift struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Bitwise struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Term struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Factor struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Power struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

// Position = 0 is Left, 1 is Right
type Increment struct {
	Expression Expr
	Operator   l.Token
	Position   bool
}

type Pointer struct {
	Operator l.Token
	Right    Expr
}

type Unary struct {
	Operator l.Token
	Right    Expr
}

type Access struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type PositionAccess struct {
	Expression Expr
	Pos        Expr
}

type Elvis struct {
	Left       Expr
	ReturnZero bool
	Right      Expr
}

type Check struct {
	Left       Expr
	HaveReturn bool
	Right      Expr
}

type Cast struct {
	Left     Expr
	Operator l.Token
	TypeCast Expr
}

type Identifier struct {
	Name l.Token
}

type ArrayLiteral struct {
	Typing l.Token
	Size   Expr
	Values []Expr
}

type Literal struct {
	Value interface{}
}

type Type struct {
	Name l.Token
}

type Grouping struct {
	Expression Expr
}
