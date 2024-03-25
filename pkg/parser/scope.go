package parser

import (
	"fmt"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Variable struct {
	Value       any
	Type        int
	TypeDefined bool
	Mutable     bool
	Nullable    bool
	Initialized bool
}

type Scope struct {
	Statements []Stmt
	Values     map[string]Variable
	Parent     *Scope // Cactus-Stack
}

func (s *Scope) Init() {
	s.Values = make(map[string]Variable)
}

func (s *Scope) Define(l LetStmt, value any) (any, error) {
	name := l.Name.Lexeme
	tmp := s.Values[name]

shadowing: // just some weird optimization
	if _, found := s.Values[name]; found {
		name = "§" + name
		goto shadowing
	} else {
		s.Values[name] = tmp
	}

	defined := l.Type != UNDEFINED && l.Type != UNKNOWN && l.Type != NIL
	s.Values[l.Name.Lexeme] = Variable{Type: l.Type, Value: value, TypeDefined: defined, Mutable: l.Mutable, Nullable: l.Nullable, Initialized: l.Initializer != nil}

	return value, nil
}

// return => type, value, isDefined, error
func (s *Scope) Get(name l.Token) (int, any, bool, error) {
	currentScope := s
	for currentScope != nil {
		value, found := currentScope.Values[name.Lexeme]
		if found {
			if !value.Initialized {
				return value.Type, value.Value, false, e.Error(name.Line, name.Column, name.Lexeme, e.RUNTIME, "variable created but not defined")
			}
			return value.Type, value.Value, true, nil
		}
		currentScope = currentScope.Parent
	}
	return UNKNOWN, nil, false, e.Error(name.Line, name.Column, name.Lexeme, e.RUNTIME, "variable not found")
}

func (s *Scope) Set(target l.Token, newValue any) (any, error) {
	currentScope := s
	var v Variable
	var found bool

	// Walk up the scope until find
	for currentScope != nil {
		v, found = currentScope.Values[target.Lexeme]
		if found {
			break
		}
		currentScope = currentScope.Parent
	}

	if currentScope == nil || !found {
		return nil, e.Error(target.Line, target.Column, target.Lexeme, e.RUNTIME, "variable not found")
	}

	if v.Initialized && !v.Mutable {
		return nil, e.Error(target.Line, target.Column, target.Lexeme, e.RUNTIME, "cannot assign because "+target.Lexeme+" is immutable")
	}

	if newValue == nil && !v.Nullable {
		return nil, e.Error(target.Line, target.Column, target.Lexeme, e.RUNTIME, target.Lexeme+" cannot be set to nil")
	}

	tp := getType(newValue)
	if v.TypeDefined && v.Type != tp && newValue != nil {
		return nil, e.Error(target.Line, target.Column, target.Lexeme, e.RUNTIME, "expected type to assign: "+typeToString(v.Type)+", found: "+typeToString(tp))
	}

	v.Type = getType(newValue)
	v.TypeDefined = true
	v.Initialized = true
	v.Value = newValue
	currentScope.Values[target.Lexeme] = v

	return v.Value, nil
}

func (s *Scope) Debug() {
	var prompt string
	for i, j := range s.Values {
		prompt = "%s = {Type: %s, Value: %v, TypeDefined: %v, Mutable: %v, Nullable: %v, Initialized: %v}\n"
		fmt.Printf(prompt, i, typeToString(j.Type), j.Value, j.TypeDefined, j.Mutable, j.Nullable, j.Initialized)
	}
}
