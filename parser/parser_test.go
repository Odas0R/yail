package parser

import (
	"fmt"
	"testing"

	"github.com/odas0r/yail/ast"
	"github.com/odas0r/yail/lexer"
)

func TestVariableStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedType  string
		expectedName  string
		expectedValue interface{}
	}{
		{"int x;", "int", "x", 0},
		{"bool y;", "bool", "y", false},
		{"float z;", "float", "z", 0.0},
		{"int j = 5;", "int", "j", 5},
		{"bool k = true;", "bool", "k", true},
		{"float l = 5.24;", "float", "l", 5.24},
		{"bool foobar = y;", "bool", "foobar", "y"},
		{"foobar = 123;", "", "foobar", 123},
		{"foobar = \"123\";", "", "foobar", "123"},
		{"foobar = false;", "", "foobar", false},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}
		stmt := program.Statements[0]
		if !testVariableStatement(t, stmt, tt.expectedType, tt.expectedName) {
			return
		}

		val := stmt.(*ast.VariableStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestVectorStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedType  string
		expectedName  string
		expectedSize  interface{}
		expectedValue interface{}
	}{
		{"int a[];", "int", "a", 1, []int64{0}},
		{"bool b[];", "bool", "b", 1, []bool{false}},
		{"float c[];", "float", "c", 1, []float64{0}},
		{"int x[3]={1, 2, 3};", "int", "x", 3, []int64{1, 2, 3}},
		{"float y[2]={1.2, 2.3};", "float", "y", 2, []float64{1.2, 2.3}},
		{"int k[]={1,2,3,4,5};", "int", "k", 5, []int64{1, 2, 3, 4, 5}},
		{"j={1,2,3,4,5};", "", "j", 5, []interface{}{1, 2, 3, 4, 5}},
		{"z={\"a\",\"b\"};", "", "z", 2, []interface{}{"a", "b"}},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain 1 statements. got=%d",
				len(program.Statements))
		}

		stmt := program.Statements[0]

		if !testVectorStatement(t, stmt, tt.expectedType, tt.expectedName, tt.expectedSize) {
			return
		}

		for i, v := range stmt.(*ast.VectorStatement).Values {
			switch tt.expectedValue.(type) {
			case []int64:
				testLiteralExpression(t, v, tt.expectedValue.([]int64)[i])
			case []float64:
				testLiteralExpression(t, v, tt.expectedValue.([]float64)[i])
			case []bool:
				testLiteralExpression(t, v, tt.expectedValue.([]bool)[i])
			default:
				testLiteralExpression(t, v, tt.expectedValue.([]interface{})[i])
			}
		}
	}
}

func TestParsingVectorIndexSetterExpressions(t *testing.T) {
	input := "x[0] = 5.4;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	accessorExp, ok := stmt.Left.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.AccessorExpression. got=%T", stmt.Left)
	}
	if !testIdentifier(t, accessorExp.Left, "x") {
		return
	}
	if !testLiteralExpression(t, accessorExp.Index, 0) {
		return
	}
	if !testLiteralExpression(t, stmt.Value, 5.4) {
		return
	}
}

func TestParsingIndexExpressions(t *testing.T) {
	input := "myArray[1 + 1]"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)

	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}

	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

