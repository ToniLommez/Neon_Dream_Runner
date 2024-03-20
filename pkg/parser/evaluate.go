package parser

import (
	"fmt"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/utils"
)

func Evaluate(expr Expr) (any, error) {
	switch v := expr.(type) {
	case Sequence:
		return SequenceEval(v)
	case Assign:
		return nil, nil
	case Pipeline:
		return nil, nil
	case Ternary:
		return TernaryEval(v)
	case Range:
		return nil, nil
	case Logic:
		return LogicEval(v)
	case Equality:
		return EqualityEval(v)
	case Comparison:
		return ComparisonEval(v)
	case Bitshift:
		return nil, nil
	case Bitwise:
		return nil, nil
	case Term:
		return TermEval(v)
	case Factor:
		return FactorEval(v)
	case Power:
		return nil, nil
	case Increment:
		return nil, nil
	case Pointer:
		return nil, nil
	case Unary:
		return UnaryEval(v)
	case Access:
		return nil, nil
	case PositionAccess:
		return nil, nil
	case Elvis:
		return nil, nil
	case Check:
		return nil, nil
	case Cast:
		return nil, nil
	case Identifier:
		return nil, nil
	case Literal:
		return v.Value, nil
	case Type:
		return nil, nil
	case ArrayLiteral:
		return nil, nil
	case Grouping:
		return Evaluate(v.Expression)
	case nil:
		return nil, nil
	default:
		fmt.Println("literal evaluation not implemented, returning false")
		return nil, nil
	}
}

func Truthy(value any) (r bool, err error) {
	switch v := value.(type) {
	case nil:
		r = false
	case bool:
		r = v
	case int:
		r = v != 0
	case uint:
		r = v != 0
	case float64:
		r = v != 0.0
	case string:
		r = v != ""
	default:
		r = false
		err = errutils.Error(-1, -1, "", errutils.RUNTIME, "truthy pattern matching not implemented, returning false")
	}

	return
}

func SequenceEval(s Sequence) (any, error) {
	if _, err := Evaluate(s.Left); err != nil {
		return nil, err
	}

	return Evaluate(s.Right)
}

func TernaryEval(t Ternary) (any, error) {
	test, err := Evaluate(t.Expression)
	if err != nil {
		return nil, err
	}

	truthy, err := Truthy(test)
	if err != nil {
		return nil, err
	}

	if truthy {
		return Evaluate(t.True)
	} else {
		return Evaluate(t.False)
	}
}

func LogicEval(x Logic) (res interface{}, err error) {
	l, err := Evaluate(x.Left)
	if err != nil {
		return
	}

	r, err := Evaluate(x.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r)
	if precedence == UNKNOWN {
		err = errutils.Error(x.Operator.Line, x.Operator.Column, x.Operator.Lexeme, errutils.RUNTIME, "invalid operands")
	}

	switch x.Operator.Type {
	case lexer.OR_LOGIC:
		switch precedence {
		case BOOL:
			res = l.(bool) || r.(bool)
		case UINT:
			res = (l.(uint) > 0) || (r.(uint) > 0)
		case INT:
			res = (l.(int) > 0) || (r.(int) > 0)
		case FLOAT:
			res = (l.(float64) > 0) || (r.(float64) > 0)
		case STRING:
			err = errutils.Error(x.Operator.Line, x.Operator.Column, x.Operator.Lexeme, errutils.RUNTIME, "cannot use string as logic value")
		}
	case lexer.AND_LOGIC:
		switch precedence {
		case BOOL:
			res = l.(bool) && r.(bool)
		case UINT:
			res = (l.(uint) > 0) && (r.(uint) > 0)
		case INT:
			res = (l.(int) > 0) && (r.(int) > 0)
		case FLOAT:
			res = (l.(float64) > 0) && (r.(float64) > 0)
		case STRING:
			err = errutils.Error(x.Operator.Line, x.Operator.Column, x.Operator.Lexeme, errutils.RUNTIME, "cannot use string as logic value")
		}
	}

	return
}

