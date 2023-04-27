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

	// if a filepath is given as an argument, run the file and exit
	if len(os.Args) > 1 {
		repl.RunFile(os.Args[1])
		os.Exit(0)
		return
	}

	fmt.Printf("%s\n", repl.YAIL)
	fmt.Printf("Hello %s!, welcome to YAIL programming language!", user.Username)
	fmt.Printf(" Feel free to type in commands")

	repl.Start(os.Stdin, os.Stdout)
}