func TestStructsDeclaration(t *testing.T) {
	input := `
	structs {
		point2D { float x, float y; };
		point3D { float x, y, z; };
		point4D { float x, y, z, int j; };
		pointND { float x[]; };
		pointNDSize { float x[5]; };
		pointNDSizeM { float x[5], y[2], z[]; };
	}
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.StructsStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.StructDefinition. got=%T", program.Statements[0])
	}

	if len(stmt.Structs) != 6 {
		t.Fatalf("stmt.Structs has wrong length. got=%d", len(stmt.Structs))
	}

	expectedStructs := []struct {
		Name       string
		Attributes []struct {
			Type     string
			Name     string
			Size     interface{}
			IsVector bool
		}
	}{
		{
			Name: "point2D",
			Attributes: []struct {
				Type     string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Type: "float",
					Name: "x",
				},
				{
					Type: "float",
					Name: "y",
				},
			},
		},
		{
			Name: "point3D",
			Attributes: []struct {
				Type     string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Type: "float",
					Name: "x",
				},
				{
					Type: "float",
					Name: "y",
				},
				{
					Type: "float",
					Name: "z",
				},
			},
		},
		{
			Name: "point4D",
			Attributes: []struct {
				Type     string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Type: "float",
					Name: "x",
				},
				{
					Type: "float",
					Name: "y",
				},
				{
					Type: "float",
					Name: "z",
				},
				{
					Type: "int",
					Name: "j",
				},
			},
		},
		{
			Name: "pointND",
			Attributes: []struct {
				Type     string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Type:     "float",
					Name:     "x",
					Size:     1,
					IsVector: true,
				},
			},
		},
		{
			Name: "pointNDSize",
			Attributes: []struct {
				Type     string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Type:     "float",
					Name:     "x",
					Size:     "5",
					IsVector: true,
				},
			},
		},
		{
			Name: "pointNDSizeM",
			Attributes: []struct {
				Type     string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Type:     "float",
					Name:     "x",
					Size:     "5",
					IsVector: true,
				},
				{
					Type:     "float",
					Name:     "y",
					Size:     "2",
					IsVector: true,
				},
				{
					Type:     "float",
					Name:     "z",
					Size:     1,
					IsVector: true,
				},
			},
		},
	}

	for i, str := range stmt.Structs {
		expectedStruct := expectedStructs[i]

		if str.TokenLiteral() != expectedStruct.Name {
			t.Errorf("str.Name.Value not %s. got=%s", expectedStructs[i].Name, str.TokenLiteral())
		}

		for j, a := range str.Attributes {
			attr := expectedStruct.Attributes[j]

			if a.Type.Value != attr.Type {
				t.Errorf("a.Type.Value not %s. got=%s", attr.Type, a.Type.Value)
			}

			if a.Name.Value != attr.Name {
				t.Errorf("p.Name.Value not %s. got=%s", attr.Name, a.Name.Value)
			}

			// only check size if it is not nil
			if a.Size != nil {
				if a.Size.String() != fmt.Sprintf("%v", attr.Size) {
					t.Errorf("p.Size not %v. got=%v", attr.Size, a.Size)
				}
			}

			if a.IsVector != attr.IsVector {
				t.Errorf("p.IsVector not %t. got=%t", attr.IsVector, a.IsVector)
			}
		}
	}
}

func TestParsingAccessorExpressions(t *testing.T) {
	input := "point2D.x"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, _ := program.Statements[0].(*ast.ExpressionStatement)

	accessorExp, ok := stmt.Expression.(*ast.AccessorExpression)
	if !ok {
		t.Fatalf("exp not *ast.AccessorExpression. got=%T", stmt.Expression)
	}
	if !testIdentifier(t, accessorExp.Left, "point2D") {
		return
	}
	if !testIdentifier(t, accessorExp.Index[0], "x") {
		return
	}
}

func TestParsingStructAccessorSetterExpressions(t *testing.T) {
	input := "point2D.x = 5.4;"
	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)
	stmt, ok := program.Statements[0].(*ast.AssignmentStatement)
	if !ok {
		t.Fatalf("stmt not *ast.AssignmentStatement. got=%T", program.Statements[0])
	}

	accessorExp, ok := stmt.Left.(*ast.AccessorExpression)
	if !ok {
		t.Fatalf("exp not *ast.AccessorExpression. got=%T", stmt.Left)
	}
	if !testIdentifier(t, accessorExp.Left, "point2D") {
		return
	}
	if !testIdentifier(t, accessorExp.Index[0], "x") {
		return
	}
	if !testLiteralExpression(t, stmt.Value, 5.4) {
		return
	}
}

func TestGlobalStatement(t *testing.T) {
	input := `global {
	int x = 5;
	float y = 5.5;
	x = {1, 2, 3};
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.GlobalStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	tt := []struct {
		expectedType  string
		expectedName  string
		expectedValue interface{}
		expectedSize  interface{}
	}{
		{
			expectedType:  "int",
			expectedName:  "x",
			expectedValue: 5,
		},
		{
			expectedType:  "float",
			expectedName:  "y",
			expectedValue: 5.5,
		},
		{
			expectedType:  "",
			expectedName:  "x",
			expectedValue: []int{1, 2, 3},
			expectedSize:  3,
		},
	}

	for i, gv := range stmt.Body.Statements {
		switch gv.(type) {
		case *ast.VariableStatement:
			testVariableStatement(t, gv, tt[i].expectedType, tt[i].expectedName)
		case *ast.VectorStatement:
			testVectorStatement(t, gv, tt[i].expectedType, tt[i].expectedName, tt[i].expectedSize)
		default:
			t.Fatalf("gv is not ast.VariableStatement or ast.VectorStatement. got=%T", gv)
		}
	}
}

func TestGlobalErrorStatement(t *testing.T) {
	input := `global {
	int x = 5;
	add(int x,y) int {
		add = x + y;
	}
}`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	expectedErrors := []string{
		"only variable declarations are allowed in variable blocks",
	}

	testParserErrors(t, p, expectedErrors, false)
}

func TestConstStatement(t *testing.T) {
	input := `const {
	int x = 5;
	float y = 5.5;
	x = {1, 2, 3};
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ConstStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	tt := []struct {
		expectedType  string
		expectedName  string
		expectedValue interface{}
		expectedSize  interface{}
	}{
		{
			expectedType:  "int",
			expectedName:  "x",
			expectedValue: 5,
		},
		{
			expectedType:  "float",
			expectedName:  "y",
			expectedValue: 5.5,
		},
		{
			expectedType:  "",
			expectedName:  "x",
			expectedValue: []int{1, 2, 3},
			expectedSize:  3,
		},
	}

	for i, gv := range stmt.Body.Statements {
		switch gv.(type) {
		case *ast.VariableStatement:
			testVariableStatement(t, gv, tt[i].expectedType, tt[i].expectedName)
		case *ast.VectorStatement:
			testVectorStatement(t, gv, tt[i].expectedType, tt[i].expectedName, tt[i].expectedSize)
		default:
			t.Fatalf("gv is not ast.VariableStatement or ast.VectorStatement. got=%T", gv)
		}
	}
}

func TestConstErrorStatement(t *testing.T) {
	input := `const {
	int x = 5;
	add(int x,y) int { add = x + y; }
}`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	expectedErrors := []string{
		"only variable declarations are allowed in variable blocks",
	}

	testParserErrors(t, p, expectedErrors, false)
}

func TestLocalStatement(t *testing.T) {
	input := `add() int {
	local {
		int x = 5;
	}

	add = x;
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	tt := []struct {
		expectedType  string
		expectedName  string
		expectedValue interface{}
		expectedSize  interface{}
	}{
		{
			expectedType:  "int",
			expectedName:  "x",
			expectedValue: 5,
		},
	}

	ls, ok := stmt.Body.Statements[0].(*ast.LocalStatement)
	if !ok {
		t.Fatalf("stmt.Body.Statements[0] is not ast.LocalStatement. got=%T",
			stmt.Body.Statements[0])
	}

	for i, s := range ls.Body.Statements {
		switch s.(type) {
		case *ast.VariableStatement:
			testVariableStatement(t, s, tt[i].expectedType, tt[i].expectedName)
		case *ast.VectorStatement:
			testVectorStatement(t, s, tt[i].expectedType, tt[i].expectedName, tt[i].expectedSize)
		default:
			t.Fatalf("gv is not ast.VariableStatement or ast.VectorStatement. got=%T", s)
		}
	}
}

func TestLocalErrorStatement(t *testing.T) {
	input := `local {
	int x = 5;
}`

	l := lexer.New(input)
	p := New(l)
	p.ParseProgram()

	expectedErrors := []string{
		"no prefix parse function for LOCAL found",
		"no prefix parse function for { found",
		"no prefix parse function for } found",
	}

	testParserErrors(t, p, expectedErrors, false)
}

func TestCommentStatement(t *testing.T) {
	input := `
	# this is a very interesting comment wow
`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 0 {
		t.Fatalf("program has statements. got=%d",
			len(program.Statements))
	}
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}

	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != 5 {
		t.Errorf("ident.Value not %d. got=%d", 5, ident.Value)
	}
	if ident.TokenLiteral() != "5" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "5", ident.TokenLiteral())
	}
}

