package parser

import "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"

type Stmt interface {
	// String() string
}

type Scope struct {
	Statements []Stmt
	Values     Environment
	Previous   *Scope
}

type ExprStmt struct {
	Expr Expr
}

type PutStmt struct {
	Value Expr
}

type LetStmt struct {
	Name        lexer.Token
	Mutable     bool
	Nullable    bool
	Type        int
	Initializer Expr
}
