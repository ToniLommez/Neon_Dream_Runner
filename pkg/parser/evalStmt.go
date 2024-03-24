package parser

import (
	"fmt"
	"strings"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
)

// TODO: auto convert initializer type if let have a preset type
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

	return nil, s.Values.Define(l, value)
}

func (s *Scope) PutEval(p PutStmt) (any, error) {
	expr, err := s.evaluate(p.Value)
	if err != nil {
		return nil, err
	}

	tmp := fmt.Sprintf("%v", expr)
	fmt.Printf("%v", strings.Replace(tmp, "\\n", "\n", -1))

	// TODO: remove this after implement printf
	fmt.Printf("\n")

	return expr, nil
}

func (s *Scope) ExprEval(e ExprStmt) (any, error) {
	return s.evaluate(e.Expr)
}
