package errutils

import (
	"fmt"
	"os"
	"strings"
)

const (
	LEXER   = "lexer"
	PARSER  = "parser"
	RUNTIME = "runtime"
)

type neonError struct {
	line      int
	column    int
	lexeme    string
	errorType string
	message   string
}

func (e neonError) Error() string {
	space := strings.Repeat(" ", e.column+1)
	pointer := strings.Repeat("^", len(e.lexeme))

	message := space + pointer + "\n"
	message += fmt.Sprintf("| %s\n", e.message)
	message += fmt.Sprintf("| [Line %d, Column %d] - %s error", e.line, e.column, e.errorType)
	return message
}

func Error(line int, column int, lexeme string, errorType string, message string) error {
	return neonError{
		line:      line,
		column:    column,
		lexeme:    lexeme,
		errorType: errorType,
		message:   message,
	}
}

func Deal(err error) (fatal error) {
	if myErr, ok := err.(neonError); ok {
		fmt.Fprintf(os.Stderr, "%s\n", myErr)
	} else {
		fatal = fmt.Errorf("fatal error: %s", err)
	}
	return
}
