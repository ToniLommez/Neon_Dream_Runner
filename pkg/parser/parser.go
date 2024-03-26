package parser

import (
	"fmt"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Parser struct {
	Tokens  []l.Token
	Current int
	Depth   int
}

func (p *Parser) Parse() ([]Stmt, error) {
	stmt := make([]Stmt, 0)

	for !p.isLastToken() && !p.isAtEnd() {
		if p.peek().Type == l.NEW_LINE {
			p.consume(l.NEW_LINE)
		}

		s, err := p.declaration()
		if err != nil {
			return stmt, err
		}

		stmt = append(stmt, s)
	}

	return stmt, nil
}

func (p *Parser) declaration() (Stmt, error) {
	if p.match(l.LET) {
		s, err := p.letStatement()

		if err != nil {
			p.Synchronize()
			return s, err
		}

		return s, err
	}

	return p.statement()
}

func (p *Parser) letStatement() (Stmt, error) {
	var mutable, nullable bool
	var initializer Expr
	var name, varType l.Token
	var err error

	if p.match(l.BANG) {
		mutable = true
	}
	if p.match(l.CHECK) {
		nullable = true
	}

	name, err = p.consume(l.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	varType = l.Token{Type: l.UNDEFINED, Lexeme: "", Literal: "", Line: name.Line, Column: name.Column}
	if p.match(l.COLON) {
		varType = p.advance()
		if !varType.Type.IsValidType() {
			return nil, e.Error(varType.Line, varType.Column, varType.Lexeme, e.PARSER, "expect type in let statement")
		}
	}

	if p.match(l.ASSIGN) {
		initializer, err = p.expression()
		if err != nil {
			return initializer, err
		}
	}

	// Bad smell, but it's working... So, it's a problem for the future when it broke the whole thing
	if p.peek().Type != l.RIGHT_BRACE {
		if _, err := p.consume(l.NEW_LINE); err != nil {
			t := p.peek()
			return nil, e.Error(t.Line, t.Column, t.Lexeme, e.PARSER, "expect new line after let statement")
		}
	}

	return LetStmt{Name: name, Mutable: mutable, Nullable: nullable, Type: tokenToType(varType), Initializer: initializer}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(l.PUT) {
		return p.putStatement()
	}

	return p.expressionStatement()
}

func (p *Parser) ifStatement() (Expr, error) {
	var condition, thenBranch, elseBranch Expr
	var err error

	condition, err = p.expression()
	if err != nil {
		return condition, err
	}

	if err := p.ensureNotUnterminated(); err != nil {
		return nil, err
	}

	// Should we add an optional p.consule(l.THEN) here?
	// If so, just create the token in the lexer package...

	if thenBranch, err = p.block(true); err != nil {
		return nil, err
	}

	// the word if is consumed before the ifStatement is called,
	// so, we consume the elif here, and then call recursively
	if p.match(l.ELIF) {
		elseBranch, err = p.ifStatement()
	} else if p.match(l.ELSE) {
		elseBranch, err = p.block(true)
	}
	if err != nil {
		return nil, err
	}

	return IfStmt{Condition: condition, Then: thenBranch, Else: elseBranch}, nil
}

func (p *Parser) putStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return expr, err
	}

	t, err := p.consume(l.NEW_LINE)
	if err != nil {
		t = p.peek()
		return nil, e.Error(t.Line, t.Column, t.Lexeme, e.PARSER, "expect new line after print")
	}

	return PutStmt{Value: expr}, nil
}

// A state that contains an expression
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return expr, err
	}

	if p.peek().Type != l.RIGHT_BRACE {
		t, err := p.consume(l.NEW_LINE)
		if err != nil {
			t = p.peek()
			return nil, e.Error(t.Line, t.Column, t.Lexeme, e.PARSER, "expect new line before new expression")
		}
	}

	return ExprStmt{Expr: expr}, nil
}

func (p *Parser) expression() (Expr, error) {
	return p.statementExpression()
}

// Statements that ARE expressions
func (p *Parser) statementExpression() (Expr, error) {
	if p.match(l.IF) {
		return p.ifStatement()
	}

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
		op := p.previous()
		right, err := p.expression()
		if err != nil {
			return expr, err
		}

		switch i := expr.(type) {
		case Identifier:
			expr = Assign{Target: i.Name, Operator: op, Value: right}
		default:
			return expr, e.Error(op.Line, op.Column, op.Lexeme, e.PARSER, "assignment target should be a identifier")
		}
	}

	return expr, nil
}

func (p *Parser) pipeline() (Expr, error) {
	expr, err := p.ternary()
	if err != nil {
		return expr, err
	}

	for p.match(l.PIPELINE_LEFT, l.PIPELINE_RIGHT) {
		op := p.previous()

		var right Expr
		if op.Type == l.PIPELINE_LEFT {
			right, err = p.expression()
		} else {
			right, err = p.ternary()
		}

		if err != nil {
			return expr, err
		}
		expr = Pipeline{Left: expr, Operator: op, Right: right}
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
		op := p.previous()
		right, err := p.equality()
		if err != nil {
			return expr, err
		}

		expr = Logic{Left: expr, Operator: op, Right: right}
	}

	return expr, nil
}

