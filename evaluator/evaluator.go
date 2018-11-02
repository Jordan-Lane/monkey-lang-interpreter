package evaluator

import (
	"fmt"
	"monkeylang/ast"
	"monkeylang/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

func Eval(node ast.Node) object.Object {
	switch castedNode := node.(type) {
	case *ast.Program:
		return evalProgram(castedNode)
	case *ast.ExpressionStatement:
		return Eval(castedNode.Expression)
	case *ast.InfixExpression:
		left := Eval(castedNode.Left)
		if isError(left) {
			return left
		}
		right := Eval(castedNode.Right)
		if isError(right) {
			return right
		}
		return evalInfixExpression(castedNode.Operator, left, right)
	case *ast.PrefixExpression:
		right := Eval(castedNode.Right)
		if isError(right) {
			return right
		}
		return evalPrefixExpression(castedNode.Operator, right)
	case *ast.IntegerLiteral:
		return &object.Integer{Value: castedNode.Value}
	case *ast.BooleanLiteral:
		return nativeBoolToBooleanObject(castedNode.Value)
	case *ast.BlockStatement:
		return evalBlockStatement(castedNode)
	case *ast.ReturnStatement:
		value := Eval(castedNode.ReturnValue)
		if isError(value) {
			return value
		}
		return &object.ReturnValue{Value: value}
	case *ast.IfExpression:
		return evalIfExpression(castedNode)
	}

	return nil
}

func evalProgram(program *ast.Program) object.Object {
	var result object.Object

	for _, statement := range program.Statements {
		result = Eval(statement)

		switch result := result.(type) {
		case *object.ReturnValue:
			return result.Value
		case *object.Error:
			return result
		}
	}

	return result
}

func evalStatements(statements []ast.Statement) object.Object {
	var result object.Object

	for _, statement := range statements {
		result = Eval(statement)

		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue.Value
		}
	}
	return result
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var result object.Object

	for _, statement := range block.Statements {
		result = Eval(statement)

		if result != nil {
			resultType := result.Type()
			if resultType == object.RETURN_VALUE_OBJ || resultType == object.ERROR_OBJ {
				return result
			}
		}
	}

	return result
}

func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangPrefixExpression(right)
	case "-":
		return evalMinusPrefixExpression(right)
	default:
		return newError("Unknown operator: %s%s", operator, right.Type())
	}
}

func evalInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ:
		return evalIntegerInfixExpression(operator, left, right)
	case left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ:
		return evalBooleanInfixExpression(operator, left, right)
	case left.Type() != right.Type():
		return newError("Mismatch types: %s %s %s", left.Type(), operator, right.Type())
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIntegerInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Integer).Value
	rightValue := right.(*object.Integer).Value

	switch operator {
	case "+":
		return &object.Integer{Value: leftValue + rightValue}
	case "-":
		return &object.Integer{Value: leftValue - rightValue}
	case "*":
		return &object.Integer{Value: leftValue * rightValue}
	case "/":
		return &object.Integer{Value: leftValue / rightValue}
	case "<":
		return nativeBoolToBooleanObject(leftValue < rightValue)
	case ">":
		return nativeBoolToBooleanObject(leftValue > rightValue)
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalBooleanInfixExpression(operator string, left object.Object, right object.Object) object.Object {
	leftValue := left.(*object.Boolean).Value
	rightValue := right.(*object.Boolean).Value

	//TODO Add support for & and |

	switch operator {
	case "==":
		return nativeBoolToBooleanObject(leftValue == rightValue)
	case "!=":
		return nativeBoolToBooleanObject(leftValue != rightValue)
	default:
		return newError("Unknown operator: %s %s %s", left.Type(), operator, right.Type())
	}
}

func evalIfExpression(ifExpression *ast.IfExpression) object.Object {
	condition := Eval(ifExpression.Condition)
	if isError(condition) {
		return condition
	}

	if isTruthy(condition) {
		return Eval(ifExpression.Consequence)
	} else if ifExpression.Alternative != nil {
		return Eval(ifExpression.Alternative)
	} else {
		return NULL
	}
}

func evalBangPrefixExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		return FALSE
	}
}

func evalMinusPrefixExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("Unknown operator: -%s", right.Type())
	}
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

func nativeBoolToBooleanObject(boolean bool) *object.Boolean {
	if boolean {
		return TRUE
	}
	return FALSE
}

func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		return true
	}
}

func newError(message string, args ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(message, args...)}
}

func isError(obj object.Object) bool {
	if obj != nil {
		return obj.Type() == object.ERROR_OBJ
	}
	return false
}
