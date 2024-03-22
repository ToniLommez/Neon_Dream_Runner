package parser

import "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"

type Neon struct {
	Tokens       []lexer.Token
	TokensBuffer []lexer.Token
	Main         []Stmt
}