func EqualityEval(e Equality) (res interface{}, err error) {
	l, err := Evaluate(e.Left)
	if err != nil {
		return
	}

	r, err := Evaluate(e.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r)
	if precedence == UNKNOWN {
		err = errutils.Error(e.Operator.Line, e.Operator.Column, e.Operator.Lexeme, errutils.RUNTIME, "invalid operands")
	}

	switch e.Operator.Type {
	case lexer.EQUAL:
		switch precedence {
		case BOOL:
			res = l.(bool) == r.(bool)
		case UINT:
			res = l.(uint) == r.(uint)
		case INT:
			res = l.(int) == r.(int)
		case FLOAT:
			res = l.(float64) == r.(float64)
		case STRING:
			res = l.(string) == r.(string)
		}
	case lexer.NOT_EQUAL:
		switch precedence {
		case BOOL:
			res = l.(bool) != r.(bool)
		case UINT:
			res = l.(uint) != r.(uint)
		case INT:
			res = l.(int) != r.(int)
		case FLOAT:
			res = l.(float64) != r.(float64)
		case STRING:
			res = l.(string) != r.(string)
		}
	}

	return
}

func ComparisonEval(c Comparison) (res interface{}, err error) {
	l, err := Evaluate(c.Left)
	if err != nil {
		return
	}

	r, err := Evaluate(c.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r)
	if precedence == UNKNOWN {
		err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "invalid operands")
	}

	switch c.Operator.Type {
	case lexer.GREATER_EQUAL:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) >= boolToInt(r.(bool))
		case UINT:
			res = l.(uint) >= r.(uint)
		case INT:
			res = l.(int) >= r.(int)
		case FLOAT:
			res = l.(float64) >= r.(float64)
		case STRING:
			err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot direct compare strings")
		}
	case lexer.LESS_EQUAL:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) <= boolToInt(r.(bool))
		case UINT:
			res = l.(uint) <= r.(uint)
		case INT:
			res = l.(int) <= r.(int)
		case FLOAT:
			res = l.(float64) <= r.(float64)
		case STRING:
			err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot direct compare strings")
		}
	case lexer.GREATER:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) > boolToInt(r.(bool))
		case UINT:
			res = l.(uint) > r.(uint)
		case INT:
			res = l.(int) > r.(int)
		case FLOAT:
			res = l.(float64) > r.(float64)
		case STRING:
			err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot direct compare strings")
		}
	case lexer.LESS:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) < boolToInt(r.(bool))
		case UINT:
			res = l.(uint) < r.(uint)
		case INT:
			res = l.(int) < r.(int)
		case FLOAT:
			res = l.(float64) < r.(float64)
		case STRING:
			err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot direct compare strings")
		}
	}

	return
}

func UnaryEval(u Unary) (res any, err error) {
	v, err := Evaluate(u.Right)
	if err != nil {
		return nil, err
	}

	t := getType(v)
	o := u.Operator

	switch o.Type {
	case lexer.BANG:
		switch t {
		case INT:
			res = utils.Ternary(v.(int) == 0, true, false)
		case BOOL:
			res = !v.(bool)
		case FLOAT:
			res = utils.Ternary(v.(float64) == 0, true, false)
		case STRING:
			err = errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, fmt.Sprintf("expect number after !, received string: \"%v\"", v))
		default:
			err = errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, fmt.Sprintf("expect number after !, received: %v", v))
		}
	case lexer.NOT_BITWISE:
		switch t {
		case INT:
			res = ^v.(int)
		case BOOL:
			res = !v.(bool)
		case FLOAT:
			err = errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, "operator ~ not defined on float")
		default:
			return nil, errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, "expect number after ~")
		}
	case lexer.PLUS:
		res = v
	case lexer.MINUS:
		switch t {
		case INT:
			res = -v.(int)
		case BOOL:
			err = errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, "operator - not defined on bool")
		case FLOAT:
			res = -v.(float64)
		default:
			return nil, errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, "expect number after ~")
		}
	case lexer.GO_IN:
		err = errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, "parallelism not yet implemented")
	default:
		err = errutils.Error(o.Line, o.Column, o.Lexeme, errutils.RUNTIME, "unknown operator")
	}

	return res, err
}

