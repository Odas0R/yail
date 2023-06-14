package repl

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"

	"github.com/odas0r/yail/compiler"
	"github.com/odas0r/yail/lexer"
	"github.com/odas0r/yail/object"
	"github.com/odas0r/yail/parser"
	"github.com/odas0r/yail/token"
	"github.com/odas0r/yail/vm"
)

const YAIL = `
██    ██  █████  ██ ██ 
 ██  ██  ██   ██ ██ ██ 
  ████   ███████ ██ ██ 
   ██    ██   ██ ██ ██ 
   ██    ██   ██ ██ ███████ 
`

const PROMPT = "yail> "
const PROMPT_KEEP_WRITING = " ...> "

func RunAst(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)
	var inputLines []string

	fmt.Fprint(out, "\n")
	fmt.Fprint(out, PROMPT)
	for {
		scanned := scanner.Scan()
		if !scanned {
			return
		}

		line := scanner.Text()

		if line == "exit" || line == "quit" {
			os.Exit(0)
		}

		// If an empty line is encountered, parse the accumulated input
		if line == "" {
			input := strings.Join(inputLines, "\n")
			l := lexer.New(input)
			p := parser.New(l)

			program := p.ParseProgram()
			if len(p.Errors()) != 0 {
				printParserErrors(out, p.Errors())
			} else {
				io.WriteString(out, "\n")
				io.WriteString(out, program.Stringify(1))
			}

			// Reset
			inputLines = []string{}
			fmt.Fprint(out, "\n")
			fmt.Fprint(out, PROMPT)
		} else {
			// Accumulate non-empty lines
			inputLines = append(inputLines, line)
			fmt.Fprint(out, PROMPT_KEEP_WRITING)
		}
	}
}

func RunVm(in io.Reader, out io.Writer) {
	scanner := bufio.NewScanner(in)

	constants := []object.Object{}
	globals := make([]object.Object, vm.GlobalsSize)
	symbolTable := compiler.NewSymbolTable()

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
		comp := compiler.NewWithState(symbolTable, constants)
		err := comp.Compile(program)
		if err != nil {
			fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
			continue
		}
		code := comp.Bytecode()
		constants = code.Constants
		machine := vm.NewWithGlobalsStore(code, globals)
		err = machine.Run()
		if err != nil {
			fmt.Fprintf(out, "Woops! Executing bytecode failed:\n %s\n", err)
			continue
		}
		lastPopped := machine.LastPoppedStackElem()
		io.WriteString(out, lastPopped.Inspect())
		io.WriteString(out, "\n")
	}
}

func RunFileVm(path string) {
	constants := []object.Object{}
	symbolTable := compiler.NewSymbolTable()

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

	if len(p.Errors()) != 0 {
		out.WriteString("\n")
		printParserErrors(out, p.Errors())
		out.WriteString("\n")
	}

	comp := compiler.NewWithState(symbolTable, constants)
	err = comp.Compile(program)
	if err != nil {
		fmt.Fprintf(out, "Woops! Compilation failed:\n %s\n", err)
	}

	code := comp.Bytecode()

	out.WriteString("=========================================")
	out.WriteString(" CODE ")
	out.WriteString("=========================================\n")
	out.WriteString(string(content))

	out.WriteString("=========================================")
	out.WriteString(" Constants ")
	out.WriteString("=========================================\n")
	for i, constant := range code.Constants {
		fn, ok := constant.(*object.CompiledFunction)
		if ok {
			fmt.Fprintf(out, "CompiledFunction[%d] Instructions:\n", i)
			out.WriteString(fn.Instructions.String())
			out.WriteString("\n")
			continue
		}

		fmt.Fprintf(out, "%d\t%s\n", i, constant.Inspect())

	}

	out.WriteString("\n")

	out.WriteString("=========================================")
	out.WriteString(" Instructions ")
	out.WriteString("=========================================\n")
	out.WriteString(code.Instructions.String())
	// if the instruction is a compiled fuction, print the instructions
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
		out.WriteString(program.Stringify(1))
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
