package parser

import (
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Expr interface {
	// interpret()
	String() string
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