func TermEval(t Term) (res interface{}, err error) {
	l, err := Evaluate(t.Left)
	if err != nil {
		return
	}

	r, err := Evaluate(t.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r)
	if precedence == UNKNOWN {
		err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot convert operands")
	}

	switch t.Operator.Type {
	case lexer.PLUS:
		switch precedence {
		case BOOL:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot sum bool values")
		case UINT:
			res = l.(uint) + r.(uint)
		case INT:
			res = l.(int) + r.(int)
		case FLOAT:
			res = l.(float64) + r.(float64)
		case STRING:
			res = l.(string) + r.(string)
		}
	case lexer.MINUS:
		switch precedence {
		case BOOL:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot sum bool values")
		case UINT:
			res = l.(uint) - r.(uint)
		case INT:
			res = l.(int) - r.(int)
		case FLOAT:
			res = l.(float64) - r.(float64)
		case STRING:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot subtract strings")
		}
	}

	return
}

func FactorEval(t Factor) (res interface{}, err error) {
	l, err := Evaluate(t.Left)
	if err != nil {
		return
	}

	r, err := Evaluate(t.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r)
	if precedence == UNKNOWN {
		err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot implicit convert operands")
	}

	switch t.Operator.Type {
	case lexer.STAR:
		switch precedence {
		case BOOL:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot multiply bool values")
		case UINT:
			res = l.(uint) * r.(uint)
		case INT:
			res = l.(int) * r.(int)
		case FLOAT:
			res = l.(float64) * r.(float64)
		case STRING:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot multiply string values")
		}
	case lexer.SLASH:
		switch precedence {
		case BOOL:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot sum bool values")
		case UINT:
			res = l.(uint) - r.(uint)
		case INT:
			res = l.(int) - r.(int)
		case FLOAT:
			res = l.(float64) - r.(float64)
		case STRING:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot subtract strings")
		}
	case lexer.MOD:
		switch precedence {
		case BOOL:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot sum bool values")
		case UINT:
			res = l.(uint) - r.(uint)
		case INT:
			res = l.(int) - r.(int)
		case FLOAT:
			res = l.(float64) - r.(float64)
		case STRING:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot subtract strings")
		}
	}

	return
}

// TODO: implementar erro
func typePrecedence(l any, r any) (any, any, int) {
	lType := getType(l)
	rType := getType(r)

	switch lType {
	case FLOAT:
		switch rType {
		case FLOAT:
			return l, r, FLOAT
		case INT:
			return l, float64(r.(int)), FLOAT
		case BOOL:
			return l, utils.Ternary(r.(bool), 1.0, 0.0), FLOAT
		default:
			return l, r, UNKNOWN
		}
	case INT:
		switch rType {
		case FLOAT:
			return float64(l.(int)), r, FLOAT
		case INT:
			return l, r, INT
		case BOOL:
			return l, utils.Ternary(r.(bool), 1, 0), INT
		default:
			return l, r, UNKNOWN
		}
	case UINT:
		switch rType {
		case UINT:
			return l, r, UINT
		case BOOL:
			return l, utils.Ternary(r.(bool), uint(1), uint(0)), UINT
		default:
			return l, r, UNKNOWN
		}
	case BOOL:
		switch rType {
		case FLOAT:
			return utils.Ternary(l.(bool), 1.0, 0.0), r, FLOAT
		case INT:
			return utils.Ternary(l.(bool), 1, 0), r, INT
		case UINT:
			return utils.Ternary(l.(bool), uint(1), uint(0)), r, UINT
		case BOOL:
			return l, r, BOOL
		default:
			return l, r, UNKNOWN
		}
	case STRING:
		switch rType {
		case STRING:
			return l, r, STRING
		default:
			return l, r, UNKNOWN
		}
	default:
		return l, r, UNKNOWN
	}
}