func TestFloatLiteralExpression(t *testing.T) {
	input := "5.5;"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program has not enough statements. got=%d",
			len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	ident, ok := stmt.Expression.(*ast.FloatLiteral)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	if ident.Value != 5.5 {
		t.Errorf("ident.Value not %f. got=%f", 5.5, ident.Value)
	}
	if ident.TokenLiteral() != "5.5" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "5.5", ident.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}

	for _, tt := range prefixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("exp is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		if !testLiteralExpression(t, exp.Right, tt.integerValue) {
			return
		}
	}
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}

	for _, tt := range infixTests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		if len(program.Statements) != 1 {
			t.Fatalf("program.Statements does not contain %d statements. got=%d\n", 1, len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}

		if !testInfixExpression(t, stmt.Expression, tt.leftValue,
			tt.operator, tt.rightValue) {
			return
		}
	}
}

func TestOperatorPrecedenceParsing(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"-a * b",
			"((-a) * b)",
		},
		{
			"!-a",
			"(!(-a))",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4; -5 * 5",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"true",
			"true",
		},
		{
			"false",
			"false",
		},
		{
			"3 > 5 == false",
			"((3 > 5) == false)",
		},
		{
			"3 < 5 == true",
			"((3 < 5) == true)",
		},
		{
			"1 + (2 + 3) + 4",
			"((1 + (2 + 3)) + 4)",
		},
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 / (5 + 5)",
			"(2 / (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
		{
			"a + add(b * c) + d",
			"((a + add((b * c))) + d)",
		},
		{
			"add(a, b, 1, 2 * 3, 4 + 5, add(6, 7 * 8))",
			"add(a, b, 1, (2 * 3), (4 + 5), add(6, (7 * 8)))",
		},
		{
			"add(a + b + c * d / f + g)",
			"add((((a + b) + ((c * d) / f)) + g))",
		},
		{
			"a * a[b * c] * d",
			"((a * (a[(b * c)])) * d)",
		},
		{
			"add(a * b[2], b[1], 2 * a[1])",
			"add((a * (b[2])), (b[1]), (2 * (a[1])))",
		},
		{
			"a * (z.a) * d",
			"((a * (z.a)) * d)",
		},
		{
			"add(a * (z.a.b) * (d.a))",
			"add(((a * (z.a.b)) * (d.a)))",
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)

		program := p.ParseProgram()
		checkParserErrors(t, p)

		actual := program.String()
		if actual != tt.expected {
			t.Errorf("expected=%q, got=%q", tt.expected, actual)
		}
	}
}

