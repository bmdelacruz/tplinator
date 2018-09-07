package tplinator

import (
	"fmt"

	"github.com/Knetic/govaluate"
)

type EvaluatorParams map[string]interface{}

type Evaluator interface {
	EvaluateBool(input string, params EvaluatorParams) (bool, error)
	EvaluateString(input string, params EvaluatorParams) (string, error)
	Evaluate(input string, params EvaluatorParams) (interface{}, error)
}

func TryEvaluateBoolOnContext(ctxNode *Node, evaluator Evaluator, inputStr string) (bool, bool, error) {
	var result bool
	var err error

	hasResult := false
	for !hasResult && ctxNode != nil {
		if ctxParams := ctxNode.ContextParams(); ctxParams != nil {
			result, err = evaluator.EvaluateBool(inputStr, ctxParams)
			hasResult = err == nil
		}
		ctxNode = ctxNode.Parent()
	}
	return hasResult, result, err
}

func TryEvaluateStringOnContext(ctxNode *Node, evaluator Evaluator, inputStr string) (bool, string, error) {
	var result string
	var err error

	hasResult := false
	for !hasResult && ctxNode != nil {
		if ctxParams := ctxNode.ContextParams(); ctxParams != nil {
			result, err = evaluator.EvaluateString(inputStr, ctxParams)
			hasResult = err == nil
		}
		ctxNode = ctxNode.Parent()
	}
	return hasResult, result, err
}

func TryEvaluateOnContext(ctxNode *Node, evaluator Evaluator, inputStr string) (bool, interface{}, error) {
	var result interface{}
	var err error

	hasResult := false
	for !hasResult && ctxNode != nil {
		if ctxParams := ctxNode.ContextParams(); ctxParams != nil {
			result, err = evaluator.Evaluate(inputStr, ctxParams)
			hasResult = err == nil
		}
		ctxNode = ctxNode.Parent()
	}
	return hasResult, result, err
}

type govaluator struct {
}

func (e *govaluator) EvaluateBool(input string, params EvaluatorParams) (bool, error) {
	expr, err := govaluate.NewEvaluableExpression(input)
	if err != nil {
		return false, err
	}
	result, err := expr.Evaluate(params)
	if err != nil {
		return false, err
	}
	boolResult, isBoolean := result.(bool)
	if !isBoolean {
		return false, fmt.Errorf("evaluator: `%v` is not a conditional expression", input)
	}
	return boolResult, nil
}

func (e *govaluator) EvaluateString(input string, params EvaluatorParams) (string, error) {
	expr, err := govaluate.NewEvaluableExpression(input)
	if err != nil {
		return "", err
	}
	result, err := expr.Evaluate(params)
	if err != nil {
		return "", err
	}
	stringResult, isString := result.(string)
	if !isString {
		return "", fmt.Errorf("evaluator: `%v` is not an expression that returns a string", input)
	}
	return stringResult, nil
}

func (e *govaluator) Evaluate(input string, params EvaluatorParams) (interface{}, error) {
	expr, err := govaluate.NewEvaluableExpression(input)
	if err != nil {
		return nil, err
	}
	result, err := expr.Evaluate(params)
	if err != nil {
		return nil, err
	}
	return result, nil
}
