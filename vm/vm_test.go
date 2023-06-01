package vm

import (
	"fmt"
	"testing"

	"github.com/odas0r/yail/ast"
	"github.com/odas0r/yail/compiler"
	"github.com/odas0r/yail/lexer"
	"github.com/odas0r/yail/object"
	"github.com/odas0r/yail/parser"
)

func parse(input string) *ast.Program {
	l := lexer.New(input)
	p := parser.New(l)
	return p.ParseProgram()
}

func testIntegerObject(expected int64, actual object.Object) error {
	result, ok := actual.(*object.Integer)
	if !ok {
		return fmt.Errorf("object is not Integer. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%d, want=%d", result.Value, expected)
	}

	return nil
}

func testFloatObject(expected float64, actual object.Object) error {
	result, ok := actual.(*object.Float)
	if !ok {
		return fmt.Errorf("object is not Float. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%f, want=%f", result.Value, expected)
	}

	return nil
}

func testBooleanObject(expected bool, actual object.Object) error {
	result, ok := actual.(*object.Boolean)
	if !ok {
		return fmt.Errorf("object is not Boolean. got=%T (%+v)", actual, actual)
	}

	if result.Value != expected {
		return fmt.Errorf("object has wrong value. got=%t, want=%t", result.Value, expected)
	}

	return nil
}

type vmTestCase struct {
	input    string
	expected interface{}
}

func runVmTests(t *testing.T, tests []vmTestCase) {
	t.Helper()

	for _, tt := range tests {
		program := parse(tt.input)

		comp := compiler.New()
		err := comp.Compile(program)
		if err != nil {
			t.Fatalf("compiler error: %s", err)
		}

		vm := New(comp.Bytecode())
		err = vm.Run()
		if err != nil {
			t.Fatalf("vm error: %s", err)
		}

		stackElem := vm.LastPoppedStackElem()

		testExpectedObject(t, tt.expected, stackElem)
	}
}

func testExpectedObject(
	t *testing.T,
	expected interface{},
	actual object.Object,
) {
	t.Helper()

	switch expected := expected.(type) {
	case int:
		err := testIntegerObject(int64(expected), actual)
		if err != nil {
			t.Errorf("testIntegerObject failed: %s", err)
		}
	case bool:
		err := testBooleanObject(expected, actual)
		if err != nil {
			t.Errorf("testBooleanObject failed: %s", err)
		}
	case *object.Null:
		if actual != Null {
			t.Errorf("object is not Null. got=%T (%+v)", actual, actual)
		}
	case string:
		if str, ok := actual.(*object.String); ok {
			if str.Value != expected {
				t.Errorf("String has wrong value. got=%q, want=%q", str.Value, expected)
			}
		} else {
			t.Errorf("object is not String. got=%T (%+v)", actual, actual)
		}
	case []int:
		array, ok := actual.(*object.Array)
		if !ok {
			t.Errorf("object is not Array. got=%T (%+v)", actual, actual)
		}

		if len(array.Elements) != len(expected) {
			t.Errorf("array has wrong num of elements. got=%d, want=%d", len(array.Elements), len(expected))
		}

		for i, expectedElem := range expected {
			err := testIntegerObject(int64(expectedElem), array.Elements[i])
			if err != nil {
				t.Errorf("testIntegerObject failed: %s", err)
			}
		}
	case [][]object.Object:
		// gotStruct, ok := actual.(*object.Struct)
		// if !ok {
		// 	t.Errorf("object is not Struct. got=%T (%+v)", actual, actual)
		// }
		// fmt.Println("STRUCT: " + actual.Inspect())
	}
}

func TestIntegerArithmetic(t *testing.T) {
	tests := []vmTestCase{
		{"1", 1},
		{"2", 2},
		{"1 + 2", 3},
		{"1 - 2", -1},
		{"1 * 2", 2},
		{"4 / 2", 2},
		{"50 / 2 * 2 + 10 - 5", 55},
		{"5 + 5 + 5 + 5 - 10", 10},
		{"2 * 2 * 2 * 2 * 2", 32},
		{"5 * 2 + 10", 20},
		{"5 + 2 * 10", 25},
		{"5 * (2 + 10)", 60},
		{"-5", -5},
		{"-10", -10},
		{"-50 + 100 + -50", 0},
		{"(5 + 10 * 2 + 15 / 3) * 2 + -10", 50},
	}
	runVmTests(t, tests)
}

func TestBooleanExpressions(t *testing.T) {
	tests := []vmTestCase{
		{"true", true},
		{"false", false},
		{"1 < 2", true},
		{"1 > 2", false},
		{"1 < 1", false},
		{"1 > 1", false},
		{"1 == 1", true},
		{"1 != 1", false},
		{"1 == 2", false},
		{"1 != 2", true},

		{"true == true", true},
		{"false == false", true},
		{"true == false", false},
		{"true != false", true},
		{"false != true", true},

		{"(1 < 2) == true", true},
		{"(1 < 2) == false", false},
		{"(1 > 2) == true", false},
		{"(1 > 2) == false", true},

		{"!true", false},
		{"!false", true},
		{"!5", false},
		{"!!true", true},
		{"!!false", false},
		{"!!5", true},
	}
	runVmTests(t, tests)
}

