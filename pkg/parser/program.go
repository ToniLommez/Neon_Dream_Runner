package parser

import (
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Program struct {
	IsLive       bool
	Text         []string
	Tokens       []lexer.Token
	TokensBuffer []lexer.Token
	Packages     []Package
}

type Package struct {
	Functions map[string]FnStmt
	Main      Scope
}

func (p *Program) IsFunction(s Stmt) (bool, FnStmt) {
	switch f := s.(type) {
	case FnStmt:
		// TODO: implement the correct context of the package
		f.Context.Owner = &p.Packages[0]
		f.Context.Parent = &p.Packages[0].Main
		return true, f
	default:
		return false, FnStmt{} // TODO: remove this...
	}
}

func (p *Program) Init(isLive bool) {
	p.IsLive = isLive
	p.Packages = append(p.Packages, Package{Functions: make(map[string]FnStmt), Main: Scope{}})
	p.Packages[0].Main.Init()
	p.Packages[0].Main.Owner = &p.Packages[0]
}
