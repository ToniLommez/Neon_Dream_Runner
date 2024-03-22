package parser

import (
	"fmt"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
)

func Interpret(instructions []Stmt) (any, error) {
	for _, i := range instructions {
		if _, err := evaluate(i); err != nil {
			return nil, err
		}

	}

	return nil, nil
}

func evaluate(instruction any) (any, error) {
	switch i := instruction.(type) {
	case ExprStmt:
		return ExprStmtEval(i)
	case PutStmt:
		return PutStmtEval(i)
	case Sequence:
		return SequenceEval(i)
	case Assign:
		return nil, nil
	case Pipeline:
		return nil, nil
	case Ternary:
		return TernaryEval(i)
	case Range:
		return nil, nil
	case Logic:
		return LogicEval(i)
	case Equality:
		return EqualityEval(i)
	case Comparison:
		return ComparisonEval(i)
	case Bitshift:
		return BitshiftEval(i)
	case Bitwise:
		return BitwiseEval(i)
	case Term:
		return TermEval(i)
	case Factor:
		return FactorEval(i)
	case Power:
		return PowerEval(i)
	case Increment:
		return nil, nil
	case Pointer:
		return nil, nil
	case Unary:
		return UnaryEval(i)
	case Access:
		return nil, nil
	case PositionAccess:
		return nil, nil
	case Elvis:
		return nil, nil
	case Check:
		return nil, nil
	case Cast:
		return CastEval(i)
	case Identifier:
		return nil, nil
	case Literal:
		return i.Value, nil
	case Type:
		return i, nil
	case ArrayLiteral:
		return nil, nil
	case Grouping:
		return evaluate(i.Expression)
	case nil:
		return nil, e.Error(0, 0, "", e.RUNTIME, "nil instruction should not be found")
	default:
		return nil, e.Error(0, 0, "", e.RUNTIME, fmt.Sprintf("evaluation not implemented for expression: %v", i))
	}
}
