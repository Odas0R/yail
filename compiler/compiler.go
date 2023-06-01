package compiler

import (
	"fmt"

	"github.com/odas0r/yail/ast"
	"github.com/odas0r/yail/code"
	"github.com/odas0r/yail/object"
)

type EmittedInstruction struct {
	Opcode   code.Opcode
	Position int
}

type Bytecode struct {
	Instructions code.Instructions
	Constants    []object.Object
}

type CompilationScope struct {
	instructions        code.Instructions
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
}

type Compiler struct {
	instructions        code.Instructions
	constants           []object.Object
	lastInstruction     EmittedInstruction
	previousInstruction EmittedInstruction
	symbolTable         *SymbolTable
	scopes              []CompilationScope
	scopeIndex          int
}

func New() *Compiler {
	mainScope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}

	symbolTable := NewSymbolTable()

	for i, v := range object.Builtins {
		symbolTable.DefineBuiltin(i, v.Name)
	}

	return &Compiler{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
		constants:           []object.Object{},
		symbolTable:         symbolTable,
		scopes:              []CompilationScope{mainScope},
		scopeIndex:          0,
	}
}

func NewWithState(s *SymbolTable, constants []object.Object) *Compiler {
	compiler := New()
	compiler.symbolTable = s
	compiler.constants = constants
	return compiler
}

