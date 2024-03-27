package parser

import (
	"fmt"
	"math"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
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
	case rune:
		r = v != 0
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
		err = e.Error(-1, -1, "", e.RUNTIME, "truthy pattern matching not implemented, returning false")
	}

	return
}

func (s *Scope) SequenceEval(x Sequence) (any, error) {
	if _, err := s.evaluate(x.Left); err != nil {
		return nil, err
	}

	return s.evaluate(x.Right)
}

func (s *Scope) AssignEval(a Assign) (res any, err error) {
	v, err := s.evaluate(a.Value)
	if err != nil {
		return
	}

	_, tv, d, err := s.Get(a.Target)
	if d {
		if err != nil {
			return
		}

		op := func(t lexer.Token, o lexer.TokenType) lexer.Token {
			t.Type = o
			return t
		}

		switch a.Operator.Type {
		case lexer.ADD_ASSIGN:
			v, err = s.TermEval(Term{Left: Literal{tv}, Operator: op(a.Operator, lexer.PLUS), Right: Literal{v}})
		case lexer.SUB_ASSIGN:
			v, err = s.TermEval(Term{Left: Literal{tv}, Operator: op(a.Operator, lexer.MINUS), Right: Literal{v}})
		case lexer.MUL_ASSIGN:
			v, err = s.FactorEval(Factor{Left: Literal{tv}, Operator: op(a.Operator, lexer.STAR), Right: Literal{v}})
		case lexer.DIV_ASSIGN:
			v, err = s.FactorEval(Factor{Left: Literal{tv}, Operator: op(a.Operator, lexer.SLASH), Right: Literal{v}})
		case lexer.MOD_ASSIGN:
			v, err = s.FactorEval(Factor{Left: Literal{tv}, Operator: op(a.Operator, lexer.MOD), Right: Literal{v}})
		case lexer.POW_ASSIGN:
			v, err = s.PowerEval(Power{Left: Literal{tv}, Operator: op(a.Operator, lexer.POW), Right: Literal{v}})
		case lexer.BITSHIFT_LEFT_ASSIGN:
			v, err = s.BitshiftEval(Bitshift{Left: Literal{tv}, Operator: op(a.Operator, lexer.SHIFT_LEFT), Right: Literal{v}})
		case lexer.BITSHIFT_RIGHT_ASSIGN:
			v, err = s.BitshiftEval(Bitshift{Left: Literal{tv}, Operator: op(a.Operator, lexer.SHIFT_RIGHT), Right: Literal{v}})
		case lexer.ROUNDSHIFT_LEFT_ASSIGN:
			v, err = s.BitshiftEval(Bitshift{Left: Literal{tv}, Operator: op(a.Operator, lexer.ROUNDSHIFT_LEFT), Right: Literal{v}})
		case lexer.ROUNDSHIFT_RIGHT_ASSIGN:
			v, err = s.BitshiftEval(Bitshift{Left: Literal{tv}, Operator: op(a.Operator, lexer.ROUNDSHIFT_RIGHT), Right: Literal{v}})
		case lexer.AND_ASSIGN:
			v, err = s.BitwiseEval(Bitwise{Left: Literal{tv}, Operator: op(a.Operator, lexer.AND_BITWISE), Right: Literal{v}})
		case lexer.OR_ASSIGN:
			v, err = s.BitwiseEval(Bitwise{Left: Literal{tv}, Operator: op(a.Operator, lexer.OR_BITWISE), Right: Literal{v}})
		case lexer.XOR_ASSIGN:
			v, err = s.BitwiseEval(Bitwise{Left: Literal{tv}, Operator: op(a.Operator, lexer.XOR_BITWISE), Right: Literal{v}})
		case lexer.NAND_ASSIGN:
			v, err = s.BitwiseEval(Bitwise{Left: Literal{tv}, Operator: op(a.Operator, lexer.NAND_BITWISE), Right: Literal{v}})
		case lexer.NOR_ASSIGN:
			v, err = s.BitwiseEval(Bitwise{Left: Literal{tv}, Operator: op(a.Operator, lexer.NOR_BITWISE), Right: Literal{v}})
		case lexer.XNOR_ASSIGN:
			v, err = s.BitwiseEval(Bitwise{Left: Literal{tv}, Operator: op(a.Operator, lexer.XNOR_BITWISE), Right: Literal{v}})
		}

		if err != nil {
			return nil, err
		}
	}

	return s.Set(a.Target, v)
}

