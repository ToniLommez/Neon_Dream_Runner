package neon

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

func runFile(path string) error {
	file, err := os.Open(path)
	defer file.Close()
	if err != nil {
		return fmt.Errorf("Error opening file: %s", err)
	}
	content, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("Erro reading file: %s", err)
	}

	fmt.Println(string(content))

	return nil
}

func runPrompt() (err error) {
	s := bufio.NewScanner(os.Stdin)

	for {
		fmt.Print("> ")
		input := s.Text()

		switch input {
		case "exit":
			return nil
		case "clear":
			clearScreen()
		default:
			err = run(input)
			if err != nil {
				fatal := Deal(err)
				if fatal != nil {
					return fatal
				}
			}
		}
	}
}

func run(input string) error {
	fmt.Println(input)
	return nil
}
