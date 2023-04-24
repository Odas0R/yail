package parser

import (
	"fmt"
	"testing"

	"github.com/odas0r/yail/ast"
	"github.com/odas0r/yail/lexer"
)

func TestVarDeclarationStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedToken string
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
		if !testVarDeclaration(t, stmt, tt.expectedToken, tt.expectedName) {
			return
		}

		val := stmt.(*ast.VarStatement).Value
		if !testLiteralExpression(t, val, tt.expectedValue) {
			return
		}
	}
}

func TestVectorDeclarationStatements(t *testing.T) {
	tests := []struct {
		input         string
		expectedToken string
		expectedName  string
		expectedSize  int64
		expectedValue interface{}
	}{
		{"int a[];", "int", "a", 1, []int64{0}},
		{"bool b[];", "bool", "b", 1, []bool{false}},
		{"float c[];", "float", "c", 1, []float64{0}},
		{"int x[3]={1, 2, 3};", "int", "x", 3, []int64{1, 2, 3}},
		{"float y[2]={1.2, 2.3};", "float", "y", 2, []float64{1.2, 2.3}},
		{"int k[]={1,2,3,4,5};", "int", "k", 5, []int64{1, 2, 3, 4, 5}},
		// {"j={1,2,3,4,5};", "int", "j", 5, []int64{1, 2, 3, 4, 5}},
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

		if stmt.TokenLiteral() != tt.expectedToken {
			t.Errorf("s.TokenLiteral not '%s'. got=%q", tt.expectedToken, stmt.TokenLiteral())
		}

		vecStmt, ok := stmt.(*ast.VectorStatement)
		if !ok {
			t.Errorf("vecStmt not *ast.VectorDeclaration. got=%T", vecStmt)
		}

		if vecStmt.Name.Value != tt.expectedName {
			t.Errorf("vecStmt.Name.Value not '%s'. got=%s", tt.expectedName, vecStmt.Name.Value)
		}

		// test the vector values
		if len(vecStmt.Values) != int(tt.expectedSize) {
			t.Errorf("vecStmt.Values has wrong length. got=%d", len(vecStmt.Values))
		}

		// Check the values
		switch tt.expectedToken {
		case "int":
			expected := tt.expectedValue.([]int64)
			for i, v := range vecStmt.Values {
				val := v.(*ast.IntegerLiteral)
				if val.Value != expected[i] {
					t.Errorf("vecStmt.Values[%d] not '%d'. got=%d", i, expected[i], val.Value)
				}
			}
		case "bool":
			expected := tt.expectedValue.([]bool)
			for i, v := range vecStmt.Values {
				val := v.(*ast.Boolean)
				if val.Value != expected[i] {
					t.Errorf("vecStmt.Values[%d] not '%t'. got=%t", i, expected[i], val.Value)
				}
			}
		case "float":
			expected := tt.expectedValue.([]float64)
			for i, v := range vecStmt.Values {
				val := v.(*ast.FloatLiteral)
				if val.Value != expected[i] {
					t.Errorf("vecStmt.Values[%d] not '%f'. got=%f", i, expected[i], val.Value)
				}
			}
		default:
			t.Errorf("vecStmt.Values has wrong type. got=%T", vecStmt.Values)
		}
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

	stmt, ok := program.Statements[0].(*ast.StructsDefinition)
	if !ok {
		t.Fatalf("program.Statements[0] is not *ast.StructDefinition. got=%T", program.Statements[0])
	}

	if len(stmt.Structs) != 6 {
		t.Fatalf("stmt.Structs has wrong length. got=%d", len(stmt.Structs))
	}

	expectedStructs := []struct {
		Name       string
		Attributes []struct {
			Token    string
			Name     string
			Size     interface{}
			IsVector bool
		}
	}{
		{
			Name: "point2D",
			Attributes: []struct {
				Token    string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Token: "float",
					Name:  "x",
				},
				{
					Token: "float",
					Name:  "y",
				},
			},
		},
		{
			Name: "point3D",
			Attributes: []struct {
				Token    string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Token: "float",
					Name:  "x",
				},
				{
					Token: "float",
					Name:  "y",
				},
				{
					Token: "float",
					Name:  "z",
				},
			},
		},
		{
			Name: "point4D",
			Attributes: []struct {
				Token    string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Token: "float",
					Name:  "x",
				},
				{
					Token: "float",
					Name:  "y",
				},
				{
					Token: "float",
					Name:  "z",
				},
				{
					Token: "int",
					Name:  "j",
				},
			},
		},
		{
			Name: "pointND",
			Attributes: []struct {
				Token    string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Token:    "float",
					Name:     "x",
					Size:     1,
					IsVector: true,
				},
			},
		},
		{
			Name: "pointNDSize",
			Attributes: []struct {
				Token    string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Token:    "float",
					Name:     "x",
					Size:     "5",
					IsVector: true,
				},
			},
		},
		{
			Name: "pointNDSizeM",
			Attributes: []struct {
				Token    string
				Name     string
				Size     interface{}
				IsVector bool
			}{
				{
					Token:    "float",
					Name:     "x",
					Size:     "5",
					IsVector: true,
				},
				{
					Token:    "float",
					Name:     "y",
					Size:     "2",
					IsVector: true,
				},
				{
					Token:    "float",
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

			if a.TokenLiteral() != attr.Token {
				t.Errorf("p.Token.Literal not %s. got=%s", attr.Token, a.Token.Literal)
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

// func TestFunctionParameterParsing(t *testing.T) {
// 	tests := []struct {
// 		input          string
// 		expectedParams []string
// 	}{
// 		{input: "fn() {};", expectedParams: []string{}},
// 		{input: "fn(x) {};", expectedParams: []string{"x"}},
// 		{input: "fn(x, y, z) {};", expectedParams: []string{"x", "y", "z"}},
// 	}
//
// 	for _, tt := range tests {
// 		l := lexer.New(tt.input)
// 		p := New(l)
// 		program := p.ParseProgram()
// 		checkParserErrors(t, p)
//
// 		stmt := program.Statements[0].(*ast.ExpressionStatement)
// 		function := stmt.Expression.(*ast.FunctionLiteral)
//
// 		if len(function.Parameters) != len(tt.expectedParams) {
// 			t.Fatalf("length parameters wrong. want %d, got=%d\n",
// 				len(tt.expectedParams), len(function.Parameters))
// 		}
//
// 		for i, ident := range tt.expectedParams {
// 			testLiteralExpression(t, function.Parameters[i], ident)
// 		}
// 	}
// }
//
// func TestFunctionLiteralParsing(t *testing.T) {
// 	input := `fn(x, y) { x + y; }`
//
// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)
//
// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
// 			1, len(program.Statements))
// 	}
// 	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T",
// 			program.Statements[0])
// 	}
//
// 	function, ok := stmt.Expression.(*ast.FunctionLiteral)
// 	if !ok {
// 		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T",
// 			stmt.Expression)
// 	}
//
// 	if len(function.Parameters) != 2 {
// 		t.Fatalf("function literal parameters wrong. want 2, got=%d\n",
// 			len(function.Parameters))
// 	}
//
// 	testLiteralExpression(t, function.Parameters[0], "x")
// 	testLiteralExpression(t, function.Parameters[1], "y")
//
// 	if len(function.Body.Statements) != 1 {
// 		t.Fatalf("function.Body.Statements has not 1 statements. got=%d\n",
// 			len(function.Body.Statements))
// 	}
//
// 	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T",
// 			function.Body.Statements[0])
// 	}
//
// 	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
// }
//
// func TestCallExpressionParsing(t *testing.T) {
// 	input := "add(1, 2 * 3, 4 + 5);"
//
// 	l := lexer.New(input)
// 	p := New(l)
// 	program := p.ParseProgram()
// 	checkParserErrors(t, p)
//
// 	if len(program.Statements) != 1 {
// 		t.Fatalf("program.Statements does not contain %d statements. got=%d\n",
// 			1, len(program.Statements))
// 	}
//
// 	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
// 	if !ok {
// 		t.Fatalf("stmt is not ast.ExpressionStatement. got=%T",
// 			program.Statements[0])
// 	}
//
// 	exp, ok := stmt.Expression.(*ast.CallExpression)
// 	if !ok {
// 		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T",
// 			stmt.Expression)
// 	}
//
// 	if !testIdentifier(t, exp.Function, "add") {
// 		return
// 	}
// 	if len(exp.Arguments) != 3 {
// 		t.Fatalf("wrong length of arguments. got=%d", len(exp.Arguments))
// 	}
//
// 	testLiteralExpression(t, exp.Arguments[0], 1)
// 	testInfixExpression(t, exp.Arguments[1], 2, "*", 3)
// 	testInfixExpression(t, exp.Arguments[2], 4, "+", 5)
// }

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

func testVarDeclaration(t *testing.T, s ast.Statement, st string, name string) bool {
	if s.TokenLiteral() != st {
		t.Errorf("s.TokenLiteral not '%s'. got=%q", st, s.TokenLiteral())
		return false
	}

	varStmt, ok := s.(*ast.VarStatement)
	if !ok {
		t.Errorf("s not *ast.VarDeclaration. got=%T", s)
		return false
	}

	if varStmt.Name.Value != name {
		t.Errorf("varStmt.Name.Value not '%s'. got=%s", name, varStmt.Name.Value)
		return false
	}

	if varStmt.Name.TokenLiteral() != name {
		t.Errorf("s.Name not '%s'. got=%s", name, varStmt.Name)
		return false
	}

	return true
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
		return testIdentifier(t, exp, v)
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