func (s *Scope) TernaryEval(t Ternary) (any, error) {
	test, err := s.evaluate(t.Expression)
	if err != nil {
		return nil, err
	}

	truthy, err := Truthy(test)
	if err != nil {
		return nil, err
	}

	if truthy {
		return s.evaluate(t.True)
	} else {
		return s.evaluate(t.False)
	}
}

func (s *Scope) LogicEval(x Logic) (any, error) {
	var left bool
	var err error
	var tmp any

	if tmp, err = s.evaluate(x.Left); err != nil {
		return nil, err
	}
	if left, err = Truthy(tmp); err != nil {
		return nil, err
	}

	if x.Operator.Type == lexer.AND_LOGIC {
		if !left {
			return false, nil
		}
	} else {
		if left {
			return true, nil
		}
	}

	if tmp, err = s.evaluate(x.Right); err != nil {
		return nil, err
	}

	return Truthy(tmp)
}

func (s *Scope) EqualityEval(eq Equality) (res any, err error) {
	l, err := s.evaluate(eq.Left)
	if err != nil {
		return
	}

	r, err := s.evaluate(eq.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, true)
	if precedence == UNKNOWN {
		err = e.Error(eq.Operator.Line, eq.Operator.Column, eq.Operator.Lexeme, e.RUNTIME, "invalid operands")
	}

	switch eq.Operator.Type {
	case lexer.EQUAL:
		switch precedence {
		case BOOL:
			res = l.(bool) == r.(bool)
		case CHAR:
			res = l.(rune) == r.(rune)
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
		case CHAR:
			res = l.(rune) != r.(rune)
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

func (s *Scope) ComparisonEval(c Comparison) (res any, err error) {
	l, err := s.evaluate(c.Left)
	if err != nil {
		return
	}

	r, err := s.evaluate(c.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, true)
	if precedence == UNKNOWN {
		err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "invalid operands")
	}

	switch c.Operator.Type {
	case lexer.GREATER_EQUAL:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) >= boolToInt(r.(bool))
		case CHAR:
			res = l.(rune) >= r.(rune)
		case UINT:
			res = l.(uint) >= r.(uint)
		case INT:
			res = l.(int) >= r.(int)
		case FLOAT:
			res = l.(float64) >= r.(float64)
		case STRING:
			err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot direct compare strings")
		}
	case lexer.LESS_EQUAL:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) <= boolToInt(r.(bool))
		case CHAR:
			res = l.(rune) <= r.(rune)
		case UINT:
			res = l.(uint) <= r.(uint)
		case INT:
			res = l.(int) <= r.(int)
		case FLOAT:
			res = l.(float64) <= r.(float64)
		case STRING:
			err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot direct compare strings")
		}
	case lexer.GREATER:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) > boolToInt(r.(bool))
		case CHAR:
			res = l.(rune) > r.(rune)
		case UINT:
			res = l.(uint) > r.(uint)
		case INT:
			res = l.(int) > r.(int)
		case FLOAT:
			res = l.(float64) > r.(float64)
		case STRING:
			err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot direct compare strings")
		}
	case lexer.LESS:
		switch precedence {
		case BOOL:
			res = boolToInt(l.(bool)) < boolToInt(r.(bool))
		case CHAR:
			res = l.(rune) < r.(rune)
		case UINT:
			res = l.(uint) < r.(uint)
		case INT:
			res = l.(int) < r.(int)
		case FLOAT:
			res = l.(float64) < r.(float64)
		case STRING:
			err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot direct compare strings")
		}
	}

	return
}

