package main

import (
	"fmt"
	"os"
)

func main() {
	if len(os.Args) > 1 {
		os.Exit(64) //TODO: better response
	} else if len(os.Args) == 1 {
		if err := runFile(os.Args[0]); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := runPrompt(); err != nil {
			fmt.Println(err)
		}
	}
}