func (c *Compiler) Compile(node ast.Node) error {
	switch node := node.(type) {
	case *ast.Program:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.ExpressionStatement:
		err := c.Compile(node.Expression)
		if err != nil {
			return err
		}
		c.emit(code.OpPop)

	case *ast.BlockStatement:
		for _, s := range node.Statements {
			err := c.Compile(s)
			if err != nil {
				return err
			}
		}
	case *ast.Identifier:
		symbol, ok := c.symbolTable.Resolve(node.Value)
		if !ok {
			return fmt.Errorf("undefined variable %s", node.Value)
		}

		c.loadSymbol(symbol)
	case *ast.InfixExpression:
		// swap operands for < operator
		if node.Operator == "<" {
			err := c.Compile(node.Right)
			if err != nil {
				return err
			}

			err = c.Compile(node.Left)
			if err != nil {
				return err
			}
			c.emit(code.OpGreaterThan)
			return nil
		}

		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "+":
			c.emit(code.OpAdd)
		case "-":
			c.emit(code.OpSub)
		case "*":
			c.emit(code.OpMul)
		case "/":
			c.emit(code.OpDiv)
		case ">":
			c.emit(code.OpGreaterThan)
		case "==":
			c.emit(code.OpEqual)
		case "!=":
			c.emit(code.OpNotEqual)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.PrefixExpression:
		err := c.Compile(node.Right)
		if err != nil {
			return err
		}
		switch node.Operator {
		case "!":
			c.emit(code.OpBang)
		case "-":
			c.emit(code.OpMinus)
		default:
			return fmt.Errorf("unknown operator %s", node.Operator)
		}
	case *ast.IfExpression:
		err := c.Compile(node.Condition)
		if err != nil {
			return err
		}

		// emit an OpJumpNotTruthy with a bogus value
		jumpNotTruthyPos := c.emit(code.OpJumpNotTruthy, 9999)

		err = c.Compile(node.Consequence)
		if err != nil {
			return err
		}

		if c.lastInstructionIs(code.OpPop) {
			c.removeLastPop()
		}

		// Emit an `OpJump` with a bogus value
		jumpPos := c.emit(code.OpJump, 9999)

		afterConsequencePos := len(c.currentInstructions())
		c.changeOperand(jumpNotTruthyPos, afterConsequencePos)

		if node.Alternative == nil {
			c.emit(code.OpNull)
		} else {
			err = c.Compile(node.Alternative)
			if err != nil {
				return err
			}

			if c.lastInstructionIs(code.OpPop) {
				c.removeLastPop()
			}
		}

		afterAlternativePos := len(c.currentInstructions())
		c.changeOperand(jumpPos, afterAlternativePos)

	case *ast.GlobalStatement:
		// compile the block statement
		for _, stmt := range node.Body.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}

			switch s := stmt.(type) {
			case *ast.VariableStatement:
				symbol := c.symbolTable.Define(s.Name.Value)
				c.emit(code.OpSetGlobal, symbol.Index)
			case *ast.ArrayStatement:
				symbol := c.symbolTable.Define(s.Name.Value)
				c.emit(code.OpSetGlobal, symbol.Index)
			default:
				return fmt.Errorf("global statement must contain only variable statements")
			}
		}
	case *ast.StructsStatement:
		for _, st := range node.Structs {
			err := c.Compile(st)
			if err != nil {
				return err
			}
		}
	case *ast.Struct:
		structName := node.Name.Value

		for _, attr := range node.Attributes {
			attrName := attr.Name.Value
			// compile the attribute value
			err := c.Compile(attr.Value)
			if err != nil {
				return err
			}

			var obj object.Object
			switch attr.Value.(type) {
			case *ast.IntegerLiteral:
				obj = &object.Integer{Value: 0}
			case *ast.Boolean:
				obj = &object.Boolean{Value: false}
			case *ast.StringLiteral:
				obj = &object.String{Value: ""}
			case *ast.ArrayStatement:
				// based on the array type, we should create an array object with
				// the correct type
				arrType := attr.Type.Value
				switch arrType {
				case "int":
					obj = &object.Array{Elements: []object.Object{
						&object.Integer{Value: 0},
					}}
				case "bool":
					obj = &object.Array{Elements: []object.Object{
						&object.Boolean{Value: false},
					}}
				case "float":
					obj = &object.Array{Elements: []object.Object{
						&object.Float{Value: 0.0},
					}}
				case "string":
					obj = &object.Array{Elements: []object.Object{
						&object.String{Value: ""},
					}}
				default:
					obj = &object.Array{Elements: []object.Object{
						&object.Null{},
					}}
				}
			}

			c.symbolTable.DefineAttribute(structName, attrName, obj)
		}

		// emit the OpStruct instruction
		c.emit(code.OpStruct, len(node.Attributes))
	case *ast.AccessorExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}

		// TODO: this is not working!

		// attr, ok := node.Index[0].(*ast.Identifier)
		// if !ok {
		// 	return fmt.Errorf("accessor expression must contain an identifier")
		// }

		// emit the OpGetAttr instruction
		c.emit(code.OpGetAttribute, 0)

	case *ast.ForStatement:
		// TODO
	case *ast.WhileStatement:
		// TODO
	case *ast.IncrementStatement:
		// TODO
	case *ast.DecrementStatement:
		// TODO
	case *ast.PlusEqualsStatement:
		// TODO
	case *ast.MinusEqualsStatement:
		// TODO
	case *ast.MultEqualsStatement:
		// TODO

	case *ast.ConstStatement:
		// compile the block statement
		for _, stmt := range node.Body.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}

			switch s := stmt.(type) {
			case *ast.VariableStatement:
				// check if exists, if so return error
				_, ok := c.symbolTable.Resolve(s.Name.Value)
				if ok {
					return fmt.Errorf("variable %s already defined", s.Name.Value)
				}

				symbol := c.symbolTable.Define(s.Name.Value)
				c.emit(code.OpSetGlobal, symbol.Index)
			case *ast.ArrayStatement:
				_, ok := c.symbolTable.Resolve(s.Name.Value)
				if ok {
					return fmt.Errorf("variable %s already defined", s.Name.Value)
				}
				symbol := c.symbolTable.Define(s.Name.Value)
				c.emit(code.OpSetGlobal, symbol.Index)
			default:
				return fmt.Errorf("const statement must contain only variable statements")
			}
		}

	case *ast.LocalStatement:
		// compile the block statement
		for _, stmt := range node.Body.Statements {
			err := c.Compile(stmt)
			if err != nil {
				return err
			}

			switch s := stmt.(type) {
			case *ast.VariableStatement:
				symbol := c.symbolTable.Define(s.Name.Value)
				c.emit(code.OpSetLocal, symbol.Index)
			case *ast.ArrayStatement:
				symbol := c.symbolTable.Define(s.Name.Value)
				c.emit(code.OpSetLocal, symbol.Index)
			default:
				return fmt.Errorf("local statement must contain only variable statements")
			}
		}

	case *ast.VariableStatement:
		err := c.Compile(node.Value)
		if err != nil {
			return err
		}

		// check if varType exists in the symbol table
		symb, ok := c.symbolTable.ResolveStruct(node.Type.Value)
		if ok {
			// we should compile the struct and emit an OpConstant with the
			// object struct
			strct := &object.Struct{
				Attributes: symb.Atributes,
			}

			// emit the OpConstant
			c.emit(code.OpConstant, c.addConstant(strct))
		}

	case *ast.IndexExpression:
		err := c.Compile(node.Left)
		if err != nil {
			return err
		}
		err = c.Compile(node.Index)
		if err != nil {
			return err
		}

		c.emit(code.OpIndex)

	case *ast.FunctionStatement:
		c.enterScope()

		for _, p := range node.Parameters {
			c.symbolTable.Define(p.Name.Value)
		}

		hasReturnValue := false
		for _, stmt := range node.Body.Statements {
			switch s := stmt.(type) {
			case *ast.AssignmentStatement:
				left, ok := s.Left.(*ast.Identifier)
				if ok && left.Value == node.Name.Value {
					// compile the return value
					err := c.Compile(s.Value)
					if err != nil {
						return err
					}

					// emit a OpReturnValue
					c.emit(code.OpReturnValue)
					hasReturnValue = true
				}
			default:
				err := c.Compile(stmt)
				if err != nil {
					return err
				}
			}
		}

		if !hasReturnValue {
			// emit a OpReturnValue with the default value (from the type)
			switch node.ReturnType.Type.Value {
			case "int":
				if node.ReturnType.IsArray {
					err := c.Compile(&ast.ArrayStatement{
						Elements: []ast.Expression{
							&ast.IntegerLiteral{Value: 0},
						},
					})
					if err != nil {
						return err
					}
				} else {
					err := c.Compile(&ast.IntegerLiteral{Value: 0})
					if err != nil {
						return err
					}
				}
			case "bool":
				if node.ReturnType.IsArray {
					err := c.Compile(&ast.ArrayStatement{
						Elements: []ast.Expression{
							&ast.Boolean{Value: false},
						},
					})
					if err != nil {
						return err
					}
				} else {
					err := c.Compile(&ast.Boolean{Value: false})
					if err != nil {
						return err
					}
				}
			case "float":
				if node.ReturnType.IsArray {
					err := c.Compile(&ast.ArrayStatement{
						Elements: []ast.Expression{
							&ast.FloatLiteral{Value: 0.0},
						},
					})
					if err != nil {
						return err
					}
				} else {
					err := c.Compile(&ast.FloatLiteral{Value: 0.0})
					if err != nil {
						return err
					}
				}
			default:
				return fmt.Errorf("unknown return type %s", node.ReturnType.Type.Value)
			}
			// emit a OpReturnValue
			c.emit(code.OpReturnValue)
		}

		if c.lastInstructionIs(code.OpPop) {
			c.replaceLastPopWithReturn()
		}
		if !c.lastInstructionIs(code.OpReturnValue) {
			c.emit(code.OpReturn)
		}

		numLocals := c.symbolTable.numDefinitions
		instructions := c.leaveScope()

		compiledFn := &object.CompiledFunction{
			Instructions:  instructions,
			NumLocals:     numLocals,
			NumParameters: len(node.Parameters),
		}
		c.emit(code.OpConstant, c.addConstant(compiledFn))

		// define the function name
		symbol := c.symbolTable.Define(node.Name.Value)
		c.emit(code.OpSetGlobal, symbol.Index)

		// pop the function
		// c.emit(code.OpPop)

	case *ast.CallExpression:
		err := c.Compile(node.Function)
		if err != nil {
			return err
		}

		for _, arg := range node.Arguments {
			err := c.Compile(arg)
			if err != nil {
				return err
			}
		}

		c.emit(code.OpCall, len(node.Arguments))

	// Types
	case *ast.ArrayStatement:
		for _, e := range node.Elements {
			err := c.Compile(e)
			if err != nil {
				return err
			}
		}
		c.emit(code.OpArray, len(node.Elements))
	case *ast.Boolean:
		if node.Value {
			c.emit(code.OpTrue)
		} else {
			c.emit(code.OpFalse)
		}
	case *ast.StringLiteral:
		string := &object.String{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(string))
	case *ast.FloatLiteral:
		float := &object.Float{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(float))
	case *ast.IntegerLiteral:
		integer := &object.Integer{Value: node.Value}
		c.emit(code.OpConstant, c.addConstant(integer))
	}
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instructions: c.currentInstructions(),
		Constants:    c.constants,
	}
}

