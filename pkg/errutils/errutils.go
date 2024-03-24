package errutils

import (
	"fmt"
	"os"
	"strings"
)

const (
	LEXER                  = "lexer"
	PARSER                 = "parser"
	RUNTIME                = "runtime"
	UNTERMINATED_STATEMENT = "unterminated_statement"
)

type NeonError struct {
	Line      int
	Column    int
	Lexeme    string
	ErrorType string
	Message   string
}

func (e NeonError) Error() string {
	space := strings.Repeat(" ", e.Column+1)
	pointer := strings.Repeat("^", len(e.Lexeme))

	message := space + pointer
	if len(message) == 1 {
		message = ""
	} else {
		message += "\n"
	}

	message += fmt.Sprintf("| %s\n", e.Message)
	message += fmt.Sprintf("| [Line %d, Column %d] - %s error", e.Line, e.Column, e.ErrorType)
	return message
}

func Error(line int, column int, lexeme string, errorType string, message string) error {
	return NeonError{
		Line:      line,
		Column:    column,
		Lexeme:    lexeme,
		ErrorType: errorType,
		Message:   message,
	}
}

func Deal(err error, line string) (fatal error) {
	if myErr, ok := err.(NeonError); ok {
		if line != "" {
			fmt.Fprintf(os.Stderr, "> %s\n", line)
		}
		fmt.Fprintf(os.Stderr, "%s\n", myErr)
	} else {
		fatal = fmt.Errorf("fatal error: %s", err)
	}
	return
}
