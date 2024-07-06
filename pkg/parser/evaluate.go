package parser

import (
	"fmt"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
)

func (s *Scope) Interpret() (any, error) {
	var v any
	var err error
	for _, stmt := range s.Statements {
		if v, err = s.evaluate(stmt); err != nil {
			return nil, err
		}
	}

	// p.Main.Values.Debug()

	return v, nil
}

func (s *Scope) evaluate(instruction any) (any, error) {
	switch i := instruction.(type) {
	case FnStmt:
		return s.FnEval(i)
	case LetStmt:
		return s.LetEval(i)
	case IfStmt:
		return s.IfEval(i)
	case PutStmt:
		return s.PutEval(i)
	case WhileStmt:
		return s.WhileEval(i)
	case ExprStmt:
		return s.ExprEval(i)
	case Sequence:
		return s.SequenceEval(i)
	case Assign:
		return s.AssignEval(i)
	case Pipeline:
		return nil, nil
	case Ternary:
		return s.TernaryEval(i)
	case Range:
		return nil, nil
	case Logic:
		return s.LogicEval(i)
	case Equality:
		return s.EqualityEval(i)
	case Comparison:
		return s.ComparisonEval(i)
	case Bitshift:
		return s.BitshiftEval(i)
	case Bitwise:
		return s.BitwiseEval(i)
	case Term:
		return s.TermEval(i)
	case Factor:
		return s.FactorEval(i)
	case Power:
		return s.PowerEval(i)
	case Increment:
		return nil, nil
	case Pointer:
		return nil, nil
	case Unary:
		return s.UnaryEval(i)
	case Access:
		return nil, nil
	case PositionAccess:
		return s.PositionAccess(i)
	case Elvis:
		return nil, nil
	case Check:
		return nil, nil
	case Cast:
		return s.CastEval(i)
	case Identifier:
		return s.IdentifierEval(i)
	case Caller:
		return s.CallerEval(i)
	case Literal:
		return i.Value, nil
	case Slice:
		return i, nil
	case ArrayConstructor:
		return s.ArrayConstructorEval(i)
	case ArrayType:
		return s.ArrayTypeEval(i)
	case Type:
		return i, nil
	case Grouping:
		return s.evaluate(i.Expression)
	case Block:
		return s.BlockEval(i)
	case nil:
		return nil, e.Error(0, 0, "", e.RUNTIME, "nil instruction should not be found")
	default:
		return nil, e.Error(0, 0, "", e.RUNTIME, fmt.Sprintf("evaluation not implemented for expression: %v", i))
	}
}