func (p *Parser) equality() (Expr, error) {
	expr, err := p.comparison()
	if err != nil {
		return expr, err
	}

	for p.match(l.NOT_EQUAL, l.EQUAL) {
		op := p.previous()
		right, err := p.comparison()
		expr = Equality{expr, op, right}
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
		op := p.previous()
		right, err := p.bitshift()
		expr = Comparison{expr, op, right}
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
		op := p.previous()
		right, err := p.bitwise()
		if err != nil {
			return expr, err
		}

		expr = Bitshift{expr, op, right}
	}

	return expr, nil
}

func (p *Parser) bitwise() (Expr, error) {
	expr, err := p.term()
	if err != nil {
		return expr, err
	}

	for p.match(l.AND_BITWISE, l.OR_BITWISE, l.XOR_BITWISE, l.NAND_BITWISE, l.NOR_BITWISE, l.XNOR_BITWISE) {
		op := p.previous()
		right, err := p.term()
		if err != nil {
			return expr, err
		}

		expr = Bitwise{expr, op, right}
	}

	return expr, nil
}

func (p *Parser) term() (Expr, error) {
	expr, err := p.factor()
	if err != nil {
		return expr, err
	}

	for p.match(l.PLUS, l.MINUS) {
		op := p.previous()
		right, err := p.factor()
		expr = Term{expr, op, right}
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
		op := p.previous()
		right, err := p.power()
		expr = Factor{expr, op, right}
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
		op := p.previous()
		right, err := p.increment()
		if err != nil {
			return expr, err
		}

		expr = Power{Left: expr, Right: right, Operator: op}
	}

	return expr, nil
}

func (p *Parser) increment() (Expr, error) {
	if p.match(l.INCREMENT, l.DECREMENT) {
		op := p.previous()
		right, err := p.expression()
		return Increment{Expression: right, Operator: op, Position: false}, err
	}

	expr, err := p.pointer()
	if err != nil {
		return expr, err
	}

	if p.match(l.INCREMENT, l.DECREMENT) {
		op := p.previous()
		expr = Increment{Expression: expr, Operator: op, Position: true}
	}

	return expr, nil
}

func (p *Parser) pointer() (Expr, error) {
	if p.match(l.STAR, l.AND_BITWISE) {
		op := p.previous()
		right, err := p.pointer()
		return Pointer{op, right}, err
	}

	return p.unary()
}

func (p *Parser) unary() (Expr, error) {
	if p.match(l.BANG, l.NOT_BITWISE, l.PLUS, l.MINUS, l.GO_IN) {
		op := p.previous()
		right, err := p.unary()
		return Unary{op, right}, err
	}

	return p.access()
}

func (p *Parser) access() (Expr, error) {
	expr, err := p.validate()
	if err != nil {
		return expr, err
	}

	for p.match(l.CHECK_NAV, l.BANG_NAV, l.DOT, l.LEFT_BRACKET) {
		op := p.previous()

		if op.Type == l.LEFT_BRACKET {
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

			expr = Access{Left: expr, Right: right, Operator: op}
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
		op := p.previous()
		actual := p.Current
		right, err := p.primary()
		if err != nil {
			return expr, err
		}

		switch t := right.(type) {
		case Type:
			expr = Cast{Left: expr, TypeCast: t, Operator: op}
		default:
			p.Current = actual - 1
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

	return p.mapLiteral()
}

func (p *Parser) mapLiteral() (Expr, error) {
	current := p.peek()
	if current.Type == l.LEFT_BRACKET {
		if _, next := p.peekN(1); next.Type.IsType() {
			return p.arrayLiteral()
		}
	} /* else if x.Type == l.LEFT_PAREN {

	} else if x.Type == l.OR_BITWISE {

	} */
	return p.group()
}

func (p *Parser) arrayLiteral() (Expr, error) {
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

	return p.block(false)
}

func (p *Parser) block(isRequired bool) (Expr, error) {
	if p.match(l.LEFT_BRACE) {
		p.Depth++
		if err := p.ensureNotUnterminated(); err != nil {
			return nil, err
		}

		statements := make([]Stmt, 0)
		for p.hasMore() && !p.check(l.RIGHT_BRACE) {
			s, err := p.declaration()
			if err != nil {
				return nil, err
			}

			if err := p.ensureNotUnterminated(); err != nil {
				return nil, err
			}

			statements = append(statements, s)
		}

		if _, err := p.consume(l.RIGHT_BRACE); err != nil {
			return nil, err
		}
		p.Depth--

		var scope Scope
		scope.Init()
		scope.Statements = statements
		return Block{Scope: scope}, nil
	} else if isRequired {
		token := p.Tokens[p.Current]
		return nil, e.Error(token.Line, token.Column, token.Lexeme, e.PARSER, fmt.Sprintf("expected a block statement, found: %v", token))
	}

	return p.deadEnd()
}

func (p *Parser) deadEnd() (Expr, error) {
	token := p.Tokens[p.Current]
	return nil, e.Error(token.Line, token.Column, token.Lexeme, e.PARSER, fmt.Sprintf("expect expression, found: %v", token))
}
