package neon

import (
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/parser"
)

type Neon struct {
	Text         []string
	IsLive       bool
	Tokens       []lexer.Token
	TokensBuffer []lexer.Token
	Main         []parser.Stmt
}
