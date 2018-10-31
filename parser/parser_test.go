package parser

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/lexer"
	"strconv"
	"testing"
)

func TestLetStatements(t *testing.T) {
	input := `
			let x = 5;
			let y = 10;
			let foobar = 838383;
			`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program produced %d statements instead of 3", len(program.Statements))
	}

	letStatementTests := []struct {
		expectedIdentifer string
	}{
		{"x"},
		{"y"},
		{"foobar"},
	}

	for i, test := range letStatementTests {
		currentStatement := program.Statements[i]
		if !testLetStatement(t, currentStatement, test.expectedIdentifer) {
			return
		}
	}
}

func testLetStatement(t *testing.T, statement ast.Statement, name string) bool {
	if statement.TokenLiteral() != "let" {
		t.Errorf("Statement.TokenLiteral not let. Got: %q", statement.TokenLiteral())
		return false
	}

	//This line asserts that the statement is a LetStatement (https://tour.golang.org/methods/15)
	letStatement, ok := statement.(*ast.LetStatement)
	if !ok {
		t.Errorf("Statement is not *ast.LetStatement. Got: %T", statement)
		return false
	}

	if letStatement.Name.Value != name {
		t.Errorf("LetStatement.Name.Value not: %s. Got: %s", name, letStatement.Name.Value)
		return false
	}

	if letStatement.Name.TokenLiteral() != name {
		t.Errorf("LetStatement.Name.TokenLiteral() not '%s'. Got=%s", name, letStatement.Name.TokenLiteral())
		return false
	}

	// Later, test the actual expression

	return true
}

func TestReturnStatements(t *testing.T) {
	input := `
			return 5;
			return x;
			return add(x+y);
			`
	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 3 {
		t.Fatalf("Program produced %d statements instead of 3", len(program.Statements))
	}

	for _, currentStatement := range program.Statements {
		if !testReturnStatement(t, currentStatement) {
			return
		}
	}
}

func testReturnStatement(t *testing.T, statement ast.Statement) bool {
	// Check that the statement token is a return
	if statement.TokenLiteral() != "return" {
		t.Errorf("Statement.TokenLiteral not return. Got %q", statement.TokenLiteral())
		return false
	}

	// Assert that the statement is a return statement
	_, ok := statement.(*ast.ReturnStatement)
	if !ok {
		t.Errorf("Statement is not *ast.ReturnStatement. Got %T", statement)
		return false
	}

	// Later, test the actual expression
	return true
}

func checkParseErrors(t *testing.T, parser *Parser) {
	errorMsgs := parser.Errors()

	numErrors := len(errorMsgs)
	if numErrors == 0 {
		return
	}

	t.Errorf("The parser encountered %d error(s)", numErrors)
	for _, errorMsg := range errorMsgs {
		t.Errorf("Parser Error: %s", errorMsg)
	}
	t.FailNow()
}

func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is incorrect. Expected: *ast.Expression. Got: %T", statement)
	}

	ident, ok := statement.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("Statement.Expression is incorrect. Expected: *ast.Identifier. Got %T", statement.Expression)
	}
	expectedIdent := "foobar"
	if ident.Value != expectedIdent {
		t.Errorf("ident.Value is incorrect. Expected: %q. Got %q", expectedIdent, ident)
	}
	if ident.TokenLiteral() != expectedIdent {
		t.Errorf("ident.TokenLiteral is incorrect. Expected: %q. Got %q", expectedIdent, ident.TokenLiteral())
	}
}

func TestIntegerLiteralExpression(t *testing.T) {
	input := "5;"

	lexer := lexer.New(input)
	parser := New(lexer)

	program := parser.ParseProgram()
	checkParseErrors(t, parser)

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}
	if len(program.Statements) != 1 {
		t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
	}

	statement, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("Statement is incorrect. Expected: *ast.Expression. Got: %T", statement)
	}

	literal, ok := statement.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("Statement.Expression is incorrect. Expected: *ast.IntegerLiteral. Got %T", statement.Expression)
	}

	expectedLiteral := "5"
	expectedInt, _ := strconv.ParseInt(expectedLiteral, 0, 64)
	if literal.Value != expectedInt {
		t.Errorf("integer.Value is incorrect. Expected: %d. Got %q", expectedInt, literal)
	}
	if literal.TokenLiteral() != expectedLiteral {
		t.Errorf("integer.TokenLiteral is incorrect. Expected: %q. Got %q", expectedLiteral, literal.TokenLiteral())
	}
}

func TestParsingPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input        string
		operator     string
		integerValue int64
	}{
		{"!5", "!", 5},
		{"-10", "-", 10},
	}

	for _, test := range prefixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is incorrect. Expected: *ast.Expression. Got: %T", statement)
		}

		expression, ok := statement.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("Statement.Expression is incorrect. Expected: *ast.PrefixExpression. Got %T", statement.Expression)
		}
		if expression.Operator != test.operator {
			t.Fatalf("Operator is incorrect. Expected: %s. Got %s", test.operator, expression.Operator)
		}
		if !testIntegerLiteral(t, expression.Right, test.integerValue) {
			t.Fatalf("Integer value is incorrect. Expected: %d. Got %d", test.integerValue, expression.Right)
		}

	}
}

func testIntegerLiteral(t *testing.T, integerLiteral ast.Expression, expectedValue int64) bool {
	integer, ok := integerLiteral.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("IntegerLiteral is incorrect. Expected: *ast.IntegerLiteral. Got: %T", integer)
		return false
	}

	if integer.Value != expectedValue {
		t.Fatalf("IntegerLiteral.Value is incorrect. Expected: %d. Got: %d", expectedValue, integer.Value)
		return false
	}

	expectedToken := fmt.Sprintf("%d", expectedValue)
	if integer.TokenLiteral() != expectedToken {
		t.Fatalf("IntegerLiteral.TokenLiteral is incorrect. Expected %s. Got %s", expectedToken, integer.TokenLiteral())
		return false
	}

	return true
}

func TestParsingInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  int64
		operator   string
		rightValue int64
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
	}

	for _, test := range infixTests {
		lexer := lexer.New(test.input)
		parser := New(lexer)

		program := parser.ParseProgram()
		checkParseErrors(t, parser)

		if len(program.Statements) != 1 {
			t.Fatalf("Program produced %d statements instead of 1", len(program.Statements))
		}

		statement, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("Statement is incorrect. Expected: *ast.Expression. Got: %T", statement)
		}

		expression, ok := statement.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("Statement.Expression is incorrect. Expected: *ask.InfixExpression. Got: %T", statement.Expression)
		}

		if !testIntegerLiteral(t, expression.Left, test.leftValue) {
			return
		}

		if expression.Operator != test.operator {
			t.Fatalf("Expression.Operator is incorrect. Expected: %q. Got: %q", test.operator, expression.Operator)
		}

		if !testIntegerLiteral(t, expression.Right, test.rightValue) {
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
			"-a * b;",
			"((-a) * b)",
		},
		{
			"!-a;",
			"(!(-a))",
		},
		{
			"a + b + c;",
			"((a + b) + c)",
		},
		{
			"a + b - c;",
			"((a + b) - c)",
		},
		{
			"a * b * c;",
			"((a * b) * c)",
		},
		{
			"a * b / c;",
			"((a * b) / c)",
		},
		{
			"a + b / c;",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f;",
			"(((a + (b * c)) + (d / e)) - f)",
		}, {
			"3 + 4; -5 * 5;",
			"(3 + 4)((-5) * 5)",
		},
		{
			"5 > 4 == 3 < 4;",
			"((5 > 4) == (3 < 4))",
		}, {
			"5 < 4 != 3 > 4;",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5;",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
	}

	for _, test := range tests {
		lexer := lexer.New(test.input)
		parser := New(lexer)
		program := parser.ParseProgram()
		checkParseErrors(t, parser)
		actual := program.String()
		if actual != test.expected {
			t.Errorf("expected=%q, got=%q", test.expected, actual)
		}
	}

}
