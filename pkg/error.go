package neon

import (
	"fmt"
	"os"
)

type ErrorWithPosition struct {
	Line    int
	Column  int
	Message string
}

func (e *ErrorWithPosition) Error() string {
	return fmt.Sprintf("[Line %d, Column %d] Error: %s", e.Line, e.Column, e.Message)
}

func Deal(err error) (fatal error) {
	if myErr, ok := err.(*ErrorWithPosition); ok {
		fmt.Fprintf(os.Stderr, "%v\n", myErr)
	} else {
		fatal = fmt.Errorf("Fatal Error: %s", err)
	}
	return
}
