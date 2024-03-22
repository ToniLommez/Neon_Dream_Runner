package parser

import (
	"fmt"
	"math"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/utils"
)

// TODO: implement char
func typePrecedence(l any, r any, acceptUint bool) (any, any, int) {
	lType := getType(l)
	rType := getType(r)

	switch lType {
	case BOOL:
		switch rType {
		case BOOL:
			return l, r, BOOL
		case INT:
			return utils.Ternary(l.(bool), 1, 0), r, INT
		case UINT:
			return utils.Ternary(l.(bool), uint(1), uint(0)), r, UINT
		case FLOAT:
			return utils.Ternary(l.(bool), 1.0, 0.0), r, FLOAT
		default:
			return l, r, UNKNOWN
		}
	case CHAR:
		switch rType {
		case CHAR:
			return l, r, CHAR
		case STRING:
			return string(l.(rune)), r, STRING
		default:
			return l, r, UNKNOWN
		}
	case INT:
		switch rType {
		case BOOL:
			return l, utils.Ternary(r.(bool), 1, 0), INT
		case INT:
			return l, r, INT
		case UINT:
			if acceptUint {
				return uint(l.(int)), r, UINT
			} else {
				return l, r, UNKNOWN
			}
		case FLOAT:
			return float64(l.(int)), r, FLOAT
		default:
			return l, r, UNKNOWN
		}
	case UINT:
		switch rType {
		case BOOL:
			return l, utils.Ternary(r.(bool), uint(1), uint(0)), UINT
		case INT:
			if acceptUint {
				return l, uint(r.(int)), UINT
			} else {
				return l, r, UNKNOWN
			}
		case UINT:
			return l, r, UINT
		case FLOAT:
			return float64(l.(uint)), r, FLOAT
		default:
			return l, r, UNKNOWN
		}
	case FLOAT:
		switch rType {
		case BOOL:
			return l, utils.Ternary(r.(bool), 1.0, 0.0), FLOAT
		case INT:
			return l, float64(r.(int)), FLOAT
		case UINT:
			return l, float64(r.(uint)), FLOAT
		case FLOAT:
			return l, r, FLOAT
		default:
			return l, r, UNKNOWN
		}
	case STRING:
		switch rType {
		case CHAR:
			return l, string(r.(rune)), STRING
		case STRING:
			return l, r, STRING
		default:
			return l, r, UNKNOWN
		}
	default:
		return l, r, UNKNOWN
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
	if _, err := evaluate(s.Left); err != nil {
		return nil, err
	}

	return evaluate(s.Right)
}

func TernaryEval(t Ternary) (any, error) {
	test, err := evaluate(t.Expression)
	if err != nil {
		return nil, err
	}

	truthy, err := Truthy(test)
	if err != nil {
		return nil, err
	}

	if truthy {
		return evaluate(t.True)
	} else {
		return evaluate(t.False)
	}
}

func LogicEval(x Logic) (res interface{}, err error) {
	l, err := evaluate(x.Left)
	if err != nil {
		return
	}

	r, err := evaluate(x.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
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
	l, err := evaluate(e.Left)
	if err != nil {
		return
	}

	r, err := evaluate(e.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, true)
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
	l, err := evaluate(c.Left)
	if err != nil {
		return
	}

	r, err := evaluate(c.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, true)
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

func BitshiftEval(b Bitshift) (res interface{}, err error) {
	l, err := evaluate(b.Left)
	if err != nil {
		return
	}

	r, err := evaluate(b.Right)
	if err != nil {
		return
	}

	if _, _, precedence := typePrecedence(l, r, true); precedence == UNKNOWN {
		err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "invalid operands")
	}

	switch b.Operator.Type {
	case lexer.SHIFT_LEFT:
		switch getType(l) {
		case BOOL:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitshift bool")
		case INT:
			res = l.(int) << toUint(r)
		case UINT:
			res = l.(uint) << toUint(r)
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitshift float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitshift strings")
		}
	case lexer.ROUNDSHIFT_LEFT:
		switch getType(l) {
		case BOOL:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot roundshift bool")
		case INT:
			res = RotateLeftInt(l.(int), toUint(r))
		case UINT:
			res = RotateLeftUint(l.(uint), toUint(r))
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot roundshift float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot roundshift strings")
		}
	case lexer.SHIFT_RIGHT:
		switch getType(l) {
		case BOOL:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitshift bool")
		case INT:
			res = l.(int) >> toUint(r)
		case UINT:
			res = l.(uint) >> toUint(r)
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitshift float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitshift strings")
		}
	case lexer.ROUNDSHIFT_RIGHT:
		switch getType(l) {
		case BOOL:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot roundshift bool")
		case INT:
			res = RotateRightInt(l.(int), toUint(r))
		case UINT:
			res = RotateRightUint(l.(uint), toUint(r))
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot roundshift float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot roundshift strings")
		}
	}

	return
}

func BitwiseEval(b Bitwise) (res interface{}, err error) {
	l, err := evaluate(b.Left)
	if err != nil {
		return
	}

	r, err := evaluate(b.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
	if precedence == UNKNOWN {
		err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "invalid operands")
	}

	switch b.Operator.Type {
	case lexer.AND_BITWISE:
		switch precedence {
		case BOOL:
			res = l.(bool) && r.(bool)
		case UINT:
			res = l.(uint) & r.(uint)
		case INT:
			res = l.(int) & r.(int)
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise strings")
		}
	case lexer.OR_BITWISE:
		switch precedence {
		case BOOL:
			res = l.(bool) || r.(bool)
		case UINT:
			res = l.(uint) | r.(uint)
		case INT:
			res = l.(int) | r.(int)
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise strings")
		}
	case lexer.XOR_BITWISE:
		switch precedence {
		case BOOL:
			res = l.(bool) != r.(bool)
		case UINT:
			res = l.(uint) ^ r.(uint)
		case INT:
			res = l.(int) ^ r.(int)
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise strings")
		}
	case lexer.NAND_BITWISE:
		switch precedence {
		case BOOL:
			res = !(l.(bool) && r.(bool))
		case UINT:
			res = ^(l.(uint) & r.(uint))
		case INT:
			res = ^(l.(int) & r.(int))
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise strings")
		}
	case lexer.NOR_BITWISE:
		switch precedence {
		case BOOL:
			res = !(l.(bool) || r.(bool))
		case UINT:
			res = ^(l.(uint) | r.(uint))
		case INT:
			res = ^(l.(int) | r.(int))
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise strings")
		}
	case lexer.XNOR_BITWISE:
		switch precedence {
		case BOOL:
			res = l.(bool) == r.(bool)
		case UINT:
			res = ^(l.(uint) ^ r.(uint))
		case INT:
			res = ^(l.(int) ^ r.(int))
		case FLOAT:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise float")
		case STRING:
			err = errutils.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, errutils.RUNTIME, "cannot bitwise strings")
		}
	}

	return
}

func TermEval(t Term) (res interface{}, err error) {
	l, err := evaluate(t.Left)
	if err != nil {
		return
	}

	r, err := evaluate(t.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
	if precedence == UNKNOWN {
		err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot convert operands")
	}

	switch t.Operator.Type {
	case lexer.PLUS:
		switch precedence {
		case BOOL:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot sum bool values")
		case CHAR:
			res = string(l.(rune)) + string(r.(rune))
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
	l, err := evaluate(t.Left)
	if err != nil {
		return
	}

	r, err := evaluate(t.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
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
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot divide bool values")
		case UINT:
			if r.(uint) == 0 {
				err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "division by zero")
			} else {
				res = l.(uint) / r.(uint)
			}
		case INT:
			if r.(int) == 0 {
				err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "division by zero")
			} else {
				res = l.(int) / r.(int)
			}
		case FLOAT:
			if r.(float64) == 0 {
				err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "division by zero")
			} else {
				res = l.(float64) / r.(float64)
			}
		case STRING:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot divide strings")
		}
	case lexer.MOD:
		switch precedence {
		case BOOL:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot mod bool values")
		case UINT:
			res = l.(uint) % r.(uint)
		case INT:
			res = l.(int) % r.(int)
		case FLOAT:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot mod float values")
		case STRING:
			err = errutils.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, errutils.RUNTIME, "cannot mod strings")
		}
	}

	return
}

