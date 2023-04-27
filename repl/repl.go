package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/odas0r/yail/lexer"
	"github.com/odas0r/yail/parser"
	"github.com/odas0r/yail/token"
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

func RunFile(path string) {
	file, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening file: %s (%s) ", path, err)
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("Error reading file: %s (%s) ", path, err)
		return
	}

	// if file .out exists, delete it
	if _, err := os.Stat(path + ".out"); err == nil {
		err = os.Remove(path + ".out")
		if err != nil {
			fmt.Printf("Error deleting file: %s (%s) ", path+".out", err)
			return
		}
	}

	// create or re a file output based on the input file
	out, err := os.Create(path + ".out")
	if err != nil {
		fmt.Printf("Error creating file: %s (%s) ", path, err)
		return
	}
	defer out.Close()

	l := lexer.New(string(content))
	p := parser.New(l)

	// create the AST by parsing the tokens
	program := p.ParseProgram()

	out.WriteString("\n=========================================")
	out.WriteString(" AST ")
	out.WriteString("=========================================\n")

	if len(p.Errors()) != 0 {
		out.WriteString("\n")
		printParserErrors(out, p.Errors())
		out.WriteString("\n")

		// create a new lexer to print the tokens
		l = lexer.New(string(content))

		out.WriteString("=========================================")
		out.WriteString(" TOKENS ")
		out.WriteString("=========================================\n\n")
		for {
			tok := l.NextToken()
			fmt.Fprintf(out, "%+v\n", tok)
			if tok.Type == token.EOF {
				break
			}
		}

		// print the errors on the console
		printParserErrors(os.Stderr, p.Errors())
		os.Exit(1)
	} else {
		out.WriteString("\n")
		out.WriteString(program.PrintAST())
		out.WriteString("\n")
	}

	// create a new lexer to print the tokens
	l = lexer.New(string(content))

	out.WriteString("=========================================")
	out.WriteString(" TOKENS ")
	out.WriteString("=========================================\n\n")
	for {
		tok := l.NextToken()
		fmt.Fprintf(out, "%+v\n", tok)
		if tok.Type == token.EOF {
			break
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	io.WriteString(out, "\nOh boy! Something clearly went wrong!\n\n")
	io.WriteString(out, "parser errors:\n")
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
