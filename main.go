package main

import (
	"fmt"
	"os"
	"os/user"

	"github.com/odas0r/b/repl"
)

func main() {
	user, err := user.Current()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Hello %s!, welcome to B programming language!", user.Username)
	fmt.Printf("Feel free to type in commands")

	repl.Start(os.Stdin, os.Stdout)
}
