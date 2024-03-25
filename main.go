package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"

	e "github.com/ToniLommez/Neon_Dream_Runner/pkg/errutils"
	l "github.com/ToniLommez/Neon_Dream_Runner/pkg/lexer"
	p "github.com/ToniLommez/Neon_Dream_Runner/pkg/parser"
	u "github.com/ToniLommez/Neon_Dream_Runner/pkg/utils"
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

	if len(content) > 0 && content[len(content)-1] != '\n' {
		content = append(content, '\n')
	}

	var neon p.Program
	neon.Init(false)
	neon.Text = strings.Split(string(content), "\n")

	_, err = run(string(content), true, &neon)
	if err != nil {
		var fatal error
		if myErr, ok := err.(e.NeonError); ok {
			fatal = e.Deal(err, neon.Text[myErr.Line-1])
		} else {
			fatal = e.Deal(err, "")
		}

		if fatal != nil {
			return fatal
		}
	}
	return nil
}

func runRepl() error {
	s := bufio.NewScanner(os.Stdin)
	var neon p.Program
	var err error
	neon.Init(true)

	depth := 0
	prompt := ""

	for {
		if depth == 0 {
			fmt.Print("\033[34m• \033[0m")
		} else {
			fmt.Print("\033[32m• \033[0m")
		}

		if s.Scan() {
			prompt = s.Text()

			switch prompt {
			case "exit":
				return nil
			case "clear":
				u.ClearScreen()
			default:
				depth, err = run(prompt, false, &neon)
				if err != nil {
					fatal := e.Deal(err, "")
					if fatal != nil {
						return fatal
					}
				}
			}
		} else {
			return s.Err()
		}
	}
}

// TODO: transfer anything that belongs do "n" to a new package named neon and encapsulate everything
// TODO: fork into runRepl and runFile
func run(input string, isFile bool, neon *p.Program) (depth int, err error) {
	// Scan new tokens
	var ts []l.Token
	s := l.NewScanner(input)
	if ts, err = s.ScanTokens(isFile); err != nil {
		return 0, err
	}

	// Concatenate with buffered tokens to parse
	neon.TokensBuffer = append(neon.TokensBuffer, ts...)

	// Parse
	pr := p.NewParser(neon.TokensBuffer)
	statement, err := pr.Parse()

	// In REPL if a incomplete statemente is found, buffer it and wait till complete before evaluate
	if !isFile && err != nil && (err.(e.NeonError)).ErrorType == e.UNTERMINATED_STATEMENT {
		return pr.Depth, nil
	} else if err != nil {
		neon.TokensBuffer = nil
		return 0, err
	}

	// If correctly parsed save the statement
	neon.Tokens = neon.TokensBuffer
	neon.TokensBuffer = nil
	neon.Main.Statements = statement

	// Evaluate the AST
	res, err := neon.Main.Interpret()
	if err != nil {
		return 0, err
	}

	if res != nil {
		tmp := fmt.Sprintf("%v", res)
		color := "\033[38;2;150;240;240m"
		reset := "\033[0m"
		fmt.Printf("%s%v%s\n", color, strings.Replace(tmp, "\\n", "\n", -1), reset)
	}

	return 0, nil
}

func main() {
	if len(os.Args) > 2 {
		os.Exit(64) // TODO: better response
	} else if len(os.Args) == 2 {
		if err := runFile(os.Args[1]); err != nil {
			fmt.Println(err)
		}
	} else {
		if err := runRepl(); err != nil {
			fmt.Println(err)
		}
	}
}
