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

func (s *Scope) LetEval(x LetStmt) (any, error) {
	// var typing any
	var initType int
	var init any
	var err error

	if x.Initializer != nil {
		if init, err = s.evaluate(x.Initializer); err != nil {
			return nil, err
		}

		// TODO: fix this shit...
		switch i := init.(type) {
		case ArrayType:
			vs := make([]any, i.MaxSize.(int))
			init = Slice{Typing: i.Typing, Size: i.MaxSize.(int), Values: vs, IsArray: true}
		}

		initType = getType(init)
		if initType == UNKNOWN || initType == UNDEFINED {
			return nil, e.Error(x.Name.Line, 0, "", e.RUNTIME, fmt.Sprintf("let statement evaluate to unknown type: %v", init))
		} else if initType == SLICE {
			x.IsSlice = true
		}

		// x.Type could be a Type or a Slice
		if x.Type == nil {
			if initType == SLICE {
				initSlice := init.(Slice)
				x.Type = ArrayType{Typing: initSlice.Typing, MaxSize: initSlice.Size, IsSlice: !initSlice.IsArray}
			} else {
				x.Type = Type{Name: l.Token{Type: anyToToken(init)}}
			}
		} else if initType != getType(x.Type) {
			if getType(x.Type) == FLOAT && initType == INT { // HARD CODED BUG, TODO: REMOVE THIS
				initType = FLOAT
			} else {
				return nil, e.Error(x.Name.Line, 0, "", e.RUNTIME, fmt.Sprintf("let statement expected %s, found %s", typeToString(getType(x.Type)), typeToString(initType)))
			}
		}

		if initType == NIL && !x.Nullable {
			return nil, e.Error(x.Name.Line, 0, "", e.RUNTIME, "non nullable let statement received nil value")
		}
	}

	// TODO: define type without initialize

	_, err = s.Define(x, init)
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

	formatted := strings.Replace(tmp, "\\n", "\n", -1)
	formatted = strings.Replace(formatted, "\\x1b[H", "\x1b[H", -1)

	fmt.Printf("%s%v%s", color, formatted, reset)
	if p.NewLine {
		fmt.Printf("\n")
	}

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
