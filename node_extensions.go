package tplinator

import (
	"fmt"
	"strings"
)

type nodeExtension interface {
	Apply(node node, interpolator interpolator, evaluator evaluator) (*node, error)
}

type cneCondition struct {
	isSelf    bool
	condition string
	node      *node
}

type conditionalNodeExtension struct {
	conditions []*cneCondition
	elseNode   *node
}

func (cne *conditionalNodeExtension) Apply(node node, interpolator interpolator, evaluator evaluator) (*node, error) {
	for _, cneCond := range cne.conditions {
		result, err := evaluator(cneCond.condition)
		if err != nil {
			if err != errEvaluationFailed {
				return nil, err
			}
			result = false // provide default value
			logDefaultValueWarning(cneCond.condition, result)
		}
		boolResult, isBool := result.(bool)
		if !isBool {
			return nil, fmt.Errorf("error: the result of `%v` is not of boolean type", cneCond.condition)
		} else if boolResult {
			if cneCond.isSelf {
				return &node, nil
			}
			return cneCond.node, nil
		}
	}
	return cne.elseNode, nil
}

type conditionalClassNodeExtension struct {
	classConditions map[string]string
}

func (ccne *conditionalClassNodeExtension) Apply(node node, interpolator interpolator, evaluator evaluator) (*node, error) {
	hasClass, _, originalClass := node.HasAttribute("class")
	if !hasClass {
		originalClass = ""
	} else {
		originalClass = strings.TrimSpace(originalClass)
	}

	classes := strings.Fields(originalClass)

	for conditionalExpression, className := range ccne.classConditions {
		result, err := evaluator(conditionalExpression)
		if err != nil {
			if err != errEvaluationFailed {
				return nil, err
			}
			result = false // provide default value
			logDefaultValueWarning(conditionalExpression, result)
		}
		boolResult, isBool := result.(bool)
		if !isBool {
			return nil, fmt.Errorf("error: the result of `%v` is not of boolean type", conditionalExpression)
		} else if boolResult {
			classes = append(classes, className)
		}
	}

	node.ReplaceAttribute("class", strings.Join(classes, " "))

	return &node, nil
}
