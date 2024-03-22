package parser

import (
	"fmt"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Parser struct {
	tokens  []l.Token
	current int
}

func (p *Parser) Parse() ([]Stmt, error) {
	stmt := make([]Stmt, 0)

	for !p.isLastToken() && !p.isAtEnd() {
		s, err := p.statement()
		if err != nil {
			return stmt, err
		}

		stmt = append(stmt, s)
	}

	/* expr, err := p.expression()
	if len(p.tokens) != p.current+1 {
		token := p.tokens[p.current]
		err = errutils.Error(token.Line, token.Column, token.Lexeme, errutils.PARSER, "unexpected value found")
	}

	if err != nil {
		fmt.Println(err)
		return nil
	} */

	return stmt, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(l.PUT) {
		return p.putStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) putStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return expr, err
	}

	t, err := p.consume(l.NEW_LINE)
	if err != nil {
		t = p.peek()
		return nil, errutils.Error(t.Line, t.Column, t.Lexeme, errutils.RUNTIME, "expect new line after print")
	}

	return PutStmt{Value: expr}, nil
}

func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return expr, err
	}

	t, err := p.consume(l.NEW_LINE)
	if err != nil {
		t = p.peek()
		return nil, errutils.Error(t.Line, t.Column, t.Lexeme, errutils.RUNTIME, "expect new line before new expression")
	}

	return ExprStmt{Expr: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.sequence()
}

func (p *Parser) sequence() (Expr, error) {
	expr, err := p.assign()
	if err != nil {
		return expr, err
	}

	for p.match(l.SEMICOLON) {
		right, err := p.assign()
		if err != nil {
			return expr, err
		}
		expr = Sequence{expr, right}
	}

	return expr, nil
}