func TestIfExpression(t *testing.T) {
	input := `if (x < y) { x }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if exp.Alternative != nil {
		t.Errorf("exp.Alternative.Statements was not nil. got=%+v", exp.Alternative)
	}
}

func TestIfElseExpression(t *testing.T) {
	input := `if (x < y) { x } else { y }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}
	exp, ok := stmt.Expression.(*ast.IfExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.IfExpression. got=%T",
			stmt.Expression)
	}

	if !testInfixExpression(t, exp.Condition, "x", "<", "y") {
		return
	}

	if len(exp.Consequence.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(exp.Consequence.Statements))
	}

	consequence, ok := exp.Consequence.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Consequence.Statements[0])
	}

	if !testIdentifier(t, consequence.Expression, "x") {
		return
	}

	if len(exp.Alternative.Statements) != 1 {
		t.Errorf("exp.Alternative.Statements does not contain 1 statements. got=%d\n",
			len(exp.Alternative.Statements))
	}

	alternative, ok := exp.Alternative.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.ExpressionStatement. got=%T",
			exp.Alternative.Statements[0])
	}

	if !testIdentifier(t, alternative.Expression, "y") {
		return
	}
}

func TestWhileStatement(t *testing.T) {
	input := `while (x < y) { x = 5; }`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	whileStmt, ok := program.Statements[0].(*ast.WhileStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.WhileStatement. got=%T",
			program.Statements[0])
	}

	if !testInfixExpression(t, whileStmt.Condition, "x", "<", "y") {
		return
	}

	if len(whileStmt.Body.Statements) != 1 {
		t.Errorf("consequence is not 1 statements. got=%d\n",
			len(whileStmt.Body.Statements))
	}

	consequence, ok := whileStmt.Body.Statements[0].(*ast.VariableStatement)
	if !ok {
		t.Fatalf("Statements[0] is not ast.VariableStatement. got=%T",
			whileStmt.Body.Statements[0])
	}

	if !testIdentifier(t, consequence.Name, "x") {
		return
	}
}