func (s *Scope) BitshiftEval(b Bitshift) (res any, err error) {
	l, err := s.evaluate(b.Left)
	if err != nil {
		return
	}

	r, err := s.evaluate(b.Right)
	if err != nil {
		return
	}

	if _, _, precedence := typePrecedence(l, r, true); precedence == UNKNOWN {
		err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "invalid operands")
	}

	switch b.Operator.Type {
	case lexer.SHIFT_LEFT:
		switch getType(l) {
		case BOOL:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitshift bool")
		case INT:
			res = l.(int) << toUint(r)
		case UINT:
			res = l.(uint) << toUint(r)
		case FLOAT:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitshift float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitshift strings")
		}
	case lexer.ROUNDSHIFT_LEFT:
		switch getType(l) {
		case BOOL:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot roundshift bool")
		case INT:
			res = RotateLeftInt(l.(int), toUint(r))
		case UINT:
			res = RotateLeftUint(l.(uint), toUint(r))
		case FLOAT:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot roundshift float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot roundshift strings")
		}
	case lexer.SHIFT_RIGHT:
		switch getType(l) {
		case BOOL:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitshift bool")
		case INT:
			res = l.(int) >> toUint(r)
		case UINT:
			res = l.(uint) >> toUint(r)
		case FLOAT:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitshift float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitshift strings")
		}
	case lexer.ROUNDSHIFT_RIGHT:
		switch getType(l) {
		case BOOL:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot roundshift bool")
		case INT:
			res = RotateRightInt(l.(int), toUint(r))
		case UINT:
			res = RotateRightUint(l.(uint), toUint(r))
		case FLOAT:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot roundshift float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot roundshift strings")
		}
	}

	return
}

func (s *Scope) BitwiseEval(b Bitwise) (res any, err error) {
	l, err := s.evaluate(b.Left)
	if err != nil {
		return
	}

	r, err := s.evaluate(b.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
	if precedence == UNKNOWN {
		err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "invalid operands")
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
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise strings")
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
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise strings")
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
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise strings")
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
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise strings")
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
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise strings")
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
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise float")
		case STRING:
			err = e.Error(b.Operator.Line, b.Operator.Column, b.Operator.Lexeme, e.RUNTIME, "cannot bitwise strings")
		}
	}

	return
}

func (s *Scope) TermEval(t Term) (res any, err error) {
	l, err := s.evaluate(t.Left)
	if err != nil {
		return
	}

	r, err := s.evaluate(t.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
	if precedence == UNKNOWN {
		err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot convert operands")
	}

	switch t.Operator.Type {
	case lexer.PLUS:
		switch precedence {
		case BOOL:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot sum bool values")
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
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot sum bool values")
		case UINT:
			res = l.(uint) - r.(uint)
		case INT:
			res = l.(int) - r.(int)
		case FLOAT:
			res = l.(float64) - r.(float64)
		case STRING:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot subtract strings")
		}
	}

	return
}

func (s *Scope) FactorEval(t Factor) (res any, err error) {
	l, err := s.evaluate(t.Left)
	if err != nil {
		return
	}

	r, err := s.evaluate(t.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
	if precedence == UNKNOWN {
		err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot implicit convert operands")
	}

	switch t.Operator.Type {
	case lexer.STAR:
		switch precedence {
		case BOOL:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot multiply bool values")
		case UINT:
			res = l.(uint) * r.(uint)
		case INT:
			res = l.(int) * r.(int)
		case FLOAT:
			res = l.(float64) * r.(float64)
		case STRING:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot multiply string values")
		}
	case lexer.SLASH:
		switch precedence {
		case BOOL:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot divide bool values")
		case UINT:
			if r.(uint) == 0 {
				err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "division by zero")
			} else {
				res = l.(uint) / r.(uint)
			}
		case INT:
			if r.(int) == 0 {
				err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "division by zero")
			} else {
				res = l.(int) / r.(int)
			}
		case FLOAT:
			if r.(float64) == 0 {
				err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "division by zero")
			} else {
				res = l.(float64) / r.(float64)
			}
		case STRING:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot divide strings")
		}
	case lexer.MOD:
		switch precedence {
		case BOOL:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot mod bool values")
		case UINT:
			res = l.(uint) % r.(uint)
		case INT:
			res = l.(int) % r.(int)
		case FLOAT:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot mod float values")
		case STRING:
			err = e.Error(t.Operator.Line, t.Operator.Column, t.Operator.Lexeme, e.RUNTIME, "cannot mod strings")
		}
	}

	return
}