func (p *Parser) assign() (Expr, error) {
	expr, err := p.pipeline()
	if err != nil {
		return expr, err
	}

	for p.match(l.ASSIGN, l.ADD_ASSIGN, l.SUB_ASSIGN, l.MUL_ASSIGN, l.DIV_ASSIGN, l.MOD_ASSIGN, l.POW_ASSIGN, l.BITSHIFT_LEFT_ASSIGN, l.BITSHIFT_RIGHT_ASSIGN, l.ROUNDSHIFT_LEFT_ASSIGN, l.ROUNDSHIFT_RIGHT_ASSIGN, l.AND_ASSIGN, l.OR_ASSIGN, l.XOR_ASSIGN, l.NAND_ASSIGN, l.NOR_ASSIGN, l.XNOR_ASSIGN) {
		operator := p.previous()
		right, err := p.expression()
		if err != nil {
			return expr, err
		}

		expr = Assign{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) pipeline() (Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return expr, err
	}

	for p.match(l.PIPELINE_LEFT, l.PIPELINE_RIGHT) {
		operator := p.previous()

		var right Expr
		if operator.Type == l.PIPELINE_LEFT {
			right, err = p.expression()
		} else {
			right, err = p.ternary()
		}

		if err != nil {
			return expr, err
		}
		expr = Pipeline{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) ternary() (Expr, error) {
	expr, err := p.interval()
	if err != nil {
		return expr, err
	}

	for p.match(l.CHECK) {
		trueExpr, err := p.expression()
		if err != nil {
			return expr, err
		}

		_, err = p.consume(l.COLON)
		if err != nil {
			return expr, err
		}

		falseExpr, err := p.expression()
		if err != nil {
			return expr, err
		}

		expr = Ternary{Expression: expr, True: trueExpr, False: falseExpr}
	}

	return expr, nil
}

func (p *Parser) interval() (Expr, error) {
	expr, err := p.logic()
	if err != nil {
		return expr, err
	}

	if p.match(l.RANGE_DOT) {
		right, err := p.expression()
		if err != nil {
			return expr, err
		}

		expr = Range{expr, right}
	}

	return expr, nil
}

func (p *Parser) logic() (Expr, error) {
	expr, err := p.equality()
	if err != nil {
		return expr, err
	}

	for p.match(l.AND_LOGIC, l.OR_LOGIC) {
		operator := p.previous()
		right, err := p.equality()
		if err != nil {
			return expr, err
		}

		expr = Logic{Left: expr, Operator: operator, Right: right}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return expr, err
	}

	for p.match(l.NOT_EQUAL, l.EQUAL) {
		operator := p.previous()
		right, err := p.comparison()
		expr = Equality{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) comparison() (Expr, error) {
	expr, err := p.bitshift()
	if err != nil {
		return expr, err
	}

	for p.match(l.GREATER, l.GREATER_EQUAL, l.LESS, l.LESS_EQUAL) {
		operator := p.previous()
		right, err := p.bitshift()
		expr = Comparison{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) bitshift() (Expr, error) {
	expr, err := p.bitwise()
	if err != nil {
		return expr, err
	}

	for p.match(l.SHIFT_LEFT, l.ROUNDSHIFT_LEFT, l.SHIFT_RIGHT, l.ROUNDSHIFT_RIGHT) {
		operator := p.previous()
		right, err := p.bitwise()
		if err != nil {
			return expr, err
		}

		expr = Bitshift{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) bitwise() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return expr, err
	}

	for p.match(l.AND_BITWISE, l.OR_BITWISE, l.XOR_BITWISE, l.NAND_BITWISE, l.NOR_BITWISE, l.XNOR_BITWISE) {
		operator := p.previous()
		right, err := p.term()
		if err != nil {
			return expr, err
		}

		expr = Bitwise{expr, operator, right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return expr, err
	}

	for p.match(l.PLUS, l.MINUS) {
		operator := p.previous()
		right, err := p.factor()
		expr = Term{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) factor() (Expr, error) {
	expr, err := p.power()
	if err != nil {
		return expr, err
	}

	for p.match(l.STAR, l.SLASH, l.MOD) {
		operator := p.previous()
		right, err := p.power()
		expr = Factor{expr, operator, right}
		if err != nil {
			return expr, err
		}
	}

	return expr, nil
}

func (p *Parser) power() (Expr, error) {
	expr, err := p.increment()
	if err != nil {
		return expr, err
	}

	for p.match(l.POW) {
		operator := p.previous()
		right, err := p.increment()
		if err != nil {
			return expr, err
		}

		expr = Power{Left: expr, Right: right, Operator: operator}
	}

	return expr, nil
}

func (p *Parser) increment() (Expr, error) {
	if p.match(l.INCREMENT, l.DECREMENT) {
		operator := p.previous()
		right, err := p.expression()
		return Increment{Expression: right, Operator: operator, Position: false}, err
	}

	expr, err := p.pointer()
	if err != nil {
		return expr, err
	}

	if p.match(l.INCREMENT, l.DECREMENT) {
		operator := p.previous()
		expr = Increment{Expression: expr, Operator: operator, Position: true}
	}

	return expr, nil
}

func (p *Parser) pointer() (Expr, error) {
	if p.match(l.STAR, l.AND_BITWISE) {
		operator := p.previous()
		right, err := p.pointer()
		return Pointer{operator, right}, err
	}

	return p.unary()
}

func (p *Parser) unary() (Expr, error) {
	if p.match(l.BANG, l.NOT_BITWISE, l.PLUS, l.MINUS, l.GO_IN) {
		operator := p.previous()
		right, err := p.unary()
		return Unary{operator, right}, err
	}

	return p.access()
}

func (p *Parser) access() (Expr, error) {
	expr, err := p.validate()
	if err != nil {
		return expr, err
	}

	for p.match(l.CHECK_NAV, l.BANG_NAV, l.DOT, l.LEFT_BRACKET) {
		operator := p.previous()

		if operator.Type == l.LEFT_BRACKET {
			right, err := p.expression()
			if err != nil {
				return expr, err
			}

			expr = PositionAccess{Expression: expr, Pos: right}

			if _, err := p.consume(l.RIGHT_BRACKET); err != nil {
				return expr, err
			}
		} else {
			right, err := p.validate()
			if err != nil {
				return expr, err
			}

			expr = Access{Left: expr, Right: right, Operator: operator}
		}

	}

	return expr, nil
}

func (p *Parser) validate() (Expr, error) {
	expr, err := p.catch()
	if err != nil {
		return expr, err
	}

	if p.match(l.ELVIS) {
		x := p.peek()
		if x.Type == l.NEW_LINE || x.Type == l.RIGHT_BRACE || x.Type == l.RIGHT_PAREN || x.Type == l.RIGHT_BRACKET {
			expr = Elvis{Left: expr, ReturnZero: true, Right: nil}
		} else {
			right, err := p.catch()
			if err != nil {
				return expr, err
			}

			expr = Elvis{Left: expr, ReturnZero: false, Right: right}
		}
	}

	return expr, nil
}

func (p *Parser) catch() (Expr, error) {
	expr, err := p.cast()
	if err != nil {
		return expr, err
	}

	tmp := p.peek()
	if tmp.Type == l.CHECK {
		found, x := p.peekN(1)
		if !found || x.Type == l.NEW_LINE {
			p.advance()
			expr = Check{Left: expr, HaveReturn: false, Right: nil}
		} else {
			found, t := p.peekN(1)
			if found && t.Type == l.RETURN {
				p.advance()
				p.advance()
				right, err := p.expression()
				if err != nil {
					return expr, err
				}
				expr = Check{Left: expr, HaveReturn: true, Right: right}
			}
		}
	}
	return expr, nil
}

func (p *Parser) cast() (Expr, error) {
	expr, err := p.primary()
	if err != nil {
		return expr, err
	}

	for p.match(l.COLON) {
		operator := p.previous()
		actual := p.current
		right, err := p.primary()
		if err != nil {
			return expr, err
		}

		switch t := right.(type) {
		case Type:
			expr = Cast{Left: expr, TypeCast: t, Operator: operator}
		default:
			p.current = actual - 1
			return expr, nil
		}
	}

	return expr, nil
}

func (p *Parser) primary() (Expr, error) {
	if p.match(l.IDENTIFIER) {
		return Identifier{p.previous()}, nil
	}
	if p.match(l.STRING_LITERAL, l.NUMBER_LITERAL, l.FLOAT_LITERAL) {
		return Literal{p.previous().Literal}, nil
	}
	if p.match(l.TRUE) {
		return Literal{true}, nil
	}
	if p.match(l.FALSE) {
		return Literal{false}, nil
	}
	if p.match(l.NIL) {
		return Literal{nil}, nil
	}
	if p.match(l.INT, l.I8, l.I16, l.I32, l.I64, l.UINT, l.U8, l.U16, l.U32, l.U64, l.FLOAT, l.F32, l.F64, l.BOOL, l.CHAR, l.STRING, l.BYTE, l.ANY) {
		return Type{Name: p.previous()}, nil
	}

	return p.declaration()
}

func (p *Parser) declaration() (Expr, error) {
	current := p.peek()
	if current.Type == l.LEFT_BRACKET {
		if _, next := p.peekN(1); next.Type.IsType() {
			return p.arrayDeclaration()
		}
	} /* else if x.Type == l.LEFT_PAREN {

	} else if x.Type == l.OR_BITWISE {

	} */
	return p.group()
}

func (p *Parser) arrayDeclaration() (Expr, error) {
	p.advance()
	arrayType := p.advance()
	if _, err := p.consume(l.COLON); err != nil {
		return nil, err
	}

	arraySize, err := p.expression()
	if err != nil {
		return nil, err
	}

	if _, err := p.consume(l.RIGHT_BRACKET); err != nil {
		return nil, err
	}

	if _, err := p.consume(l.LEFT_BRACKET); err != nil {
		return nil, err
	}

	if p.check(l.RIGHT_BRACKET) {
		return ArrayLiteral{Typing: arrayType, Size: arraySize, Values: []Expr{}}, nil
	} else {
		var values []Expr

		for {
			value, err := p.expression()
			if err != nil {
				return nil, err
			}
			values = append(values, value)
			if !p.match(l.COMMA) {
				break
			}
		}

		if !p.match(l.RIGHT_BRACKET) {
			return nil, err
		}

		return ArrayLiteral{Typing: arrayType, Size: arraySize, Values: values}, nil
	}
}

func (p *Parser) group() (Expr, error) {
	if p.match(l.LEFT_PAREN) {
		expr, err := p.expression()
		if err != nil {
			return expr, err
		}
		if _, err := p.consume(l.RIGHT_PAREN); err != nil {
			return expr, err
		}
		return Grouping{expr}, nil
	}

	return p.block()
}

func (p *Parser) block() (Expr, error) {
	if p.match(l.LEFT_BRACE) {
		/* 		if p.isLastToken() {
			final := p.previous()
			return nil, errutils.Error(final.Line, final.Column, final.Lexeme, errutils.UNTERMINATED_STATEMENT, "unterminated statement")
		} */

		expr, err := p.expression()
		if err != nil {
			return expr, err
		}

		/* if p.isLastToken() {
			final := p.previous()
			return nil, errutils.Error(final.Line, final.Column, final.Lexeme, errutils.UNTERMINATED_STATEMENT, "unterminated statement")
		} */
		if _, err := p.consume(l.RIGHT_BRACE); err != nil {
			return expr, err
		}

		return Grouping{expr}, nil
	}

	return p.deadEnd()
}

func (p *Parser) deadEnd() (Expr, error) {
	token := p.tokens[p.current]
	return nil, errutils.Error(token.Line, token.Column, token.Lexeme, errutils.PARSER, fmt.Sprintf("expect expression, found: %v", token))
}