func TestFunctionParameterParsing(t *testing.T) {
	tests := []struct {
		input          string
		expectedParams []struct {
			Name string
			Type string
		}
	}{
		{input: "fn() int {}", expectedParams: []struct {
			Name string
			Type string
		}{}},
		{input: "fn(int x) int {}", expectedParams: []struct {
			Name string
			Type string
		}{
			{
				Name: "x",
				Type: "int",
			},
		}},
		{input: "fn(int x, y, z) int {}", expectedParams: []struct {
			Name string
			Type string
		}{
			{
				Name: "x",
				Type: "int",
			},
			{
				Name: "y",
				Type: "int",
			},
			{
				Name: "z",
				Type: "int",
			},
		},
		},
	}

	for _, tt := range tests {
		l := lexer.New(tt.input)
		p := New(l)
		program := p.ParseProgram()
		checkParserErrors(t, p)

		function, ok := program.Statements[0].(*ast.FunctionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
				program.Statements[0])
		}

		if len(function.Parameters) != len(tt.expectedParams) {
			t.Fatalf("length parameters wrong. want %d, got=%d\n",
				len(tt.expectedParams), len(function.Parameters))
		}

		for i, p := range tt.expectedParams {
			testTokenType(t, function.Parameters[i].Type, p.Type)
			testLiteralExpression(t, function.Parameters[i].Name, p.Name)
		}
	}
}

func TestFunctionLiteralParsing(t *testing.T) {
	input := `add(int x, int y) int {
	x + y;
}`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}
	function, ok := program.Statements[0].(*ast.FunctionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.FunctionStatement. got=%T",
			program.Statements[0])
	}

	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
			len(function.Parameters))
	}

	testLiteralExpression(t, function.Parameters[0].Name, "x")
	testTokenType(t, function.Parameters[0].Type, "int")

	testLiteralExpression(t, function.Parameters[1].Name, "y")
	testTokenType(t, function.Parameters[1].Type, "int")

	testTokenType(t, function.ReturnType.Type, "int")
	if function.ReturnType.IsVector {
		t.Errorf("function.ReturnType.IsVector is not false. got=%t", function.ReturnType.IsVector)
	}

	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
			len(function.Body.Statements))
	}

	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
			function.Body.Statements[0])
	}

	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

func TestCallExpressionParsing(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
			1, len(program.Statements))
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
			program.Statements[0])
	}

	exp, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
			stmt.Expression)
	}

	if !testIdentifier(t, exp.Function, "add") {
		return
	}
	if len(exp.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
	}

	testLiteralExpression(t, exp.Arguments[0], 1)
	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
}

func TestStringLiteralExpression(t *testing.T) {
	input := `"hello world";`

	l := lexer.New(input)
	p := New(l)
	program := p.ParseProgram()
	checkParserErrors(t, p)

	stmt := program.Statements[0].(*ast.ExpressionStatement)
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}

	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %q. got=%q", "hello world", literal.Value)
	}
}

// ---------------------- helpers ----------------------

func testVariableStatement(t *testing.T, s ast.Statement, typ string, name string) bool {
	varStmt, ok := s.(*ast.VariableStatement)
	if !ok {
		t.Errorf("s not *ast.VariableDeclaration. got=%T", s)
		return false
	}

	if varStmt.Type.Value != typ {
		t.Errorf("varStmt.Type.Value not '%s'. got=%q", typ, varStmt.Type.Value)
		return false
	}

	if varStmt.Type.Value == "" {
		if varStmt.Type.TokenLiteral() != name {
			t.Errorf("vecStmt.Type.TokenLiteral() not '%s'. got=%q", name, varStmt.Type.TokenLiteral())
			return false
		}
	} else {
		if varStmt.Type.TokenLiteral() != typ {
			t.Errorf("vecStmt.Type.TokenLiteral() not '%s'. got=%q", typ, varStmt.Type.TokenLiteral())
			return false
		}
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("varStmt.Name.TokenLiteral() not '%s'. got=%s", name, varStmt.Name)
		return false
	}

	return true
}

