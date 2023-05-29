package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/odas0r/yail/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	if len(os.Args) > 1 {
		switch os.Args[1] {
		case "vm":
			fmt.Printf("%s\n", repl.YAIL)
			fmt.Printf("Hello %s!, welcome to YAIL programming language!\n", user.Username)
			fmt.Printf("Feel free to type in commands\n")

			repl.RunVm(os.Stdin, os.Stdout)
			os.Exit(0)
			return
		case "ast":
			// if a filepath is given as an argument, run the file and exit
			if os.Args[2] != "" {
				repl.RunFile(os.Args[2])
				os.Exit(0)
				return
			}

			fmt.Printf("%s\n", repl.YAIL)
			fmt.Printf("Hello %s!, welcome to YAIL programming language!\n", user.Username)
			fmt.Printf("Feel free to type in commands\n")

			repl.RunAst(os.Stdin, os.Stdout)
			os.Exit(0)
			return
		default:
			fmt.Printf("Please use either vm or ast as an argument\n")
			fmt.Printf("\nvm: run the virtual machine\n")
			fmt.Printf("ast: run the abstract syntax tree\n")
			os.Exit(0)
			return
		}
	}

	fmt.Printf("Please use either vm or ast as an argument\n")
	fmt.Printf("\nvm: run the virtual machine\n")
	fmt.Printf("ast: run the abstract syntax tree\n")
}
