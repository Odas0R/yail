package object

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"strconv"
	"strings"
)

var Builtins = []struct {
	Name    string
	Builtin *Builtin
}{
	{
		"len",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments. got=%d, want=1",
					len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				return &Integer{Value: int64(len(arg.Elements))}
			case *String:
				return &Integer{Value: int64(len(arg.Value))}
			default:
				return newError("argument to `len` not supported, got %s",
					args[0].Type())
			}
		},
		},
	},
	{
		"pow",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for 'pow'. got=%d, want=2", len(args))
			}
			var base, exponent float64

			// Handle the base argument
			switch args[0].Type() {
			case INTEGER_OBJ:
				base = float64(args[0].(*Integer).Value)
			case FLOAT_OBJ:
				base = args[0].(*Float).Value
			default:
				return newError("first argument to 'pow' must be INTEGER or FLOAT, got %s", args[0].Type())
			}

			// Handle the exponent argument
			switch args[1].Type() {
			case INTEGER_OBJ:
				exponent = float64(args[1].(*Integer).Value)
			case FLOAT_OBJ:
				exponent = args[1].(*Float).Value
			default:
				return newError("second argument to 'pow' must be INTEGER or FLOAT, got %s", args[1].Type())
			}

			result := math.Pow(base, exponent)
			return &Float{Value: result}
		}},
	},
	{
		"square_root",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for 'square_root'. got=%d, want=1", len(args))
			}
			var number float64

			// Handle the number argument
			switch args[0].Type() {
			case INTEGER_OBJ:
				number = float64(args[0].(*Integer).Value)
			case FLOAT_OBJ:
				number = args[0].(*Float).Value
			default:
				return newError("argument to 'square_root' must be INTEGER or FLOAT, got %s", args[0].Type())
			}

			if number < 0 {
				return newError("argument to 'square_root' must be non-negative")
			}

			result := math.Sqrt(number)
			return &Float{Value: result}
		}},
	},
	{
		"gen",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 2 {
				return newError("wrong number of arguments for 'gen'. got=%d, want=2", len(args))
			}
			if args[0].Type() != INTEGER_OBJ || args[1].Type() != INTEGER_OBJ {
				return newError("arguments to 'gen' must be INTEGERs, got %s and %s", args[0].Type(), args[1].Type())
			}
			start := args[0].(*Integer).Value
			end := args[1].(*Integer).Value
			result := &Array{}
			for i := start; i <= end; i++ {
				result.Elements = append(result.Elements, &Integer{Value: i})
			}
			return result
		}},
	},
	{
		"write",
		&Builtin{Fn: func(args ...Object) Object {
			var output []string
			for _, arg := range args {
				output = append(output, arg.Inspect())
			}
			fmt.Println(strings.Join(output, " "))
			return nil
		}},
	},
	{
		"write_all",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for 'write_all'. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				var output []string
				for _, el := range arg.Elements {
					output = append(output, el.Inspect())
				}
				fmt.Println(strings.Join(output, ", "))
			case *Struct:
				var output []string
				for _, el := range arg.Attributes {
					output = append(output, el.Inspect())
				}
				fmt.Println(strings.Join(output, ", "))
			default:
				return newError("argument to 'write_all' must be ARRAY or STRUCT, got %s", args[0].Type())
			}
			return nil
		}},
	}, {
		"write_string",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for 'write_string'. got=%d, want=1", len(args))
			}
			switch arg := args[0].(type) {
			case *Array:
				for _, el := range arg.Elements {
					if intEl, ok := el.(*Integer); ok {
						fmt.Print(string(rune(intEl.Value)))
					} else {
						return newError("write_string array elements must be INTEGERS, got %s", el.Type())
					}
				}
				fmt.Println()
			default:
				return newError("argument to 'write_string' must be ARRAY, got %s", args[0].Type())
			}
			return nil
		}},
	},
	{
		"read",
		&Builtin{Fn: func(args ...Object) Object {
			reader := bufio.NewReader(os.Stdin)
			line, _ := reader.ReadString('\n')
			return &String{Value: line}
		}},
	},
	{
		"read_all",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for 'read_all'. got=%d, want=1", len(args))
			}
			reader := bufio.NewReader(os.Stdin)
			switch arg := args[0].(type) {
			case *Array:
				for i := range arg.Elements {
					fmt.Printf("v[%d]: ", i)
					input, _ := reader.ReadString('\n')
					arg.Elements[i] = &Integer{Value: parseInput(input)}
				}
			case *Struct:
				for key := range arg.Attributes {
					fmt.Printf("struct %s\n%s: ", arg.Type(), key)
					input, _ := reader.ReadString('\n')
					arg.Attributes[key] = &Integer{Value: parseInput(input)}
				}
			default:
				return newError("argument to 'read_all' must be ARRAY or STRUCT, got %s", args[0].Type())
			}
			return nil
		}},
	},
	{
		"read_string",
		&Builtin{Fn: func(args ...Object) Object {
			if len(args) != 1 {
				return newError("wrong number of arguments for 'read_string'. got=%d, want=1", len(args))
			}
			arg, ok := args[0].(*Array)
			if !ok {
				return newError("argument to 'read_string' must be ARRAY, got %s", args[0].Type())
			}
			reader := bufio.NewReader(os.Stdin)
			input, _ := reader.ReadString('\n')
			input = strings.TrimSuffix(input, "\n") // remove the newline character at the end
			for i, ch := range input {
				if i < len(arg.Elements) {
					arg.Elements[i] = &Integer{Value: int64(ch)}
				} else {
					break
				}
			}
			for i := len(input); i < len(arg.Elements); i++ { // fill the rest with zeros
				arg.Elements[i] = &Integer{Value: 0}
			}
			return nil
		}},
	},
}

func GetBuiltinByName(name string) *Builtin {
	for _, def := range Builtins {
		if def.Name == name {
			return def.Builtin
		}
	}
	return nil
}

func newError(format string, a ...interface{}) *Error {
	return &Error{Message: fmt.Sprintf(format, a...)}
}

func parseInput(input string) int64 {
	trimmed := strings.TrimSpace(input)
	num, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		return 0
	}
	return num
}