func testVectorStatement(t *testing.T, s ast.Statement, typ string, name string, size interface{}) bool {
	vecStmt, ok := s.(*ast.VectorStatement)
	if !ok {
		t.Errorf("statment not *ast.VectorStatement. got=%T", s)
		return false
	}

	if fmt.Sprintf("%v", vecStmt.Size) != fmt.Sprintf("%v", size) {
		t.Errorf("vecStmt.Size not %v. got=%v", size, vecStmt.Size)
		return false
	}

	if vecStmt.Type.Value != typ {
		t.Errorf("vecStmt.Type.Value not '%s'. got=%q", typ, vecStmt.Type.Value)
		return false
	}

	if vecStmt.Type.Value == "" {
		if vecStmt.Type.TokenLiteral() != name {
			t.Errorf("vecStmt.Type.TokenLiteral() not '%s'. got=%q", name, vecStmt.Type.TokenLiteral())
			return false
		}
	} else {
		if vecStmt.Type.TokenLiteral() != typ {
			t.Errorf("vecStmt.Type.TokenLiteral() not '%s'. got=%q", typ, vecStmt.Type.TokenLiteral())
			return false
		}
	}

	if vecStmt.Name.Value != name {
		t.Errorf("vecStmt.Name.Value not '%s'. got=%s", name, vecStmt.Name.Value)
		return false
	}

	if vecStmt.Name.TokenLiteral() != name {
		t.Errorf("vecStmt.Name not '%s'. got=%s", name, vecStmt.Name)
		return false
	}

	return true
}

func testTokenType(t *testing.T, tt *ast.Identifier, expected string) bool {
	if tt.Value != expected {
		t.Errorf("tt.Literal not '%s'. got=%s", expected, tt.Value)
		return false
	}
	return false
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case float64:
		return testFloatLiteral(t, exp, v)
	case string:
		if _, ok := exp.(*ast.StringLiteral); ok {
			return testStringLiteral(t, exp, v)
		} else if _, ok := exp.(*ast.Identifier); ok {
			return testIdentifier(t, exp, v)
		} else {
			t.Errorf("type of exp not handled. got=%T", exp)
			return false
		}
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}

	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}

	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}

	return true
}

func testFloatLiteral(t *testing.T, fl ast.Expression, value float64) bool {
	flo, ok := fl.(*ast.FloatLiteral)
	if !ok {
		t.Errorf("fl not *ast.FloatLiteral. got=%T", fl)
		return false
	}

	if flo.Value != value {
		t.Errorf("flo.Value not %f. got=%f", value, flo.Value)
		return false
	}
	if flo.TokenLiteral() != fmt.Sprintf("%.6g", value) {
		t.Errorf("flo.TokenLiteral not %f. got=%s", value, flo.TokenLiteral())
		return false
	}

	return true
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{},
	operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}

	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}

	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}

	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}

	return true
}

func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}

	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}

	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}

	return true
}

func testStringLiteral(t *testing.T, exp ast.Expression, value string) bool {
	string, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if string.Value != value {
		t.Errorf("string.Value not %s. got=%s", value, string.Value)
		return false
	}

	if string.TokenLiteral() != value {
		t.Errorf("string.TokenLiteral not %s. got=%s", value, string.TokenLiteral())
		return false
	}

	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}

	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s",
			value, bo.TokenLiteral())
		return false
	}

	return true
}

func testParserErrors(t *testing.T, p *Parser, expectedErrors []string, debug bool) {
	errors := p.Errors()

	if debug {
		for _, err := range p.Errors() {
			fmt.Println(err)
		}
	}

	if len(errors) != len(expectedErrors) {
		t.Fatalf("wrong number of errors. expected=%d, got=%d", len(expectedErrors), len(errors))
	}

	for _, wantErr := range expectedErrors {
		found := false
		for _, gotErr := range errors {
			if gotErr == wantErr {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("expected error not found: \"%s\"", wantErr)
		}
	}

}

func checkParserErrors(t *testing.T, p *Parser) {
	errors := p.Errors()
	if len(errors) == 0 {
		return
	}

	t.Errorf("parser has %d errors", len(errors))
	for _, msg := range errors {
		t.Errorf("parser error: %q", msg)
	}
	t.FailNow()
}
