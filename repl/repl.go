package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/odas0r/yail/lexer"
	"github.com/odas0r/yail/parser"
)

const YAIL = `
██    ██  █████  ██ ██ 
 ██  ██  ██   ██ ██ ██ 
  ████   ███████ ██ ██ 
   ██    ██   ██ ██ ██ 
   ██    ██   ██ ██ ███████ 
`

const PROMPT = "\n>> "

func Start(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	for {
		fmt.Fprint(out, PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}

		io.WriteString(out, program.String())
		io.WriteString(out, "\n")
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "\nOh boy! Something clearly went wrong!\n\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