func PowerEval(p Power) (res interface{}, err error) {
	l, err := evaluate(p.Left)
	if err != nil {
		return
	}

	r, err := evaluate(p.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
	if precedence == UNKNOWN {
		err = errutils.Error(p.Operator.Line, p.Operator.Column, p.Operator.Lexeme, errutils.RUNTIME, "cannot implicit convert operands")
	}

	switch p.Operator.Type {
	case lexer.POW:
		switch precedence {
		case BOOL:
			err = errutils.Error(p.Operator.Line, p.Operator.Column, p.Operator.Lexeme, errutils.RUNTIME, "cannot power bool values")
		case UINT:
			res = uint(math.Pow(float64(l.(uint)), float64(r.(uint))))
		case INT:
			res = int(math.Pow(float64(l.(int)), float64(r.(int))))
		case FLOAT:
			res = math.Pow(l.(float64), r.(float64))
		case STRING:
			err = errutils.Error(p.Operator.Line, p.Operator.Column, p.Operator.Lexeme, errutils.RUNTIME, "cannot power string values")
		}
	}

	return
}

func UnaryEval(u Unary) (res any, err error) {
	v, err := evaluate(u.Right)
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

func CastEval(c Cast) (res any, err error) {
	l, err := evaluate(c.Left)
	if err != nil {
		return nil, err
	}

	r, err := evaluate(c.TypeCast)
	if err != nil {
		return nil, err
	}

	switch t := r.(type) {
	case Type:
		switch t.Name.Type {
		case lexer.BOOL:
			switch getType(l) {
			case BOOL:
				res = l
			case INT:
				res = l.(int) != 0
			case UINT:
				res = l.(uint) != 0
			case FLOAT:
				res = l.(float64) != 0
			case STRING:
				res = len(l.(string)) == 0
			default:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot convert type to bool")
			}
		case lexer.INT:
			switch getType(l) {
			case BOOL:
				res = boolToInt(l.(bool))
			case INT:
				res = l
			case UINT:
				res = int(l.(uint))
			case FLOAT:
				res = int(l.(float64))
			case STRING:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot conver string to int")
			default:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot convert type to int")
			}
		case lexer.UINT:
			switch getType(l) {
			case BOOL:
				res = uint(boolToInt(l.(bool)))
			case INT:
				res = uint(l.(int))
			case UINT:
				res = l
			case FLOAT:
				res = uint(l.(float64))
			case STRING:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot conver string to uint")
			default:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot convert type to uint")
			}
		case lexer.FLOAT:
			switch getType(l) {
			case BOOL:
				res = float64(boolToInt(l.(bool)))
			case INT:
				res = float64(l.(int))
			case UINT:
				res = float64(l.(uint))
			case FLOAT:
				res = l
			case STRING:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot conver string to float")
			default:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot convert type to float")
			}
		case lexer.STRING:
			switch getType(l) {
			case BOOL:
				res = utils.Ternary(l.(bool), "true", "false")
			default:
				res = fmt.Sprintf("%v", l)
			}
		case lexer.CHAR:
			switch getType(l) {
			case INT:
				res = rune(l.(int))
			default:
				err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "???????????????")
			}
		default:
			err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "cannot convert to type")
		}
	default:
		err = errutils.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, errutils.RUNTIME, "type of typecast not found")
	}

	return res, err
}