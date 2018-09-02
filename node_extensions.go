package tplinator

import (
	"strconv"
	"strings"
)

type nodeExtension interface {
	Apply(node node, evaluator evaluator, boolEvaluator boolEvaluator) (*node, error)
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

func (cne *conditionalNodeExtension) Apply(node node, evaluator evaluator, boolEvaluator boolEvaluator) (*node, error) {
	for _, cneCond := range cne.conditions {
		result, err := boolEvaluator(cneCond.condition)
		if err != nil {
			return nil, err
		} else if result {
			if cneCond.isSelf {
				return &node, nil
			}
			return cneCond.node, nil
		} else {
			break
		}
	}
	return cne.elseNode, nil
}

type conditionalClassNodeExtension struct {
	classConditions map[string]string
}

func (ccne *conditionalClassNodeExtension) Apply(node node, evaluator evaluator) (*node, error) {
	originalClass, hasClass := node.attributes["class"]
	if !hasClass {
		originalClass = ""
	} else {
		originalClass = strings.TrimSpace(originalClass)
	}

	classes := []string{originalClass}
	for conditionalExpression, className := range ccne.classConditions {
		resultStr, err := evaluator(conditionalExpression)
		if err != nil {
			return nil, err
		}
		result, err := strconv.ParseBool(resultStr)
		if err != nil {
			return nil, err
		} else if result {
			classes = append(classes, className)
		}
	}
	node.attributes["class"] = strings.Join(classes, " ")

	return &node, nil
}
