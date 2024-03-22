package parser

type Stmt interface {
	// String() string
}

type ExprStmt struct {
	Expr Expr
}

type PutStmt struct {
	Value Expr
}
