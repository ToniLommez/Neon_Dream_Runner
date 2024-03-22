package parser

import (
	"fmt"
	"strings"
)

func ExprStmtEval(s ExprStmt) (any, error) {
	return evaluate(s.Expr)
}

func PutStmtEval(p PutStmt) (any, error) {
	expr, err := evaluate(p.Value)
	if err != nil {
		return nil, err
	}

	tmp := fmt.Sprintf("%v", expr)
	fmt.Printf("%v", strings.Replace(tmp, "\\n", "\n", -1))

	// TODO: remove this after implement printf
	fmt.Printf("\n")

	return expr, nil
}
