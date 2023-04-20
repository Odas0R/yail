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

	fmt.Printf("%s\n", repl.YAIL)
	fmt.Printf("Hello %s!, welcome to YAIL programming language!", user.Username)
	fmt.Printf(" Feel free to type in commands")

	repl.Start(os.Stdin, os.Stdout)
}