func (s *Scope) PowerEval(p Power) (res any, err error) {
	l, err := s.evaluate(p.Left)
	if err != nil {
		return
	}

	r, err := s.evaluate(p.Right)
	if err != nil {
		return
	}

	l, r, precedence := typePrecedence(l, r, false)
	if precedence == UNKNOWN {
		err = e.Error(p.Operator.Line, p.Operator.Column, p.Operator.Lexeme, e.RUNTIME, "cannot implicit convert operands")
	}

	switch p.Operator.Type {
	case lexer.POW:
		switch precedence {
		case BOOL:
			err = e.Error(p.Operator.Line, p.Operator.Column, p.Operator.Lexeme, e.RUNTIME, "cannot power bool values")
		case UINT:
			res = uint(math.Pow(float64(l.(uint)), float64(r.(uint))))
		case INT:
			res = int(math.Pow(float64(l.(int)), float64(r.(int))))
		case FLOAT:
			res = math.Pow(l.(float64), r.(float64))
		case STRING:
			err = e.Error(p.Operator.Line, p.Operator.Column, p.Operator.Lexeme, e.RUNTIME, "cannot power string values")
		}
	}

	return
}

func (s *Scope) UnaryEval(u Unary) (res any, err error) {
	v, err := s.evaluate(u.Right)
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
			err = e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, fmt.Sprintf("expect number after !, received string: \"%v\"", v))
		default:
			err = e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, fmt.Sprintf("expect number after !, received: %v", v))
		}
	case lexer.NOT_BITWISE:
		switch t {
		case INT:
			res = ^v.(int)
		case BOOL:
			res = !v.(bool)
		case FLOAT:
			err = e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, "operator ~ not defined on float")
		default:
			return nil, e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, "expect number after ~")
		}
	case lexer.PLUS:
		res = v
	case lexer.MINUS:
		switch t {
		case INT:
			res = -v.(int)
		case BOOL:
			err = e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, "operator - not defined on bool")
		case FLOAT:
			res = -v.(float64)
		default:
			return nil, e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, "expect number after ~")
		}
	case lexer.GO_IN:
		err = e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, "parallelism not yet implemented")
	default:
		err = e.Error(o.Line, o.Column, o.Lexeme, e.RUNTIME, "unknown operator")
	}

	return res, err
}

func (s *Scope) CastEval(c Cast) (res any, err error) {
	l, err := s.evaluate(c.Left)
	if err != nil {
		return nil, err
	}

	r, err := s.evaluate(c.TypeCast)
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
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot convert type to bool")
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
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot conver string to int")
			default:
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot convert type to int")
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
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot conver string to uint")
			default:
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot convert type to uint")
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
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot conver string to float")
			default:
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot convert type to float")
			}
		case lexer.STRING:
			switch getType(l) {
			case BOOL:
				res = utils.Ternary(l.(bool), "true", "false")
			case CHAR:
				res = string(l.(rune))
			default:
				res = fmt.Sprintf("%v", l)
			}
		case lexer.CHAR:
			switch getType(l) {
			case INT:
				res = rune(l.(int))
			case STRING:
				res = rune(l.(string)[0])
			default:
				err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "???????????????")
			}
		default:
			err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "cannot convert to type")
		}
	default:
		err = e.Error(c.Operator.Line, c.Operator.Column, c.Operator.Lexeme, e.RUNTIME, "type of typecast not found")
	}

	return res, err
}

func (s *Scope) IdentifierEval(i Identifier) (res any, err error) {
	_, res, _, err = s.Get(i.Name)
	if err != nil {
		return nil, err
	}

	return
}