func TestConditionals(t *testing.T) {
	tests := []vmTestCase{
		{"if (true) { 10 }", 10},
		{"if (true) { 10 } else { 20 }", 10},
		{"if (false) { 10 } else { 20 }", 20},
		{"if (1) { 10 }", 10},
		{"if (1 < 2) { 10 }", 10},
		{"if (1 < 2) { 10 } else { 20 }", 10},
		{"if (1 > 2) { 10 } else { 20 }", 20},
		{"if (1 > 2) { 10 }", Null},
		{"if (false) { 10 }", Null},
		{"!(if (false) { 5; })", true},
		{"if ((if (false) { 10 })) { 10 } else { 20 }", 20},
	}
	runVmTests(t, tests)
}

func TestStringExpressions(t *testing.T) {
	tests := []vmTestCase{
		{`"monkey"`, "monkey"},
		{`"mon" + "key"`, "monkey"},
		{`"mon" + "key" + "banana"`, "monkeybanana"},
	}
	runVmTests(t, tests)
}

func TestArrayStatements(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			global {
				int v[];
			}
			v;
			`, []int{0},
		},
		{
			`
			global {
				int v[] = {1,2,3};
			}
			v;
			`, []int{1, 2, 3},
		},
		{
			`
			global {
				int v[] = {1+2,3-4,5*6};
			}
			v;
			`, []int{3, -1, 30},
		},
	}
	runVmTests(t, tests)
}

func TestGlobalVarStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			global {
				int one = 1;
			}
			one;
			`, 1,
		},
		{
			`
			global {
				int one = 1;
				int two = 2;
			}
			one + two;
			`, 3,
		},
		{
			`
			global {
				int one = 1;
				int two = one + one;
			}
			one + two;
			`, 3,
		},
	}
	runVmTests(t, tests)
}

func TestIndexStatement(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			global {
				int v[3] = {1,2,3};
			}
			v[1]
			`, 2,
		},
		{
			`
			global {
				int v[3] = {1,2,3};
			}
			v[0 + 2]
			`, 3,
		},
		{
			`
			global {
				int v[];
			}
			v[0]
			`, 0,
		},
		{
			`
			global {
				int v[] = {1,2,3};
			}
			v[99]
			`, Null,
		},
		{
			`
			global {
				int v[] = {1,2,3};
			}
			v[-1]
			`, Null,
		},
	}
	runVmTests(t, tests)
}

// TODO: finish this later
// func TestStructStatement(t *testing.T) {
// 	tests := []vmTestCase{
// 		{
// 			`
// 			structs {
// 				point2D {int x;};
// 			}
// 			point2D p;
// 			`, 0,
// 		},
// 		{
// 			`
// 			structs {
// 				circle {int center, int radius;};
// 				point3D {float x, y, z;};
// 			}
// 			circle c;
// 			`, 0,
// 		},
// 	}
// 	runVmTests(t, tests)
// }

func TestCallingFunctionsWithoutArguments(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			add() int {
				add = 5 + 10;
			}
			add();
			`, 15,
		},
		{
			`
			add() int {
				add = 5 + 10;
				# FIXME: DOESN'T WORK
			  # add = add + 15;
			}
			add();
			`, 15,
		},
		{
			`
			a() int {
				a = 1;
			}
			b() int {
				b = a() + 1;
			}
			c() int {
				c = b() + 1;
			}
			c();
			`, 3,
		},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithoutReturn(t *testing.T) {
	tests := []vmTestCase{
		{
			`
			noReturn() int { }
			noReturn();
			`, 0,
		},
		{
			`
			noReturn() int[] { }
			noReturn();
			`, []int{0},
		},
		{
			`
			noReturn() int { }
			noReturnTwo() int { noReturnTwo = noReturn(); }
			noReturn();
			noReturnTwo();
			`, 0,
		},
		{
			`
			noReturn() bool { }
			noReturnTwo() int { noReturnTwo = noReturn(); }
			noReturn();
			noReturnTwo();
			`, false,
		},
		{
			`
			noReturn() float { }
			noReturnTwo() int { noReturnTwo = noReturn(); }
			noReturn();
			noReturnTwo();
			`, 0.0,
		},
	}
	runVmTests(t, tests)
}

func TestCallingFunctionsWithBindings(t *testing.T) {
	tests := []vmTestCase{

		{
			`
			wow() int {
				local {
					int a = 30;
					int b = 30;
				}
				wow = a + b;
			}
			wow();
			`, 60,
		},
		{
			`
			oneAndTwo() int {
				local {
					int one = 1;
					int two = 2;
				}
				oneAndTwo = one + two;
			}
			threeAndFour() int {
				local {
					int three = 3;
					int four = 4;
				}
				threeAndFour = three + four;
			}
			oneAndTwo() + threeAndFour();
			`, 10,
		},
		{
			`
			global {
				int globalSeed = 50;
			}

			minusOne() int {
				local {
					int num = 1;
				}
				minusOne = globalSeed - num;
			}

			minusTwo() int {
				local {
					int num = 2;
				}
				minusTwo = globalSeed - num;
			}
			minusOne() + minusTwo();
			`, 97,
		},
	}
	runVmTests(t, tests)
}
