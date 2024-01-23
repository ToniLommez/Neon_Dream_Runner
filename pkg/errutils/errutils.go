package errutils

import (
	"fmt"
	"os"
)

type neonError struct {
	line    int
	column  int
	message string
}

func (e neonError) Error() string {
	return fmt.Sprintf("[Line %d, Column %d] %s", e.line, e.column, e.message)
}

func Error(line int, column int, message string) error {
	return neonError{
		line:    line,
		column:  column,
		message: message,
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
