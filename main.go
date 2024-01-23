package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
	"github.com/ToniLommez/Neon_Dream_Runner/pkg/utils"
)

func runFile(path string) error {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("error opening file: %s", err)
	}
	defer file.Close()
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("erro reading file: %s", err)
	}

	return run(string(content))
}

func runPrompt() (err error) {
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		if s.Scan() {
			input := s.Text()

			switch input {
			case "exit":
				return nil
			case "clear":
				utils.ClearScreen()
			default:
				err = run(input)
				if err != nil {
					fatal := errutils.Deal(err)
					if fatal != nil {
						return fatal
					}
				}
			}
		} else {
			return s.Err() // Retorna um erro se a leitura falhar
		}
	}
}

func run(input string) error {
	s := lexer.NewScanner(input)
	ts, err := s.ScanTokens()
	if err != nil {
		return err
	}
	for _, t := range ts {
		fmt.Printf("%s ", t)
	}
	fmt.Printf("\n")
	return nil
}

func main() {
	if len(os.Args) > 2 {
		os.Exit(64) //TODO: better response
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := runPrompt(); err != nil {
			fmt.Println(err)
		}
	}
}
