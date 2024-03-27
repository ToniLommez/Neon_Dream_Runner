package parser

import "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"

type Stmt interface {
	// String() string
}

type WhileStmt struct {
	Condition Expr
	Body      Expr
}

// All `if/else` statements need to have blocks, and all blocks are statements
// so then/else have expressions, so `if` can act as a ternary too
// therefore `if` is actually an expression, regardless of the stmt in its name
type IfStmt struct {
	Condition Expr
	Then      Expr
	Else      Expr
}

type Block struct {
	Scope Scope
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