func (c *Compiler) emit(op code.Opcode, operands ...int) int {
	ins := code.Make(op, operands...)
	pos := c.addInstruction(ins)

	c.setLastInstruction(op, pos)

	return pos
}

func (c *Compiler) currentInstructions() code.Instructions {
	return c.scopes[c.scopeIndex].instructions
}

func (c *Compiler) addInstruction(ins []byte) int {
	posNewInstruction := len(c.currentInstructions())
	updatedInstructions := append(c.currentInstructions(), ins...)

	c.scopes[c.scopeIndex].instructions = updatedInstructions

	return posNewInstruction
}

func (c *Compiler) addConstant(obj object.Object) int {
	c.constants = append(c.constants, obj)
	return len(c.constants) - 1
}

func (c *Compiler) setLastInstruction(op code.Opcode, pos int) {
	previous := c.scopes[c.scopeIndex].lastInstruction
	last := EmittedInstruction{Opcode: op, Position: pos}

	c.scopes[c.scopeIndex].previousInstruction = previous
	c.scopes[c.scopeIndex].lastInstruction = last
}

func (c *Compiler) lastInstructionIs(op code.Opcode) bool {
	if len(c.currentInstructions()) == 0 {
		return false
	}

	return c.scopes[c.scopeIndex].lastInstruction.Opcode == op
}

