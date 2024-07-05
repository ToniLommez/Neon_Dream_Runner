package parser

import (
	"fmt"
	"strings"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

func (s *Scope) FnEval(f FnStmt) (any, error) {
	var err error
	args := make([]Stmt, len(f.Args))
	for i, fnArg := range f.Args {
		var literal any
		if literal, err = s.evaluate(fnArg.Initializer); err != nil {
			return nil, err
		}
		fnArg.Initializer = Literal{Value: literal}
		args[i] = fnArg
	}
	f.Context.Statements = append(args, f.Context.Statements...)

	res, err := f.Context.Interpret()
	if err != nil {
		return nil, err
	}

	if len(f.Return) == 0 {
		return nil, nil
	}

	// Creating the type cast of the result
	tokenType := l.TokenType(typeToString(f.Return[0]))
	token := l.Token{Type: tokenType}
	typeCast := Type{Name: token}
	result := Literal{Value: res}
	x := Cast{Left: result, TypeCast: typeCast}

	return s.evaluate(x)
}

func (s *Scope) IfEval(i IfStmt) (any, error) {
	test, err := s.evaluate(i.Condition)
	if err != nil {
		return nil, err
	}

	truthy, err := Truthy(test)
	if err != nil {
		return nil, err
	}

	if truthy {
		return s.evaluate(i.Then)
	} else if i.Else != nil {
		return s.evaluate(i.Else)
	}

	return nil, nil
}

func (s *Scope) LetEval(l LetStmt) (any, error) {
	var value any
	var err error

	if l.Initializer != nil {
		value, err = s.evaluate(l.Initializer)
		if err != nil {
			return nil, err
		}

		valueType := getType(value)
		if valueType == UNKNOWN || valueType == UNDEFINED {
			return nil, e.Error(l.Name.Line, 0, "", e.RUNTIME, fmt.Sprintf("let statement evaluate to unknown type: %v", value))
		}

		if l.Type == UNDEFINED {
			l.Type = valueType
		} else if l.Type != valueType {
			return nil, e.Error(l.Name.Line, 0, "", e.RUNTIME, fmt.Sprintf("let statement expected %s, found %s", typeToString(l.Type), typeToString(valueType)))
		}

		if valueType == NIL && !l.Nullable {
			return nil, e.Error(l.Name.Line, 0, "", e.RUNTIME, "non nullable let statement received nil value")
		}
	}

	_, err = s.Define(l, value)
	return nil, err
}

func (s *Scope) PutEval(p PutStmt) (any, error) {
	expr, err := s.evaluate(p.Value)
	if err != nil {
		return nil, err
	}

	tmp := fmt.Sprintf("%v", expr)
	color := "\033[38;2;150;240;240m"
	reset := "\033[0m"
	fmt.Printf("%s%v%s", color, strings.Replace(tmp, "\\n", "\n", -1), reset)

	// TODO: remove this after implement printf
	// fmt.Printf("\n")

	return nil, nil
}

func (s *Scope) WhileEval(w WhileStmt) (any, error) {
	var err error
	var c any  // Raw condition
	var b bool // Truthy(Raw condition)
	var i int

	for {
		if c, err = s.evaluate(w.Condition); err != nil {
			return nil, err
		}

		if b, err = Truthy(c); err != nil {
			return nil, err
		}

		if !b {
			break
		}

		if _, err = s.evaluate(w.Body); err != nil {
			return nil, err
		}

		i++
	}

	return nil, nil
}

func (s *Scope) ExprEval(e ExprStmt) (any, error) {
	return s.evaluate(e.Expr)
}

func (s *Scope) BlockEval(b Block) (any, error) {
	b.Scope.Parent = s
	b.Scope.Owner = s.Owner
	return b.Scope.Interpret()
}
