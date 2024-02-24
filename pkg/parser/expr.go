package parser

import (
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Expr interface {
	// interpret()
	String() string
	Type() string
}

type Ternary struct {
	Expression Expr
	True       Expr
	False      Expr
}

type Binary struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Unary struct {
	Operator l.Token
	Right    Expr
}

type Grouping struct {
	Expression Expr
}

type Literal struct {
	Value interface{}
}

type Type struct {
	Name string
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

type Identifier struct {
	Name l.Token
}

type Pipeline struct {
	Left     Expr
	Operator l.Token
	Right    Expr
}

type Cast struct {
	Left     Expr
	TypeCast l.TokenType
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

type Power struct {
	Left  Expr
	Right Expr
}

// Position = 0 is Left
//
// Position = 1 is Right
type Increment struct {
	Expression Expr
	Operator   l.Token
	Position   bool
}

type Pointer struct {
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

type ArrayLiteral struct {
	Typing l.Token
	Size   Expr
	Values []Expr
}