func (c *Compiler) removeLastPop() {
	last := c.scopes[c.scopeIndex].lastInstruction
	previous := c.scopes[c.scopeIndex].previousInstruction

	old := c.currentInstructions()
	new := old[:last.Position]

	c.scopes[c.scopeIndex].instructions = new
	c.scopes[c.scopeIndex].lastInstruction = previous
}

func (c *Compiler) replaceInstruction(pos int, newInstruction []byte) {
	instructions := c.currentInstructions()

	for i := 0; i < len(newInstruction); i++ {
		instructions[pos+i] = newInstruction[i]
	}
}

func (c *Compiler) changeOperand(opPos int, operand int) {
	op := code.Opcode(c.currentInstructions()[opPos])
	newInstruction := code.Make(op, operand)

	c.replaceInstruction(opPos, newInstruction)
}

func (c *Compiler) enterScope() {
	scope := CompilationScope{
		instructions:        code.Instructions{},
		lastInstruction:     EmittedInstruction{},
		previousInstruction: EmittedInstruction{},
	}
	c.scopes = append(c.scopes, scope)
	c.scopeIndex++
	c.symbolTable = NewEnclosedSymbolTable(c.symbolTable)
}

func (c *Compiler) leaveScope() code.Instructions {
	instructions := c.currentInstructions()

	c.scopes = c.scopes[:len(c.scopes)-1]
	c.scopeIndex--

	c.symbolTable = c.symbolTable.Outer

	return instructions
}

func (c *Compiler) replaceLastPopWithReturn() {
	lastPos := c.scopes[c.scopeIndex].lastInstruction.Position
	c.replaceInstruction(lastPos, code.Make(code.OpReturnValue))

	c.scopes[c.scopeIndex].lastInstruction.Opcode = code.OpReturnValue
}

func (c *Compiler) loadSymbol(s Symbol) {
	switch s.Scope {
	case GlobalScope:
		c.emit(code.OpGetGlobal, s.Index)
	case LocalScope:
		c.emit(code.OpGetLocal, s.Index)
	case BuiltinScope:
		c.emit(code.OpGetBuiltin, s.Index)
	}
}
