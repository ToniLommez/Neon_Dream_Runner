package parser

import (
	"fmt"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/utils"
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

	if p.match(l.FN) {
		s, err := p.function()

		if err != nil {
			// TODO: change to a synchronize proper to function
			p.Synchronize()
			return s, err
		}

		return s, err
	}

	return p.statement()
}

func (p *Parser) function() (Stmt, error) {
	var name l.Token
	var args []LetStmt
	var rets []int
	var err error
	var block Expr

	// Function name
	name, err = p.consume(l.IDENTIFIER)
	if err != nil {
		return nil, err
	}

	// Args
	if p.match(l.LEFT_PAREN) {
		for {
			var varName, varType l.Token

			if varName, err = p.consume(l.IDENTIFIER); err != nil {
				return nil, err
			}

			if _, err = p.consume(l.COLON); err != nil {
				return nil, err
			}

			varType = p.advance()
			if !varType.Type.IsValidType() {
				return nil, e.Error(varType.Line, varType.Column, varType.Lexeme, e.PARSER, "expect type")
			}

			// TODO: add mutability, nullability and etc
			args = append(args, LetStmt{Name: varName, Mutable: true, Nullable: true, Type: varType, Initializer: nil})

			if !p.match(l.COMMA) {
				break
			}
		}

		if _, err = p.consume(l.RIGHT_PAREN); err != nil {
			return nil, err
		}
	}

	// Return types
	if p.match(l.RETURN) {
		for {
			varType := p.advance()
			if !varType.Type.IsValidType() {
				return nil, e.Error(varType.Line, varType.Column, varType.Lexeme, e.PARSER, "expect type")
			}

			rets = append(rets, tokenToType(varType))

			if !p.match(l.COMMA) {
				break
			}
		}
	}

	if block, err = p.block(true); err != nil {
		return nil, err
	}

	return FnStmt{Name: name, Args: args, Return: rets, Context: block.(Block).Scope}, nil
}

func (p *Parser) letStatement() (Stmt, error) {
	var mutable, nullable bool
	var initializer Expr
	var varType Expr
	var name l.Token
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

	if p.match(l.COLON) {
		if varType, err = p.primary(); err != nil {
			return nil, err
		}
	}

	if p.match(l.ASSIGN) {
		initializer, err = p.expression()
		if err != nil {
			return initializer, err
		}
	}

	// Bad smell, but it's working... So, it's a problem for the future when it broke the whole thing
	if t := p.peek().Type; t != l.RIGHT_BRACE && t != l.SEMICOLON { // Semicolon is because of declaration on for statement
		if _, err := p.consume(l.NEW_LINE); err != nil {
			t := p.peek()
			return nil, e.Error(t.Line, t.Column, t.Lexeme, e.PARSER, "expect new line after let statement")
		}
	}

	return LetStmt{Name: name, Mutable: mutable, Nullable: nullable, Type: varType, Initializer: initializer, IsSlice: false}, nil
}

