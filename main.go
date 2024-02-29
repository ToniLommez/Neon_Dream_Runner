package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
	p "github.com/ToniLommez/Neon_Dream_Runner/pkg/parser"
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
	prompt := ""

	for {
		fmt.Print("> ")
		if s.Scan() {
			prompt = s.Text()

			switch prompt {
			case "exit":
				return nil
			case "clear":
				utils.ClearScreen()
			case "":
			default:
				err = run(prompt)
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

func run(input string) (err error) {
	var ts []l.Token
	s := l.NewScanner(input)
	if ts, err = s.ScanTokens(); err != nil {
		return err
	}

	parser := p.NewParser(ts)
	expr := parser.Parse()
	result, err := p.Evaluate(expr)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(result)
	}

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
