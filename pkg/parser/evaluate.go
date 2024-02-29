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
		return nil, nil
	case Assign:
		return nil, nil
	case Pipeline:
		return nil, nil
	case Ternary:
		return TernaryEval(v)
	case Range:
		return nil, nil
	case Logic:
		return nil, nil
	case Equality:
		return nil, nil
	case Comparison:
		return nil, nil
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
	case int, int8, int16, int32, int64:
		r = v != 0
	case uint, uint8, uint16, uint32, uint64:
		r = v != 0
	case float32, float64:
		r = v != 0.0
	case string:
		r = v != ""
	default:
		r = false
		err = errutils.Error(-1, -1, "", errutils.RUNTIME, "truthy pattern matching not implemented, returning false")
	}

	return
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
	}

	return res, err
}

func TermEval(t Term) (res interface{}, err error) {
	left, err := Evaluate(t.Left)
	if err != nil {
		return
	}

	right, err := Evaluate(t.Right)
	if err != nil {
		return
	}

	leftConv, rightConv, isFloat := floatIntPrecedence(left, right)

	switch t.Operator.Type {
	case lexer.PLUS:
		if isFloat {
			res = Literal{leftConv.(float64) + rightConv.(float64)}
		} else {
			res = Literal{leftConv.(int) + rightConv.(int)}
		}
	case lexer.MINUS:
		if isFloat {
			res = Literal{leftConv.(float64) - rightConv.(float64)}
		} else {
			res = Literal{leftConv.(int) - rightConv.(int)}
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

	leftConv, rightConv, isFloat := floatIntPrecedence(l, r)

	switch t.Operator.Type {
	case lexer.STAR:
		if isFloat {
			res = leftConv.(float64) * rightConv.(float64)
		} else {
			res = leftConv.(int) * rightConv.(int)
		}
	case lexer.SLASH:
		if isFloat {
			res = leftConv.(float64) / rightConv.(float64)
		} else {
			res = leftConv.(int) / rightConv.(int)
		}
	case lexer.MOD:
		left, err := toInt(l)
		if err != nil {
			return nil, err
		}
		right, err := toInt(r)
		if err != nil {
			return nil, err
		}
		res = left % right
	}

	return
}

// TODO: implementar erro
func floatIntPrecedence(l any, r any) (any, any, bool) {
	lType := getType(l)
	rType := getType(r)

	if lType == FLOAT || rType == FLOAT {
		left, _ := toFloat(l)
		right, _ := toFloat(r)
		return left, right, true
	} else {
		left, _ := toInt(l)
		right, _ := toInt(r)
		return left, right, false
	}
}