func (p *Parser) statement() (Stmt, error) {
	if p.match(l.FOR) {
		return p.forStatement()
	} else if p.match(l.PUT) {
		return p.putStatement(false)
	} else if p.match(l.PUTLN) {
		return p.putStatement(true)
	} else if p.match(l.WHILE) {
		return p.whileStatement()
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

func (p *Parser) forStatement() (Stmt, error) {
	var condition, increment, body Expr
	var initializer Stmt
	var err error

	// Declaration
	if p.match(l.LET) {
		initializer, err = p.letStatement()
	} else if p.peek().Type != l.SEMICOLON {
		initializer, err = p.expressionStatement()
	}
	if err != nil {
		return nil, err
	}
	if _, err = p.consume(l.SEMICOLON); err != nil {
		return nil, err
	}

	// Condition
	if !p.check(l.SEMICOLON) {
		if condition, err = p.expression(); err != nil {
			return nil, err
		}
	}
	if _, err = p.consume(l.SEMICOLON); err != nil {
		return nil, err
	}

	// Increment
	if !p.check(l.SEMICOLON) {
		if increment, err = p.expression(); err != nil {
			return nil, err
		}
	}

	if body, err = p.block(true); err != nil {
		return nil, err
	}

	if increment != nil {
		s := (body.(Block))
		s.Scope.Statements = append(s.Scope.Statements, increment)
		body = s
	}

	if condition == nil {
		condition = Literal{true}
	}

	body = WhileStmt{Condition: condition, Body: body}

	if initializer != nil {
		var scope Scope
		scope.Init()
		scope.Statements = []Stmt{initializer, body}
		body = Block{Scope: scope}
	}

	return body, nil
}

func (p *Parser) putStatement(line bool) (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return expr, err
	}

	if p.peek().Type == l.NEW_LINE {
		p.advance()
	}

	/* t, err := p.consume(l.NEW_LINE)
	if err != nil {
		t = p.peek()
		return nil, e.Error(t.Line, t.Column, t.Lexeme, e.PARSER, "expect new line after print")
	} */

	return PutStmt{Value: expr, NewLine: line}, nil
}

func (p *Parser) whileStatement() (Stmt, error) {
	var expr, body Expr
	var err error

	if expr, err = p.expression(); err != nil {
		return nil, err
	}
	if err = p.ensureNotUnterminated(); err != nil {
		return nil, err
	}

	if body, err = p.block(true); err != nil {
		return nil, err
	}

	return WhileStmt{Condition: expr, Body: body}, nil
}

// A state that contains an expression
func (p *Parser) expressionStatement() (Stmt, error) {
	expr, err := p.expression()
	if err != nil {
		return expr, err
	}

	if t := p.peek().Type; t != l.RIGHT_BRACE && t != l.SEMICOLON {
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

// Statements that ARE expressions - TODO: put blockStatement here
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

	for p.match(l.COMMA) {
		left, err := p.assign()
		if err != nil {
			return expr, err
		}
		expr = Sequence{Left: left, Right: expr}
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

		expr = Assign{Target: expr, Operator: op, Value: right}
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

	for p.match(l.CHECK_NAV, l.BANG_NAV, l.DOT) {
		op := p.previous()

		right, err := p.validate()
		if err != nil {
			return expr, err
		}

		expr = Access{Left: expr, Right: right, Operator: op}
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
	expr, err := p.identifier(true)
	if err != nil {
		return expr, err
	}

	for p.match(l.COLON) {
		op := p.previous()
		actual := p.Current
		right, err := p.identifier(true)
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

func (p *Parser) identifier(couldBeFunction bool) (Expr, error) {
	if p.match(l.IDENTIFIER) {
		id := p.previous()
		args := []Expr{}

		// Consider it could be a Position Access
		var isAccess bool
		var expr Expr
		var err error
		var pos Expr

		expr = Identifier{id}
		for p.match(l.LEFT_BRACKET) {
			isAccess = true
			if pos, err = p.expression(); err != nil {
				return nil, err
			}

			if _, err := p.consume(l.RIGHT_BRACKET); err != nil {
				return nil, err
			}

			expr = PositionAccess{Expression: expr, Pos: pos}
		}

		if isAccess {
			return expr, nil
		} else if couldBeFunction {
			// if it has other identifiers after, it is a caller
			// otherwise just a normal identifier
			for {
				expr, err := p.identifier(false)
				// In this case is not an error, just ended the function arguments
				// Error driven architecture, sorry...
				if err != nil {
					if (err.(e.NeonError)).ErrorType == e.PARSER_DEAD_END {
						break
					} else {
						return nil, err
					}
				}

				args = append(args, expr)
			}
		}

		if len(args) == 0 {
			expr = Identifier{id}
			return expr, nil
		} else {
			return Caller{Name: id, Args: args}, nil
		}
	}

	return p.primary()
}

func (p *Parser) primary() (Expr, error) {
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

	return p.typeName(false)
}

func (p *Parser) typeName(onlyType bool) (Expr, error) {
	if p.match(l.INT, l.I8, l.I16, l.I32, l.I64, l.UINT, l.U8, l.U16, l.U32, l.U64, l.FLOAT, l.F32, l.F64, l.BOOL, l.CHAR, l.STRING, l.BYTE, l.ANY) {
		return Type{Name: p.previous()}, nil
	}

	if onlyType {
		token := p.peek()
		return nil, e.Error(token.Line, token.Column, token.Lexeme, e.PARSER, fmt.Sprintf("expected a type: %v", token))
	} else {
		return p.array()
	}
}

func (p *Parser) array() (Expr, error) {
	if p.match(l.LEFT_BRACKET) {
		var expr Expr
		var err error

		if expr, err = p.expression(); err != nil {
			return nil, err
		}

		switch expr.(type) {
		case Type:
			return p.arrayType(expr)
		case ArrayType:
			return p.arrayType(expr)
		default:
			return p.arrayLiteral(expr, true)
		}
	}

	return p.group()
}

func (p *Parser) arrayLiteral(expr Expr, inferDeclaration bool) (Expr, error) {
	var values []Expr
	var literal Expr
	var err error

loop:
	switch e := expr.(type) {
	case Sequence:
		values = append(values, e.Left)
		expr = e.Right
		goto loop
	default:
		values = append(values, e)
	}

	utils.Reverse[Expr](values)

	if _, err = p.consume(l.RIGHT_BRACKET); err != nil {
		return nil, err
	}

	literal = ArrayLiteralRaw{Values: values}
	if inferDeclaration {
		return ArrayConstructor{Typing: nil, Values: literal, ActualSize: len(values)}, nil
	} else {
		return literal, nil
	}
}

func (p *Parser) arrayType(expr Expr) (Expr, error) {
	var arraySize Expr
	var arrayType Expr
	var isSlice bool
	var literal Expr
	var err error

	// If have size is an array, if not is a slice
	if p.match(l.COLON) {
		arraySize, err = p.expression()
		isSlice = false
		if err != nil {
			return nil, err
		}
	} else {
		isSlice = true
	}

	if _, err := p.consume(l.RIGHT_BRACKET); err != nil {
		return nil, err
	}

	arrayType = ArrayType{Typing: expr, MaxSize: arraySize, IsSlice: isSlice}

	if p.match(l.LEFT_BRACKET) {
		var expr Expr
		var err error

		if expr, err = p.expression(); err != nil {
			return nil, err
		}

		if literal, err = p.arrayLiteral(expr, false); err != nil {
			return nil, err
		}

		return ArrayConstructor{Typing: arrayType, Values: literal, ActualSize: len(literal.(ArrayLiteralRaw).Values)}, nil
	} else {
		return arrayType, nil //ArrayType
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
		token := p.peek()
		return nil, e.Error(token.Line, token.Column, token.Lexeme, e.PARSER, fmt.Sprintf("expected a block statement, found: %v", token))
	}

	return p.deadEnd()
}

func (p *Parser) deadEnd() (Expr, error) {
	token := p.Tokens[p.Current]
	return nil, e.Error(token.Line, token.Column, token.Lexeme, e.PARSER_DEAD_END, fmt.Sprintf("expect expression, found: %v", token))
}
