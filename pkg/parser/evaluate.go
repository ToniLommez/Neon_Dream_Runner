package parser

import (
	"fmt"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
)

func (p *Program) Interpret(instructions []Stmt) (any, error) {
	for _, i := range instructions {
		if _, err := p.Main.evaluate(i); err != nil {
			return nil, err
		}
	}

	// p.Main.Values.Debug()

	return nil, nil
}

func (s *Scope) evaluate(instruction any) (any, error) {
	switch i := instruction.(type) {
	case LetStmt:
		return s.LetEval(i)
	case PutStmt:
		return s.PutEval(i)
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
		return nil, nil
	case Elvis:
		return nil, nil
	case Check:
		return nil, nil
	case Cast:
		return s.CastEval(i)
	case Identifier:
		return s.IdentifierEval(i)
	case Literal:
		return i.Value, nil
	case Type:
		return i, nil
	case ArrayLiteral:
		return nil, nil
	case Grouping:
		return s.evaluate(i.Expression)
	case nil:
		return nil, e.Error(0, 0, "", e.RUNTIME, "nil instruction should not be found")
	default:
		return nil, e.Error(0, 0, "", e.RUNTIME, fmt.Sprintf("evaluation not implemented for expression: %v", i))
	}
}
