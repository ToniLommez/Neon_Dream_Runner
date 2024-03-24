package parser

import (
	"fmt"

	er "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
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

type Environment struct {
	Values map[string]Variable
}

func (e *Environment) Init() {
	e.Values = make(map[string]Variable)
}

func (e *Environment) Define(l LetStmt, value any) error {
	name := l.Name.Lexeme
	tmp := e.Values[name]

shadowing:
	if _, found := e.Values[name]; found {
		name = "§" + name
		goto shadowing
	} else {
		e.Values[name] = tmp
	}

	defined := l.Type != UNDEFINED && l.Type != UNKNOWN && l.Type != NIL
	e.Values[l.Name.Lexeme] = Variable{Type: l.Type, Value: value, TypeDefined: defined, Mutable: l.Mutable, Nullable: l.Nullable, Initialized: l.Initializer != nil}

	return nil
}

// return => type, value, isDefined, error
func (e *Environment) Get(name l.Token) (int, any, bool, error) {
	value, found := e.Values[name.Lexeme]
	if !found {
		return UNKNOWN, nil, true, er.Error(name.Line, name.Column, name.Lexeme, er.RUNTIME, "variable not found")
	}

	if !value.Initialized {
		return value.Type, value.Value, false, er.Error(name.Line, name.Column, name.Lexeme, er.RUNTIME, "variable created but not defined")
	}
	return value.Type, value.Value, true, nil
}

func (e *Environment) Set(target l.Token, newValue any) (any, error) {
	v, found := e.Values[target.Lexeme]
	if !found {
		return nil, er.Error(target.Line, target.Column, target.Lexeme, er.RUNTIME, "variable not found")
	}

	if v.Initialized && !v.Mutable {
		return nil, er.Error(target.Line, target.Column, target.Lexeme, er.RUNTIME, fmt.Sprintf("cannot assign because %s is immutable", target.Lexeme))
	}

	if newValue == nil && !v.Nullable {
		return nil, er.Error(target.Line, target.Column, target.Lexeme, er.RUNTIME, fmt.Sprintf("%s cannot be set to nil", target.Lexeme))
	}

	tp := getType(newValue)
	if !v.TypeDefined {
		v.Type = tp
		v.TypeDefined = true
	} else if v.Type != tp && newValue != nil {
		msg := fmt.Sprintf("expected type to assign: %s, found: %s", typeToString(v.Type), typeToString(tp))
		return nil, er.Error(target.Line, target.Column, target.Lexeme, er.RUNTIME, msg)
	}

	v.Initialized = true // it's faster to assign than to verify
	v.Value = newValue

	e.Values[target.Lexeme] = v

	return v.Value, nil
}

func (e *Environment) Debug() {
	var prompt string
	for i, j := range e.Values {
		prompt = "%s = {Type: %s, Value: %v, TypeDefined: %v, Mutable: %v, Nullable: %v, Initialized: %v}\n"
		fmt.Printf(prompt, i, typeToString(j.Type), j.Value, j.TypeDefined, j.Mutable, j.Nullable, j.Initialized)
	}
}
