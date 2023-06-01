package compiler

import "github.com/odas0r/yail/object"

type SymbolScope string

const (
	GlobalScope  SymbolScope = "GLOBAL"
	LocalScope   SymbolScope = "LOCAL"
	BuiltinScope SymbolScope = "BUILTIN"
)

type Symbol struct {
	Name  string
	Scope SymbolScope
	Index int
}

type SymbolStruct struct {
	Name      string
	Scope     SymbolScope
	Index     int
	Atributes map[string]object.Object
}

type SymbolTable struct {
	Outer          *SymbolTable
	store          map[string]Symbol
	structs        map[string]SymbolStruct
	numDefinitions int
}

func NewSymbolTable() *SymbolTable {
	s := make(map[string]Symbol)
	structs := make(map[string]SymbolStruct)
	return &SymbolTable{store: s, structs: structs}
}

func NewEnclosedSymbolTable(outer *SymbolTable) *SymbolTable {
	s := NewSymbolTable()
	s.Outer = outer
	return s
}

func (s *SymbolTable) Define(name string) Symbol {
	symbol := Symbol{Name: name, Index: s.numDefinitions}
	if s.Outer == nil {
		symbol.Scope = GlobalScope
	} else {
		symbol.Scope = LocalScope
	}
	s.store[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) DefineStruct(name string) SymbolStruct {
	symbol := SymbolStruct{
		Name:      name,
		Index:     s.numDefinitions,
		Scope:     GlobalScope,
		Atributes: make(map[string]object.Object),
	}
	s.structs[name] = symbol
	s.numDefinitions++
	return symbol
}

func (s *SymbolTable) ResolveStruct(name string) (SymbolStruct, bool) {
	obj, ok := s.structs[name]
	return obj, ok
}

func (s *SymbolTable) DefineAttribute(structName string, name string, attr object.Object) SymbolStruct {
	var symbol SymbolStruct
	// find the struct
	sy, ok := s.structs[structName]
	if !ok {
		symbol = s.DefineStruct(structName)
	} else {
		symbol = sy
	}

	// define the attribute
	symbol.Atributes[name] = attr

	return symbol
}

func (s *SymbolTable) DefineBuiltin(index int, name string) Symbol {
	symbol := Symbol{Name: name, Index: index, Scope: BuiltinScope}
	s.store[name] = symbol
	return symbol
}

func (s *SymbolTable) Resolve(name string) (Symbol, bool) {
	obj, ok := s.store[name]
	if !ok && s.Outer != nil {
		obj, ok = s.Outer.Resolve(name)
		return obj, ok
	}
	return obj, ok
}
