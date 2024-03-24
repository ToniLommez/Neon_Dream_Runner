package parser

import (
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
)

type Program struct {
	IsLive       bool
	Text         []string
	Tokens       []lexer.Token
	TokensBuffer []lexer.Token
	Main         Scope
}

func (p *Program) Init(isLive bool) {
	p.IsLive = isLive
	p.Main.Values.Init()
}
